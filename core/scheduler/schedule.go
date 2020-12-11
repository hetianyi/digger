///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package scheduler

import (
	"digger/common"
	"digger/models"
	"digger/services/service"
	"fmt"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/queue"
	"github.com/hetianyi/gox/timer"
	"sync"
	"time"
)

type StopFunc interface {
	Stop(taskId int)
}

var (
	scheduleLock = new(sync.Mutex)
	// blockQueueCache 是存放task的queue的缓冲容器
	blockQueueCache = make(map[int]*queue.NoneBlockQueue)
	// queueNotifier 是dispatcher用于通知schedule冲数据库查询queue的信号channel
	queueNotifier = make(map[int]*queue.NoneBlockQueue)
	// queueNotifier 是schedule通知dispatcher可以从schedule取queue的信号channel
	backPushNotifier = make(map[int]*queue.NoneBlockQueue)
	queryLock        = make(map[int]*sync.Mutex)
	// taskFinishTimer 是一个定时任务，用于定时检测任务是否达到完成条件
	taskFinishTimer = make(map[int]*timer.Timer)
	// stopFunc 在dispatcher中实现，用来给schedule通知dispatcher任务可以停止了
	stopFunc StopFunc
)

// 为dispatcher注册停止函数
func RegisterStopFunc(_stopFunc StopFunc) {
	stopFunc = _stopFunc
}

// 通知dispatcher任务停止
func Stop(taskId int) {
	stopFunc.Stop(taskId)
}

// 启动任务启动扫描任务和queue状态定时重置
func StartScheduler() {
	scheduleScanTask()
	scheduleResetQueue()
	scheduleScanPushTask()
}

// 定时扫描任务，适用于程序启动时恢复task状态
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
			Schedule(task)
		}
	})
}

// 定时重置任务状态
func scheduleResetQueue() {
	timer.Start(0, time.Second*10, 0, func(t *timer.Timer) {
		err := service.QueueService().ResetQueuesStatus()
		if err != nil {
			logger.Error("cannot reset queues: ", err)
		}
	})
}

// dispatcher调用，用于获取当前task的可用queue，用于worker消费
func FetchQueue(taskId int) *models.Queue {
	if blockQueueCache[taskId] == nil {
		return nil
	}
	q, s := blockQueueCache[taskId].Fetch()
	if s && q != nil {
		return q.(*models.Queue)
	}
	return nil
}

// dispatcher调用，用于注册并获取通知channel
func RegisterNotifier(taskId int) *queue.NoneBlockQueue {
	scheduleLock.Lock()
	defer scheduleLock.Unlock()

	if queueNotifier[taskId] == nil {
		queueNotifier[taskId] = queue.NewNoneBlockQueue(1)
	}
	return queueNotifier[taskId]
}

// dispatcher取消注册通知channel
func DeRegisterNotifier(taskId int) {
	scheduleLock.Lock()
	defer scheduleLock.Unlock()

	delete(queueNotifier, taskId)
}

// dispatcher处理一条queue成功，通知schedule可以继续读取queue
func BackPushNotify(taskId int) {
	q := backPushNotifier[taskId]
	if q != nil {
		q.Put(1)
	}
}

// 轮询backPushNotifier，接收到信号就查询queue
func backPushListener(taskId, fetchSize, queueExpireSeconds int) {
	for {
		time.Sleep(time.Millisecond * 100)
		q := backPushNotifier[taskId]
		if q == nil {
			break
		}
		fetch, s := q.Fetch()
		if !s || fetch == nil {
			continue
		}
		// 查询是否还有库存，有则忽略，否则查询queue
		store := FetchQueue(taskId)
		if store != nil {
			doQueue(store)
			continue
		}
		selectTask, _ := service.TaskService().SelectTask(taskId)
		// 如果task不存在或者不是在运行中，则忽略
		if selectTask == nil || selectTask.Status != 1 {
			continue
		}
		fetchQueue(taskId, fetchSize, queueExpireSeconds)
	}
}

// 为一个task启动一个调度器，
// 调度器的任务是从数据库t_queue表轮询查询队列，并分配任务给worker节点
func Schedule(task *models.Task) error {
	scheduleLock.Lock()
	defer scheduleLock.Unlock()

	//
	if blockQueueCache[task.Id] != nil {
		return nil
	}

	// 获取配置快照
	project, err := service.CacheService().GetSnapshotConfig(task.Id)
	if err != nil {
		return err
	}

	// 并发配置
	conSize := project.GetIntSetting(common.SETTINGS_CONCURRENT_REQUESTS, 5)
	// queue过期配置
	queueExpireSeconds := project.GetIntSetting(common.SETTINGS_QUEUE_EXPIRE_SECONDS, 10)
	blockQueueCache[task.Id] = queue.NewNoneBlockQueue(conSize * 5)
	queryLock[task.Id] = new(sync.Mutex)
	backPushNotifier[task.Id] = queue.NewNoneBlockQueue(1)
	// 启动定时器检测任务是否达到完成状态
	checkStatusFinishStatus(task.Id)

	go backPushListener(task.Id, conSize, queueExpireSeconds)

	logger.Info("开始调度任务：", task.Id)

	timer.Start(time.Second*2, time.Second*10, 0, func(t *timer.Timer) {
		for {
			selectTask, err := service.TaskService().SelectTask(task.Id)
			if err != nil {
				logger.Error(err)
				return
			}
			if selectTask == nil {
				scheduleLock.Lock()
				logger.Info("任务不存在：", task.Id)
				// 清理task相关数据
				cleanTaskData(task.Id)
				scheduleLock.Unlock()
				t.Destroy()
				return
			}
			if selectTask.Status == 1 {
				if fetchQueue(task.Id, conSize, queueExpireSeconds) {
					continue
				}
			} else if selectTask.Status == 2 || selectTask.Status == 3 {
				scheduleLock.Lock()
				if selectTask.Status == 2 {
					logger.Info("停止任务：", task.Id)
				} else {
					logger.Info("任务已完成：", task.Id)
				}
				// 清理task相关数据
				cleanTaskData(task.Id)
				scheduleLock.Unlock()
				t.Destroy()
			}
			break
		}
	})

	return nil
}

func cleanTaskData(taskId int) {
	// 通知dispatcher停止任务
	Stop(taskId)
	// 释放blockQueueCache的任务
	releaseTaskQueueCache(blockQueueCache[taskId])
	delete(blockQueueCache, taskId)
	delete(backPushNotifier, taskId)
	delete(queueNotifier, taskId)
	delete(queryLock, taskId)
	delete(taskFinishTimer, taskId)
}

// 清空task的block queue
func releaseTaskQueueCache(blockQueue *queue.NoneBlockQueue) {
	if blockQueue == nil {
		return
	}
	for {
		i, s := blockQueue.Fetch()
		if s {
			logger.Debug("释放queue：", i)
			continue
		}
		break
	}
}

func fetchQueue(taskId int, fetchSize, queueExpireSeconds int) bool {
	lock := queryLock[taskId]
	if lock != nil {
		lock.Lock()
		defer lock.Unlock()
	} else {
		return false
	}

	var queues []*models.Queue
	queues, err := service.QueueService().SelectQueues(models.QueueQueryVO{
		TaskId:             taskId,
		LockStatus:         true,
		Status:             0,
		QueueExpireSeconds: queueExpireSeconds,
		Limit:              fetchSize,
	})
	if err != nil {
		logger.Debug("error query queue from database: ", err)
		return false
	}
	if len(queues) > 0 {
		logger.Info(fmt.Sprintf("加载了%d条queue", len(queues)))
	}
	doQueue(queues...)
	if len(queues) > 0 {
		return true
	}
	return false
}

// 将查询的queue列表放入队列，如果队列已满，则会阻塞等待
func doQueue(queues ...*models.Queue) {
	for _, q := range queues {
		blockQueue := blockQueueCache[q.TaskId]
		if blockQueue == nil {
			break
		}
		for !blockQueue.Put(q) {
			time.Sleep(time.Millisecond * 100)
		}
		// 通知dispatcher有活干了
		if queueNotifier[q.TaskId] != nil {
			queueNotifier[q.TaskId].Put(1)
		}
	}
}

func checkStatusFinishStatus(taskId int) {
	if taskFinishTimer[taskId] != nil {
		return
	}
	checkExist := func(t *timer.Timer) bool {
		scheduleLock.Lock()
		defer scheduleLock.Unlock()
		task, err := service.TaskService().SelectTask(taskId)
		if err != nil {
			logger.Error(err)
			return true
		}
		if task == nil || task.Status == 2 || task.Status == 3 {
			t.Destroy()
			return false
		}
		return true
	}
	taskFinishTimer[taskId] = timer.Start(0, time.Second*10, 0, func(t *timer.Timer) {
		logger.Debug("检查任务是否完成")
		f, err := service.TaskService().CheckTaskFinish(taskId)
		if err != nil {
			logger.Error(err)
			return
		}
		if f {
			if err = service.TaskService().FinishTask(taskId); err != nil {
				if !checkExist(t) {
					logger.Info("任务已完成：", taskId)
					Stop(taskId)
				} else {
					logger.Error("无法完成任务", taskId, "：", err)
				}
				return
			}
			Stop(taskId)
			t.Destroy()
			logger.Info("任务已完成：", taskId)
		}
	})
	logger.Debug("开启任务完成检测")
}
