///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package event

import (
	"digger/models"
	"digger/scheduler"
	"digger/services/service"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/queue"
	"sync"
	"time"
)

var (
	lock                   = new(sync.Mutex)
	queueFinishNotifyChans = make(map[int]*queue.NoneBlockQueue)
)

func getBlockQueue(taskId int) *queue.NoneBlockQueue {
	lock.Lock()
	defer lock.Unlock()

	q := queueFinishNotifyChans[taskId]
	if q == nil {
		q = queue.NewNoneBlockQueue(100)
		queueFinishNotifyChans[taskId] = q
	}
	return q
}

// 处理task创建的事件
func handleTaskCreatedEvent(event *models.RedisEvent) {
	if event.Body != nil {
		taskId, _ := convert.StrToInt(event.Body["taskId"])
		logger.Debug("start create task ", taskId)
		task, err := service.TaskService().SelectTask(taskId)
		if err != nil {
			logger.Error(err)
			return
		}
		scheduler.Schedule(task)
	}
}

// 处理task暂停的事件
func handleTaskPauseEvent(event *models.RedisEvent) {
	if event.Body != nil {
		taskId, _ := convert.StrToInt(event.Body["taskId"])
		logger.Debug("start create task ", taskId)
		//scheduler.Pause(taskId)
	}
}

// 处理task停止的事件
func handleTaskStopEvent(event *models.RedisEvent) {
	if event.Body != nil {
		taskId, _ := convert.StrToInt(event.Body["taskId"])
		logger.Debug("start create task ", taskId)
		//scheduler.Stop(taskId)
	}
}

// 处理task继续的事件
func handleTaskContinueEvent(event *models.RedisEvent) {
	if event.Body != nil {
		taskId, _ := convert.StrToInt(event.Body["taskId"])
		logger.Debug("start create task ", taskId)
		//scheduler.Continue(taskId)
	}
}

// 处理queue处理完成事件
func handleQueueFinishEvent(event *models.RedisEvent) {
	if event.Body != nil {
		taskId, _ := convert.StrToInt(event.Body["taskId"])
		requestId := event.Body["requestId"]
		getBlockQueue(taskId).Put(requestId)
	}
}

// 用于请求并发数控制，当并发数达到阈值时，
// 则阻塞等待其他请求结束发送订阅消息，然后阻塞的请求继续竞争锁
func WaitConcurrentLock(taskId int, timeout time.Duration) bool {
	//logger.Info("等待并发信号")
	expireTime := time.Now().Add(timeout)
	blockQueue := getBlockQueue(taskId)
	requestId := ""
	for {
		lock.Lock()
		qid, success := blockQueue.Fetch()
		lock.Unlock()
		if !success {
			if time.Now().Unix() > expireTime.Unix() {
				return false
			}
			time.Sleep(time.Millisecond * 100)
			continue
		}
		requestId = qid.(string)
		break
	}
	logger.Info("收到并发信号:", requestId)
	return true
}
