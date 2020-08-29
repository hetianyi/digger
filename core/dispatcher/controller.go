///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package dispatcher

import (
	"digger/common"
	"digger/scheduler"
	"github.com/hetianyi/gox/logger"
)

type stopFuncImp struct {
}

var (
	taskSig = make(map[int]int)
)

func init() {
	scheduler.RegisterStopFunc(&stopFuncImp{})
}

// 暂停任务
func Pause(taskId int) {
	dispatchLock.Lock()
	defer dispatchLock.Unlock()
	logger.Info("暂停任务：", taskId)

	if taskScheduleTimer[taskId] == nil || taskSig[taskId] == common.PAUSE {
		return
	}
	if taskSig[taskId] == common.RUNNING {
		taskSig[taskId] = common.PAUSE
		return
	}
}

// 停止任务
func (*stopFuncImp) Stop(taskId int) {
	dispatchLock.Lock()
	defer dispatchLock.Unlock()

	if taskWorkLock[taskId] != nil {
		taskWorkLock[taskId].Lock()
		defer taskWorkLock[taskId].Unlock()
	}
	logger.Info("停止任务：", taskId)

	if taskScheduleTimer[taskId] != nil {
		taskScheduleTimer[taskId].Destroy()
	}
	delete(taskSig, taskId)
	delete(taskWorkLock, taskId)
	delete(notifiers, taskId)
	file := logFileMap[taskId]
	if file != nil {
		file.Close()
	}
	delete(logFileMap, taskId)
}

// 继续暂停的任务
func Continue(taskId int) {
	dispatchLock.Lock()
	defer dispatchLock.Unlock()
	logger.Info("继续任务：", taskId)

	if taskScheduleTimer[taskId] == nil || taskSig[taskId] == common.RUNNING {
		return
	}
	taskSig[taskId] = common.RUNNING
}

func ExistClient(id int) bool {
	dispatchLock.Lock()
	defer dispatchLock.Unlock()

	client := selectClientById(id)

	if client == nil {
		return false
	}
	return true
}
