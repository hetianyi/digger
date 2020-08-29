///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package client

import (
	"digger/dispatcher"
	"digger/models"
	"digger/utils"
	"github.com/gorilla/websocket"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/timer"
	jsoniter "github.com/json-iterator/go"
	"net/url"
	"os"
	"sync"
	"time"
)

var (
	conn      *websocket.Conn
	config    *models.BootstrapConfig
	heartBeat *timer.Timer
	writeLock = new(sync.Mutex)
)

func InitWsClient(_config *models.BootstrapConfig) {
	config = _config

	managerUrl, err := utils.Parse(config.ManagerUrl)
	if err != nil {
		logger.Fatal("invalid manager url: ", config.ManagerUrl, ": ", err.Error())
	}

	u := url.URL{Scheme: gox.TValue(managerUrl.Scheme == "https", "wss", "ws").(string), Host: managerUrl.Host, Path: "/api/ws"}
	for {
		c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
		if err != nil {
			logger.Error(err)
			time.Sleep(time.Second * 5)
			continue
		}
		logger.Info("WS连接成功")
		conn = c
		if err := handleRegister(); err != nil {
			c.Close()
			logger.Error("ws注册失败:", err)
			time.Sleep(time.Second * 5)
			continue
		}
		logger.Info("ws注册成功")
		heartBeat = handleCheckHealth()
		go handleRecvMsg(c)
		break
	}
}

func handleRecvMsg(c *websocket.Conn) {
	for {
		var msg dispatcher.WsMessage
		err := c.ReadJSON(&msg)
		if err != nil {
			logger.Error(err)
			c.Close()
			break
		}

		switch msg.Command {
		case dispatcher.CMD_DISPATCH_QUEUE:
			go handleQueue(msg.Data)
			break
		case dispatcher.CMD_SHUTDOWN_WORKER:
			handleShutdown()
			break
		}
	}
	logger.Warn("ws已断开")
	if heartBeat != nil {
		heartBeat.Destroy()
		heartBeat = nil
	}
	go InitWsClient(config)
}

func handleShutdown() {
	logger.Info("服务端强制通知下线")
	os.Exit(0)
}

func handleQueue(data string) {

	var work models.DispatchWork
	if err := jsoniter.UnmarshalFromString(data, &work); err != nil {
		return
	}

	queue := work.Queue

	if queue.Id == 0 || queue.TaskId == 0 {
		return
	}

	logger.Info("处理queue：", work.RequestId)

	result, err := processQueue(queue)
	if err != nil {
		logger.Error("处理queue错误：", err)
		result = &models.QueueProcessResult{
			TaskId:    work.Queue.TaskId,
			InitUrl:   work.Queue.Url,
			QueueId:   work.Queue.Id,
			Expire:    work.Queue.Expire,
			RequestId: work.RequestId,
			Error:     "<span style=\"color:#F38F8F\">Err: " + err.Error() + "</span>\n",
			NewQueues: nil,
			Results:   nil,
		}
		s, _ := jsoniter.MarshalToString(result)
		writeResult(conn, s)
		return
	}

	if result == nil {
		result = &models.QueueProcessResult{
			TaskId:    work.Queue.TaskId,
			InitUrl:   work.Queue.Url,
			QueueId:   work.Queue.Id,
			Expire:    work.Queue.Expire,
			RequestId: work.RequestId,
			Error:     "no result",
			NewQueues: nil,
			Results:   nil,
		}
		s, _ := jsoniter.MarshalToString(result)
		writeResult(conn, s)
		return
	}
	result.RequestId = work.RequestId
	result.InitUrl = work.Queue.Url

	s, _ := jsoniter.MarshalToString(result)
	if conn != nil {
		logger.Info("queue处理成功：", work.RequestId)
		writeResult(conn, s)
	}
}

func writeResult(conn *websocket.Conn, result string) {
	writeLock.Lock()
	defer writeLock.Unlock()

	if conn != nil {
		if err := conn.WriteJSON(&dispatcher.WsMessage{
			ClientId: config.InstanceId,
			Command:  dispatcher.CMD_QUEUE_RESULT,
			Data:     result,
		}); err != nil {
			logger.Error(err)
			conn.Close()
		}
	}
}

func handleCheckHealth() *timer.Timer {
	return timer.Start(0, time.Second*10, 0, func(t *timer.Timer) {
		if conn != nil {
			//logger.Info("执行ws健康检查")

			writeLock.Lock()
			defer writeLock.Unlock()

			if err := conn.WriteJSON(&dispatcher.WsMessage{
				ClientId: config.InstanceId,
				Command:  dispatcher.CMD_HEART_BEAT,
				Data:     "hi",
			}); err != nil {
				logger.Error(err)
				conn.Close()
			}
		}
	})
}

func handleRegister() error {
	s, err := jsoniter.MarshalToString(config.Labels)
	if err != nil {
		return err
	}

	writeLock.Lock()
	defer writeLock.Unlock()

	if err = conn.WriteJSON(&dispatcher.WsMessage{
		ClientId: config.InstanceId,
		Command:  dispatcher.CMD_REG,
		Data:     s,
	}); err != nil {
		return err
	}
	var ret dispatcher.WsMessage
	err = conn.ReadJSON(&ret)
	if err != nil {
		return err
	}

	if ret.Command == dispatcher.CMD_SHUTDOWN_WORKER {
		handleShutdown()
	}
	return nil
}
