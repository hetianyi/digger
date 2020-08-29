///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service

import (
	"digger/models"
	"errors"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"sync"
	"time"
)

const (
	cache_key_task                 = "TASK_CONFIG_SNAPSHOT:"
	cache_key_task_config_snapshot = "TASK_CONFIG_SNAPSHOT:"
)

// 项目服务
type cacheServiceImp struct {
	cache     map[string]interface{}
	cacheLock *sync.Mutex
}

// 缓存任务和配置快照信息
func (c cacheServiceImp) CacheFullProjectInfo(task *models.Task) (*models.Project, error) {
	project, err := TaskService().LoadConfigSnapshot(task.ConfigSnapShotId)
	if err != nil {
		return nil, err
	}
	if project == nil {
		return nil, errors.New("no snapshot")
	}
	c.cache[cache_key_task_config_snapshot+convert.IntToStr(task.Id)] = project
	return project, nil
}

// 根据id获取当前任务详情
func (c cacheServiceImp) GetTask(taskId int) (*models.Task, error) {
	task := c.cache[cache_key_task+convert.IntToStr(taskId)]
	if task != nil {
		return task.(*models.Task), nil
	}
	selectTask, err := TaskService().SelectTask(taskId)
	if selectTask != nil {
		c.cache[cache_key_task+convert.IntToStr(taskId)] = selectTask
	}
	return selectTask, err
}

// 根据任务获取当前项目的配置快照
func (c cacheServiceImp) GetSnapshotConfig(taskId int) (*models.Project, error) {
	c.cacheLock.Lock()
	defer c.cacheLock.Unlock()
	project := c.cache[cache_key_task_config_snapshot+convert.IntToStr(taskId)]
	if project != nil {
		return project.(*models.Project), nil
	}
	task, err := TaskService().SelectTask(taskId)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, nil
	}
	return c.CacheFullProjectInfo(task)
}

// 将queue错误次数+1
func (c cacheServiceImp) IncreQueueErrorCount(queueTaskIds []int, queueIds []int64) ([]int64, error) {

	var errCacheQueueTaskIds []string
	var errCacheQueueIds []interface{}
	for i, v := range queueIds {
		errCacheQueueTaskIds = append(errCacheQueueTaskIds, fmt.Sprintf("ERR_QUEUE:%d", queueTaskIds[i]))
		errCacheQueueIds = append(errCacheQueueIds, v)
	}

	ret, err := RedisClient.Eval(`
local r = {}
for k,v in ipairs(ARGV) do
	r[k]=redis.call('HINCRBY', KEYS[k], v, 1)
end
return r
`, errCacheQueueTaskIds, errCacheQueueIds...).Result()

	if err != nil {
		return nil, err
	}

	result := make([]int64, len(queueIds))

	arr := ret.([]interface{})
	for i, v := range arr {
		result[i] = v.(int64)
	}
	return result, nil
}

// 批量获取hash值
func (c cacheServiceImp) ExistMembers(taskId int, members []interface{}) ([]bool, error) {
	ret, err := RedisClient.Eval(`
local r = {}
for k,v in ipairs(ARGV) do
	r[k]=redis.call('SISMEMBER', KEYS[1], v)
end
return r
`, []string{fmt.Sprintf("DONE_QUEUE:%d", taskId)}, members...).Result()

	if err != nil {
		return nil, err
	}

	result := make([]bool, len(members))

	arr := ret.([]interface{})
	for i, v := range arr {
		if v.(int64) == 1 {
			result[i] = true
		} else {
			result[i] = false
		}
	}
	return result, nil
}

// 分布式锁，保证一个task的下同时只能有一个manager查询queue，避免queue被多个manager同时加载
func (c cacheServiceImp) LockTaskQueueFetch(taskId int, job func()) bool {
	// Create a new lock client.
	locker := redislock.New(RedisClient)

	var err error
	var lock *redislock.Lock
	var retry = 0
	for {
		// Try to obtain lock.
		lock, err = locker.Obtain(fmt.Sprintf("FETCH_QUEUE:%d", taskId), 30*time.Second, nil)
		if err != nil {
			retry++
			time.Sleep(time.Millisecond * 100)
			// 2s内没获取到锁，则放弃
			if retry > 20 {
				logger.Error("timeout fetching fetch lock")
				return false
			}
		} else {
			break
		}
	}
	defer lock.Release()
	if job != nil {
		job()
	}
	return true
}

// 缓存已成功的queue
func (c cacheServiceImp) SaveSuccessQueueIds(reqBody *models.QueueCallbackRequestVO) {
	if len(reqBody.SuccessQueueIds) == 0 {
		return
	}

	keys := make([]string, len(reqBody.SuccessQueueIds))
	args := make([]string, len(reqBody.SuccessQueueIds))
	for i := range reqBody.SuccessQueueIds {
		keys[i] = fmt.Sprintf("DONE_QUEUE:%d", reqBody.SuccessQueueTaskIds[i])
		args[i] = convert.Int64ToStr(reqBody.SuccessQueueIds[i])
	}

	_, err := RedisClient.Eval(`
local r
for k,v in pairs(ARGV) do
	redis.call('SADD', KEYS[k], v)
end
return 'OK'
`, keys, args).Result()

	if err != nil {
		logger.Error("error cache done queues: ", err)
		return
	}
}

// 增加task的并发数
func (c cacheServiceImp) IncreConcurrentTaskCount(requestId string, taskId, concurrent int) bool {
	ret, err := RedisClient.Eval(`
local targets = redis.call('KEYS', KEYS[1]..'*')
redis.log(redis.LOG_NOTICE, 'concur: '..tostring(#targets))
if (#targets >= tonumber(KEYS[3]))
then
  return 0
end
redis.call('setex', KEYS[2], 10, '')
return 1
`, []string{fmt.Sprintf("CON_TASK:%d:", taskId), fmt.Sprintf("CON_TASK:%d:%s", taskId, requestId), convert.IntToStr(concurrent)}).Result()

	if err != nil {
		return false
	}

	if ret.(int64) == 0 {
		return false
	}
	return true
}

// 减少task的并发数
func (c cacheServiceImp) DecreConcurrentTaskCount(requestId string, taskId int) bool {
	_, err := RedisClient.Del(fmt.Sprintf("CON_TASK:%d:%s", taskId, requestId)).Result()
	if err != nil {
		return false
	}
	return true
}

// 检查资源是否已完成
func (c cacheServiceImp) IsUniqueResFinish(taskId int, res string) bool {
	b, err := RedisClient.SIsMember(fmt.Sprintf("FINISH_RES:%d", taskId), res).Result()
	if err != nil {
		return false
	}
	return b
}

// 添加已完成unique资源
func (c cacheServiceImp) AddFinishUniqueRes(taskId int, res string) error {
	_, err := RedisClient.SAdd(fmt.Sprintf("FINISH_RES:%d", taskId), res).Result()
	return err
}
