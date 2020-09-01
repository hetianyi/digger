///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package dispatcher

import (
	"digger/common"
	"digger/models"
	"digger/scheduler"
	"digger/services/service"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/queue"
	"github.com/hetianyi/gox/timer"
	"github.com/hetianyi/gox/uuid"
	jsoniter "github.com/json-iterator/go"
	"os"
	"sync"
	"time"
)

const (
	CMD_NONE            = 0 // 无
	CMD_REG             = 1 // 注册
	CMD_HEART_BEAT      = 2 // 心跳
	CMD_DISPATCH_QUEUE  = 3 // 分配任务
	CMD_QUEUE_RESULT    = 4 // 返回任务结果
	CMD_SHUTDOWN_WORKER = 5 // 强制断开worker
)

var (
	dispatchLock      = new(sync.Mutex)
	taskScheduleTimer = make(map[int]*timer.Timer)
	logFileMap        = make(map[int]*os.File)
	taskWorkLock      = make(map[int]*sync.Mutex)
	notifiers         = make(map[int]*queue.NoneBlockQueue)
	config            *models.BootstrapConfig
)

func StartDispatcher(_config *models.BootstrapConfig) {
	config = _config
	// 开启定时扫描任务
	scheduleScanTask()
}

func scheduleScanTask() {
	timer.Start(0, time.Second*30, 0, func(t *timer.Timer) {
		tasks, err := service.TaskService().SelectActiveTasks()
		if err != nil {
			logger.Error("cannot get active tasks: ", err)
		}
		if len(tasks) == 0 {
			return
		}
		for _, task := range tasks {
			DispatchTask(task)
		}
	})
}

func DispatchTask(task *models.Task) {
	dispatchLock.Lock()
	defer dispatchLock.Unlock()

	if taskScheduleTimer[task.Id] != nil {
		return
	}

	logFile := config.LogDir + "/" + convert.IntToStr(task.Id) + ".log"

	appendFile, err := file.AppendFile(logFile)
	if err != nil {
		logger.Error("error start task dispatcher: ", err)
		appendFile.Close()
	}

	// 获取配置
	project, err := service.CacheService().GetSnapshotConfig(task.Id)
	if err != nil {
		logger.Error(err)
		appendFile.Close()
		return
	}
	// 注册通知
	notifier := scheduler.RegisterNotifier(task.Id)
	if notifiers[task.Id] == nil {
		notifiers[task.Id] = notifier
	}
	if taskWorkLock[task.Id] == nil {
		taskWorkLock[task.Id] = new(sync.Mutex)
	}
	taskSig[task.Id] = 1
	// 并发配置
	conSize := project.GetIntSetting(common.SETTINGS_CONCURRENT_REQUESTS)

	var doWork = func() bool {
		taskWorkLock[task.Id].Lock()
		defer taskWorkLock[task.Id].Unlock()

		if taskSig[task.Id] != common.RUNNING {
			return false
		}

		// 执行任务必须有worker
		client := selectClient(project.NodeAffinityParsed)
		if client == nil {
			logger.Debug("没有worker，无法调度任务")
			return false
		}
		// 先取并发锁，拿不到则说明达到最大并发
		requestId := uuid.UUID()
		if !service.CacheService().IncreConcurrentTaskCount(requestId, task.Id, conSize) {
			// 达到最大并发，暂停一会
			return false
		}
		releaseLock := false
		defer func() {
			if releaseLock {
				// 释放并发锁
				service.CacheService().DecreConcurrentTaskCount(requestId, task.Id)
			}
		}()
		// 拿到并发锁，分配任务
		queue := scheduler.FetchQueue(task.Id)
		if queue == nil {
			logger.Debug("当前没有任务")
			releaseLock = true
			return false
		}
		if queue.Expire < gox.GetTimestamp(time.Now()) {
			logger.Debug("拿到的任务已过期：", queue.Id)
			releaseLock = true
			return true
		}

		// 检查是否是unique类型的stage
		stage := project.GetStageByName(queue.StageName)
		if stage == nil {
			logger.Debug("stage不存在：", stage)
			releaseLock = true
			return true
		}

		if stage.IsUnique && service.CacheService().IsUniqueResFinish(queue.TaskId, queue.Url) {
			logger.Info("资源已完成，跳过：", queue.Id)
			if err := service.ResultService().SaveProcessResultData(&models.QueueProcessResult{
				TaskId:    queue.TaskId,
				QueueId:   queue.Id,
				Expire:    queue.Expire,
				NewQueues: nil,
				Results:   nil,
			}, false); err != nil {
				logger.Error(err)
			}
			releaseLock = true
			return true
		}

		dispatchWork(requestId, queue, client)
		return true
	}

	go notify(task.Id, doWork)

	logFileMap[task.Id] = appendFile

	taskScheduleTimer[task.Id] = timer.Start(0, time.Second*3, 0, func(t *timer.Timer) {
		for {
			logger.Debug("检测任务")
			if !doWork() {
				break
			}
		}
	})
}

func notify(taskId int, doWork func() bool) {
	for {
		blockQueue := notifiers[taskId]
		lock := taskWorkLock[taskId]
		if blockQueue == nil || lock == nil {
			break
		}
		lock.Lock()
		_, s := blockQueue.Fetch()
		lock.Unlock()
		if s {
			logger.Debug("notify推送")
			gox.Try(func() {
				doWork()
			}, func(e interface{}) {
				logger.Error(e)
			})
			continue
		}
		time.Sleep(time.Second)
	}
}

func dispatchWork(requestId string, queue *models.Queue, client *WsClient) {
	logger.Debug("分配任务：", queue.Id, "，requestId=", requestId)
	push(client, CMD_DISPATCH_QUEUE, &models.DispatchWork{
		RequestId: requestId,
		Queue:     queue,
	})
}

func push(client *WsClient, cmd int, data interface{}) error {
	wsManageLock.Lock()
	oldSta := clientStatisticMap[client.ClientId]
	wsManageLock.Unlock()

	d, err := jsoniter.MarshalToString(data)
	if err != nil {
		return err
	}

	req := &WsMessage{
		ClientId: 0,
		Command:  cmd,
		Data:     d,
	}
	client.writeLock.Lock()
	defer client.writeLock.Unlock()

	if oldSta != nil {
		oldSta.Assign = oldSta.Assign + 1
	}

	return client.Connection.WriteJSON(req)
}
