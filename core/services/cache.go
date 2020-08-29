///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package services

import (
	"digger/models"
)

// 项目服务
type CacheService interface {
	// 缓存任务和配置快照信息
	CacheFullProjectInfo(task *models.Task) (*models.Project, error)
	// 根据id获取当前任务详情
	GetTask(taskId int) (*models.Task, error)
	// 根据任务获取当前项目的配置快照
	GetSnapshotConfig(taskId int) (*models.Project, error)
	// 将queue错误次数+1
	IncreQueueErrorCount(queueTaskIds []int, queueIds []int64) ([]int64, error)
	// 批量获取hash值
	ExistMembers(taskId int, members []interface{}) ([]bool, error)
	// 分布式锁，保证一个task的下同时只能有一个manager查询queue，避免queue被多个manager同时加载
	LockTaskQueueFetch(taskId int, job func()) bool
	// 缓存已成功的queue
	SaveSuccessQueueIds(reqBody *models.QueueCallbackRequestVO)
	// 增加task的并发数
	// 如果已经达到concurrent，则不进行自增，返回false
	IncreConcurrentTaskCount(requestId string, taskId, concurrent int) bool
	// 减少task的并发数
	DecreConcurrentTaskCount(requestId string, taskId int) bool
	// 检查资源是否已完成
	IsUniqueResFinish(taskId int, res string) bool
	// 添加已完成unique资源
	AddFinishUniqueRes(taskId int, res string) error
}
