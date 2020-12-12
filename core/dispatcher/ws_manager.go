///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package dispatcher

import (
	"digger/models"
	"digger/scheduler"
	"digger/services/service"
	"errors"
	"github.com/gorilla/websocket"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	jsoniter "github.com/json-iterator/go"
	"math/rand"
	"sync"
	"time"
)

var (
	clientMap          = make(map[int]*WsClient)
	clientStatisticMap = make(map[int]*models.Node)
	wsManageLock       = new(sync.Mutex)
)

type WsMessage struct {
	ClientId int    `json:"clientId"`
	Command  int    `json:"command"`
	Data     string `json:"data"`
}

type WsClient struct {
	ClientId   int
	Connection *websocket.Conn
	writeLock  *sync.Mutex
	Closed     bool
	Labels     map[string]string
}

func (c *WsClient) WriteString(cmd int, msg string) error {
	return c.Connection.WriteJSON(&WsMessage{
		ClientId: c.ClientId,
		Command:  CMD_NONE,
		Data:     msg,
	})
}

func (c *WsClient) Close() {
	c.Closed = true
	c.Connection.Close()
	c.Labels = nil
}

func (c *WsClient) handleReceive() {
listen:
	for !c.Closed {
		var msg WsMessage
		err := c.Connection.ReadJSON(&msg)
		if err != nil {
			logger.Error("ws读取错误: ", err)
			break
		}

		switch msg.Command {
		case CMD_HEART_BEAT:
			if err = handleHeartBeat(&msg, c); err != nil {
				break listen
			}
			break
		case CMD_QUEUE_RESULT:
			handleQueueResult(&msg)
			break
		}
	}
	logger.Info("ws中断监听")
	wsManageLock.Lock()
	defer wsManageLock.Unlock()
	conn := clientMap[c.ClientId]
	if conn != nil {
		conn.Close()
		delete(clientMap, c.ClientId)
	}
	oldSta := clientStatisticMap[c.ClientId]
	if oldSta != nil {
		oldSta.Down = oldSta.Down + 1
		oldSta.Status = 0
	}
}

func GetNodes() []*models.Node {
	wsManageLock.Lock()
	defer wsManageLock.Unlock()

	var ret []*models.Node
	for _, v := range clientStatisticMap {
		ret = append(ret, v)
	}
	return ret
}

func handleHeartBeat(msg *WsMessage, c *WsClient) error {
	// logger.Info("心跳，Id: ", msg.ClientId, "; data: ", msg.Data)
	c.writeLock.Lock()
	defer c.writeLock.Unlock()
	return c.Connection.WriteJSON(msg)
}

func handleQueueResult(msg *WsMessage) {
	wsManageLock.Lock()
	oldSta := clientStatisticMap[msg.ClientId]
	wsManageLock.Unlock()

	var ret models.QueueProcessResult
	err := jsoniter.UnmarshalFromString(msg.Data, &ret)
	if err != nil {
		oldSta.Error = oldSta.Error + 1
		logger.Error("无法处理queue结果: ", err)
		return
	}
	if ret.QueueId == 0 || ret.RequestId == "" && ret.Expire < gox.GetTimestamp(time.Now()) {
		oldSta.Error = oldSta.Error + 1
		logger.Info("请求非法，丢弃结果")
		return
	}

	logFile := logFileMap[ret.TaskId]
	if logFile != nil {
		if ret.Error != "" {
			logFile.WriteString(ret.Error)
		} else {
			logFile.WriteString(ret.Logs)
		}
	}
	defer func() {
		// ??
		ret.Results = nil
		ret.NewQueues = nil
		ret.Logs = ""
		if logFile != nil {
			logFile.Sync()
		}
	}()
	if ret.Error != "" {
		updateLock.Lock()
		errorRequests++
		updateLock.Unlock()
		oldSta.Error = oldSta.Error + 1
		logger.Error("worker处理返回异常：", ret.Error)
		n, err := service.CacheService().IncreQueueErrorCount([]int{ret.TaskId}, []int64{ret.QueueId})
		if err != nil {
			logger.Error("无法标记错误次数：", err)
		}
		if len(n) > 0 && n[0] > 3 {
			logger.Warn("重试超过阈值，丢弃，queueId: ", ret.QueueId)
			if err = service.ResultService().SaveProcessResultData(&ret, true); err != nil {
				logger.Info(err)
			}
		}
		// 将并发移除
		service.CacheService().DecreConcurrentTaskCount(ret.RequestId, ret.TaskId)
		return
	}

	updateLock.Lock()
	results += len(ret.Results)
	updateLock.Unlock()

	if err = service.ResultService().SaveProcessResultData(&ret, false); err != nil {
		logger.Info(err)
	}
	service.CacheService().SaveSuccessQueueIds(&models.QueueCallbackRequestVO{
		SuccessQueueIds:     []int64{ret.QueueId},
		SuccessQueueTaskIds: []int{ret.TaskId},
	})
	if ret.InitUrl != "" {
		if err := service.CacheService().AddFinishUniqueRes(ret.TaskId, ret.InitUrl); err != nil {
			logger.Info(err)
		}
	}
	// 将并发移除
	service.CacheService().DecreConcurrentTaskCount(ret.RequestId, ret.TaskId)
	logger.Debug("请求处理成功：", ret.QueueId)
	oldSta.Success = oldSta.Success + 1
	scheduler.BackPushNotify(ret.TaskId)
}

func shutdownClient(c *WsClient) error {
	c.writeLock.Lock()
	defer c.writeLock.Unlock()

	logger.Info("通知worker下线: ", c.ClientId)

	return c.Connection.WriteJSON(&WsMessage{
		ClientId: c.ClientId,
		Command:  CMD_SHUTDOWN_WORKER,
	})
}

// 注册连接
func RegisterWsConnection(conn *websocket.Conn, remoteHost string) error {
	wsManageLock.Lock()
	defer wsManageLock.Unlock()

	var req WsMessage
	err := conn.ReadJSON(&req)
	if err != nil {
		return err
	}

	if clientMap[req.ClientId] != nil {
		logger.Info("实例ID冲突:", req.ClientId)
		return errors.New("conflict instance id: " + convert.IntToStr(req.ClientId))
	}

	logger.Info("worker注册，Id: ", req.ClientId, ", 标签: ", req.Data)

	_labels := make(map[string]string)
	err = jsoniter.UnmarshalFromString(req.Data, &_labels)
	if err != nil {
		return err
	}
	client := &WsClient{
		ClientId:   req.ClientId,
		Labels:     _labels,
		Connection: conn,
		writeLock:  new(sync.Mutex),
		Closed:     false,
	}

	if err = client.WriteString(CMD_NONE, "ok"); err != nil {
		return err
	}

	logger.Info("worker注册成功")
	clientMap[req.ClientId] = client
	oldSta := clientStatisticMap[req.ClientId]
	if oldSta == nil {
		oldSta = &models.Node{
			RegisterAt: gox.GetLongDateString(time.Now()),
		}
	}
	clientStatisticMap[req.ClientId] = &models.Node{
		InstanceId: req.ClientId,
		RegisterAt: oldSta.RegisterAt,
		Address:    remoteHost,
		Down:       oldSta.Down,
		Assign:     oldSta.Assign,
		Success:    oldSta.Success,
		Error:      oldSta.Error,
		Status:     1,
		Labels:     _labels,
	}
	go client.handleReceive()
	return nil
}

// 随机选择一个worker节点
func selectClient(labels []models.KV) *WsClient {
	wsManageLock.Lock()
	defer wsManageLock.Unlock()

	var retList []*WsClient
	for _, c := range clientMap {
		if c.Labels == nil {
			continue
		}
		allMatch := true
		for _, kv := range labels {
			if c.Labels[kv.Key] != kv.Value {
				allMatch = false
				break
			}
		}
		if allMatch {
			retList = append(retList, c)
		}
	}
	if len(retList) == 0 {
		return nil
	}
	index := rand.Intn(len(retList))
	return retList[index]
}

func CountClient() int {
	wsManageLock.Lock()
	defer wsManageLock.Unlock()
	return len(clientMap)
}

func selectClientById(id int) *WsClient {
	wsManageLock.Lock()
	defer wsManageLock.Unlock()

	for _, c := range clientMap {
		if c.ClientId == id {
			return c
		}
	}
	return nil
}

func handleRecvMsg() {

}
