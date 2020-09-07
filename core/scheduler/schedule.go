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
	scheduleLock     = new(sync.Mutex)
	blockQueueCache  = make(map[int]*queue.NoneBlockQueue)
	queueNotifier    = make(map[int]*queue.NoneBlockQueue)
	backPushNotifier = make(map[int]*queue.NoneBlockQueue)
	queryLock        = make(map[int]*sync.Mutex)
	taskFinishTimer  = make(map[int]*timer.Timer)
	stopFunc         StopFunc
)

func RegisterStopFunc(_stopFunc StopFunc) {
	stopFunc = _stopFunc
}

func Stop(taskId int) {
	stopFunc.Stop(taskId)
}

func StartScheduler() {
	scheduleScanTask()
	scheduleResetQueue()
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
			Schedule(task)
		}
	})
}

func scheduleResetQueue() {
	timer.Start(0, time.Second*10, 0, func(t *timer.Timer) {
		err := service.QueueService().ResetQueuesStatus()
		if err != nil {
			logger.Error("cannot reset queues: ", err)
		}
	})
}

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

func RegisterNotifier(taskId int) *queue.NoneBlockQueue {
	scheduleLock.Lock()
	defer scheduleLock.Unlock()

	if queueNotifier[taskId] == nil {
		queueNotifier[taskId] = queue.NewNoneBlockQueue(1)
	}
	return queueNotifier[taskId]
}

func DeRegisterNotifier(taskId int) {
	scheduleLock.Lock()
	defer scheduleLock.Unlock()

	delete(queueNotifier, taskId)
}

func BackPushNotify(taskId int) {
	if backPushNotifier[taskId] != nil {
		backPushNotifier[taskId].Put(1)
	}
}

func backPushListener(taskId, fetchSize int) {
	for {
		time.Sleep(time.Millisecond * 100)
		blockQueue := backPushNotifier[taskId]
		if blockQueue == nil {
			break
		}
		fetch, s := blockQueue.Fetch()
		if !s || fetch == nil {
			continue
		}
		selectTask, _ := service.TaskService().SelectTask(taskId) //service.TaskService().SelectTask(task.Id)
		if selectTask == nil || selectTask.Status != 1 {
			continue
		}
		//logger.Info("背压")
		fetchQueue(taskId, fetchSize)
	}
}

// 为一个task启动一个调度器，
// 调度器的任务是从数据库t_queue表轮询查询队列，并分配任务给worker节点
func Schedule(task *models.Task) error {
	scheduleLock.Lock()
	defer scheduleLock.Unlock()

	if blockQueueCache[task.Id] != nil {
		return nil
	}

	project, err := service.CacheService().GetSnapshotConfig(task.Id)
	if err != nil {
		return err
	}

	conSize := project.GetIntSetting(common.SETTINGS_CONCURRENT_REQUESTS)
	if conSize == 0 {
		conSize = 5
	}
	blockQueueCache[task.Id] = queue.NewNoneBlockQueue(conSize * 5)
	queryLock[task.Id] = new(sync.Mutex)
	backPushNotifier[task.Id] = queue.NewNoneBlockQueue(1)
	checkStatusFinishStatus(task.Id)
	go backPushListener(task.Id, conSize)

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
				Stop(task.Id)
				releaseTaskQueueCache(blockQueueCache[task.Id])
				delete(blockQueueCache, task.Id)
				delete(backPushNotifier, task.Id)
				delete(queueNotifier, task.Id)
				delete(queryLock, task.Id)
				delete(taskFinishTimer, task.Id)
				scheduleLock.Unlock()
				t.Destroy()
				return
			}
			if selectTask.Status == 1 {
				if fetchQueue(task.Id, conSize) {
					continue
				}
				// 没有更多，睡一会
			} else if selectTask.Status == 2 || selectTask.Status == 3 {
				scheduleLock.Lock()
				if selectTask.Status == 2 {
					logger.Info("停止任务：", task.Id)
				} else {
					logger.Info("任务已完成：", task.Id)
				}
				Stop(task.Id)
				releaseTaskQueueCache(blockQueueCache[task.Id])
				delete(blockQueueCache, task.Id)
				delete(backPushNotifier, task.Id)
				delete(queueNotifier, task.Id)
				delete(queryLock, task.Id)
				delete(taskFinishTimer, task.Id)
				scheduleLock.Unlock()
				t.Destroy()
			}
			break
		}
	})

	return nil
}

// 清空task的block queue
func releaseTaskQueueCache(blockQueue *queue.NoneBlockQueue) {
	if blockQueue != nil {
		for {
			i, s := blockQueue.Fetch()
			if s {
				logger.Info("释放queue：", i)
				continue
			}
			break
		}
	}
}

func fetchQueue(taskId int, fetchSize int) bool {
	lock := queryLock[taskId]
	if lock != nil {
		queryLock[taskId].Lock()
		defer queryLock[taskId].Unlock()
	} else {
		return false
	}

	var queues []*models.Queue
	queues, err := service.QueueService().SelectQueues(models.QueueQueryVO{
		TaskId:     taskId,
		LockStatus: true,
		Status:     0,
		Limit:      fetchSize,
	})
	if err != nil {
		logger.Debug("error query queue from database: ", err)
		return false
	}

	if len(queues) > 0 {
		logger.Info("加载了", len(queues), "条queue")
	}
	var ids []interface{}
	for _, v := range queues {
		ids = append(ids, v.Id)
	}

	doQueue(queues)
	if len(queues) > 0 {
		return true
	}
	return false
}

// 将查询的queue列表放入队列，如果队列已满，则会阻塞等待
func doQueue(queues []*models.Queue) {
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
		if task.Status == 2 || task.Status == 3 {
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
