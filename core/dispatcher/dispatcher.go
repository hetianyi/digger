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
	// statistic
	assignRequests int
	errorRequests  int
	results        int
	updateLock = new(sync.Mutex)
)

func StartDispatcher(_config *models.BootstrapConfig) {
	config = _config
	// 开启定时扫描任务
	scheduleScanTask()
}

func scheduleScanTask() {
	// 定时任务扫描任务状态
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
	// 定时任务记录统计数据
	data := make(map[string]interface{})
	reqCount := 0
	reqCountLock := new(sync.Mutex)
	timer.Start(0, time.Second*20, 0, func(t *timer.Timer) {
		updateLock.Lock()
		reqCountLock.Lock()
		defer updateLock.Unlock()
		defer reqCountLock.Unlock()

		/*if assignRequests == 0 && assignRequests == 0 && errorRequests == 0 {
			return
		}*/
		// {"request_count":120,"error_request_count":5,"result_count":206}
		data["request_count"] = assignRequests
		data["error_request_count"] = errorRequests
		data["result_count"] = results
		reqCount += assignRequests

		if err := service.StatisticService().Save(data); err != nil {
			logger.Error(err)
		} else {
			assignRequests = 0
			errorRequests = 0
			results = 0
		}
	})

	// 定时任务统计全局数据，包括项目总数，任务总数，请求总数，结果总数等
	timer.Start(0, time.Second*10, 0, func(t *timer.Timer) {
		config, err := service.ConfigService().ListConfigs()
		if err != nil {
			logger.Error(err)
			return
		}

		var resultCountSince int64 = 0
		var resultCount int64 = -1
		var totalReqCount int64 = -1
		if config["total_request_count"] == "" {
			config["total_request_count"] = "0"
		}
		if i, err := convert.StrToInt64(config["total_request_count"]); err == nil {
			totalReqCount = i
		}
		if totalReqCount != -1 {
			reqCountLock.Lock()
			totalReqCount += int64(reqCount)
			reqCount = 0
			reqCountLock.Unlock()
			if err =service.ConfigService().UpdateConfig("total_request_count", convert.Int64ToStr(totalReqCount)); err != nil {
				logger.Error(err)
			}
		}

		if config["result_count_since"] == "" {
			config["result_count_since"] = "0"
		}
		if i, err := convert.StrToInt64(config["result_count_since"]); err == nil {
			resultCountSince = i
		}
		if config["result_count"] == "" {
			config["result_count"] = "0"
		}
		if i, err := convert.StrToInt64(config["result_count"]); err == nil {
			resultCount = i
		}

		pc, err := service.ProjectService().AllProjectCount()
		if err != nil {
			logger.Error(err)
		} else {
			if err =service.ConfigService().UpdateConfig("project_count", convert.IntToStr(pc)); err != nil {
				logger.Error(err)
			}
		}
		tc, err := service.TaskService().AllTaskCount()
		if err != nil {
			logger.Error(err)
		} else {
			if err =service.ConfigService().UpdateConfig("task_count", convert.IntToStr(tc)); err != nil {
				logger.Error(err)
			}
		}

		if resultCount == -1 {
			return
		}
		c, nextId, err := service.ResultService().ResultCountSince(resultCountSince)
		if err != nil {
			logger.Error(err)
		} else {
			resultCount += int64(c) + 1
			resultCountSince = nextId
			if err =service.ConfigService().UpdateConfig("result_count", convert.Int64ToStr(resultCount)); err != nil {
				logger.Error(err)
			}
			if err =service.ConfigService().UpdateConfig("result_count_since", convert.Int64ToStr(resultCountSince)); err != nil {
				logger.Error(err)
			}
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
	updateLock.Lock()
	assignRequests++
	updateLock.Unlock()
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
