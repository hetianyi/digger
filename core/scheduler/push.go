///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package scheduler

import (
	"digger/models"
	"digger/services/service"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/timer"
	"net/http"
	"sync"
	"time"
)

var (
	lock            = new(sync.Mutex)
	runningPushTask = make(map[int]bool)
	httpClient      = resty.New()
)

// 定时扫描任务，适用于程序启动时恢复task状态
func scheduleScanPushTask() {
	timer.Start(0, time.Second*30, 0, func(t *timer.Timer) {
		tasks, err := service.PushService().SelectPushTasks()
		if err != nil {
			logger.Error("cannot get active tasks: ", err)
		}
		if len(tasks) == 0 {
			return
		}
		for _, task := range tasks {
			SchedulePush(task)
		}
	})
}

func SchedulePush(task *models.PushTask) {
	lock.Lock()
	defer lock.Unlock()

	if runningPushTask[task.TaskId] {
		logger.Warn("push task is currently running, skip.")
		return
	}

	runningPushTask[task.TaskId] = true

	// 加载配置快照
	config, err := service.CacheService().GetSnapshotConfig(task.TaskId)
	if err != nil {
		logger.Error("cannot start push task: ", err)
		return
	}
	if config == nil {
		logger.Error("cannot start push task: snapshot config not found")
		if err := service.PushService().FinishPushTask(task.TaskId); err != nil {
			logger.Error(err)
		}
		return
	}

	logger.Info("run push task: ", task.TaskId)

	var lastResultId int64
	enableRetry := config.PushSources[0].EnableRetry
	pushSize := config.PushSources[0].PushSize
	if pushSize <= 0 {
		pushSize = 50
	}

	timer.Start(0, time.Second*3, 0, func(t *timer.Timer) {
		data := ""
		retryTimes := 0
		var results []*models.Result
		for true {
			if data == "" {
				results, err = service.PushService().SelectPushResults(task.TaskId, pushSize, lastResultId)
				if err != nil {
					logger.Error("cannot get push result: ", err)
					break
				}
				data = buildResultArray(results)
				if len(results) > 0 {
					lastResultId = results[len(results)-1].Id
				} else {
					logger.Info("push task finish: ", task.TaskId)
					t.Destroy()
					if err = service.PushService().FinishPushTask(task.TaskId); err != nil {
						logger.Error("cannot finish push task: ", err)
					}
					break
				}
			}
			// retry too many times, skip
			if retryTimes > 5 { // TODO
				retryTimes = 0
				data = ""
				logger.Info("cannot push: reach max retry times: ", task.TaskId)
				if err = service.PushService().UpdatePushTaskResultId(task.TaskId, lastResultId); err != nil {
					logger.Error("[1]error update push state: ", err)
				}
				continue
			}
			if err = push(config.PushSources[0], data); err != nil {
				logger.Error("error push result: ", err)
				if enableRetry {
					retryTimes++
					time.Sleep(time.Second * 3)
					continue
				}
			}
			logger.Info("push success: ", task.TaskId)
			if len(results) < pushSize {
				logger.Info("push task finish: ", task.TaskId)
				t.Destroy()
				if err = service.PushService().FinishPushTask(task.TaskId); err != nil {
					logger.Error("cannot finish push task: ", err)
				}
				break
			}
			if err = service.PushService().UpdatePushTaskResultId(task.TaskId, lastResultId); err != nil {
				logger.Error("[2]error update push state: ", err)
				retryTimes++
				continue
			}
			retryTimes = 0
			data = ""
		}
	})
}

func buildResultArray(results []*models.Result) string {
	if len(results) == 0 {
		return "[]"
	}
	var resultArray []interface{}
	for _, v := range results {
		var temp interface{}
		json.Unmarshal([]byte(v.Result), &temp)
		resultArray = append(resultArray, temp)
	}
	r, _ := json.Marshal(resultArray)
	return string(r)
}

func push(source *models.PushSource, data string) error {
	res, err := httpClient.R().
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"User-Agent":   "digger-push-client",
		}).
		SetBody(data).
		Execute(source.Method, source.Url)
	if err != nil {
		return err
	}
	if res.StatusCode() != http.StatusOK {
		return errors.New(fmt.Sprintf("http status: %d", res.StatusCode()))
	}
	return nil
}
