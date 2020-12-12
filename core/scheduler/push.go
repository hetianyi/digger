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
	timer.Start(0, time.Second*5, 0, func(t *timer.Timer) {
		tasks, err := service.PushService().SelectPushTasks()
		if err != nil {
			logger.Error("无法获取推送任务数据: ", err)
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
		return
	}

	runningPushTask[task.TaskId] = true

	// 加载配置快照
	config, err := service.CacheService().GetSnapshotConfig(task.TaskId)
	if err != nil {
		logger.Error(fmt.Sprintf("无法为任务%d启动结果推送: %s", task.TaskId, err))
		return
	}
	if config == nil {
		logger.Error(fmt.Sprintf("无法为任务%d启动结果推送: 配置快照不存在", task.TaskId))
		if err := service.PushService().FinishPushTask(task.TaskId); err != nil {
			logger.Error(err)
		}
		return
	}

	logger.Info(fmt.Sprintf("为任务%d启动结果推送", task.TaskId))

	var lastResultId int64
	enableRetry := config.PushSources[0].EnableRetry
	pushSize := config.PushSources[0].PushSize
	if pushSize <= 0 {
		pushSize = 50
	}

	timer.Start(0, time.Second*3, 0, func(t *timer.Timer) {
		var data []interface{}
		retryTimes := 0
		var results []*models.Result
		for true {
			if data == nil {
				results, err = service.PushService().SelectPushResults(task.TaskId, pushSize, lastResultId)
				if err != nil {
					logger.Error(fmt.Sprintf("无法获取任务%d的结果: %s", task.TaskId, err))
					break
				}
				data = buildResultArray(results)
				if len(results) > 0 {
					lastResultId = results[len(results)-1].Id
				} else {
					logger.Info(fmt.Sprintf("任务%d结果推送结束", task.TaskId))
					t.Destroy()
					if err = service.PushService().FinishPushTask(task.TaskId); err != nil {
						logger.Error(fmt.Sprintf("无法完成任务%d的结果推送: %s", task.TaskId, err))
					}
					break
				}
			}
			// retry too many times, skip
			if retryTimes > 5 { // TODO
				retryTimes = 0
				data = nil
				logger.Info(fmt.Sprintf("任务%d该批结果数据推送失败: 达到最大重试次数", task.TaskId))
				if err = service.PushService().UpdatePushTaskResultId(task.TaskId, lastResultId); err != nil {
					logger.Error(fmt.Sprintf("无法更新任务%d的更新推送进度: %s", task.TaskId, err))
				}
				continue
			}
			if err = push(config.PushSources[0], data); err != nil {
				logger.Error(fmt.Sprintf("任务%d的结果推送失败: %s", task.TaskId, err))
				if enableRetry {
					retryTimes++
					time.Sleep(time.Second * 3)
					continue
				}
			}
			// logger.Info("推送成功")
			if len(results) < pushSize {
				logger.Info(fmt.Sprintf("任务%d结果推送完成", task.TaskId))
				t.Destroy()
				if err = service.PushService().FinishPushTask(task.TaskId); err != nil {
					logger.Error(fmt.Sprintf("无法完成任务%d的结果推送: %s", task.TaskId, err))
				}
				break
			}
			if err = service.PushService().UpdatePushTaskResultId(task.TaskId, lastResultId); err != nil {
				logger.Error(fmt.Sprintf("无法更新任务%d的更新推送进度: %s", task.TaskId, err))
				retryTimes++
				continue
			}
			retryTimes = 0
			data = nil
		}
	})
}

func buildResultArray(results []*models.Result) []interface{} {
	if len(results) == 0 {
		return nil
	}
	var resultArray []interface{}
	for _, v := range results {
		var temp interface{}
		json.Unmarshal([]byte(v.Result), &temp)
		resultArray = append(resultArray, temp)
	}
	//indent, _ := json.MarshalIndent(resultArray, "", " ")
	//fmt.Println(string(indent))
	return resultArray
}

func push(source *models.PushSource, data []interface{}) error {
	res, err := httpClient.R().
		SetHeaders(map[string]string{
			"User-Agent": "digger-push-client",
		}).
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
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
