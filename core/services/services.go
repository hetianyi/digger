///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package services

import (
	"digger/models"
	"time"
)

// 项目服务
type ProjectService interface {
	// 根据项目id查询项目信息
	SelectProjectById(projectId int) (*models.Project, error)
	// 根据项目id查询项目完整信息
	SelectFullProjectInfo(projectId int) (*models.Project, error)
	// 根据项目名称查询项目信息
	SelectProjectByName(name string) (*models.Project, error)
	// 根据条件查询项目列表
	SelectProjectList(params models.ProjectQueryVO) (int64, []*models.Project, error)
	// 新增项目
	CreateProject(project models.Project) (*models.Project, error)
	// 修改项目信息
	UpdateProject(project models.Project) (bool, error)
	// 修改项目代理服务器
	UpdateProjectProxies(projectId int, proxyIds []int) error
	// 根据项目名称查询项目信息
	DeleteProject(projectId int) (bool, error)
	// 查询启用定时任务的项目
	SelectCronProjectList() ([]*models.Project, error)
	// 查询所有项目数量
	AllProjectCount() (int, error)
}

// 项目配置服务
type ProjectConfigService interface {
	// 根据项目id查询阶段列表
	SelectStages(projectId int) ([]models.Stage, error)
	// 根据项目id查询参数列表
	SelectFields(projectId int) ([]models.Field, error)
	// 保存项目字段数据
	SaveStagesAndFields(projectId int, stages []models.Stage) error
	// 从配置文件解析的数据保存
	SaveProjectConfig(project *models.Project, stages []models.Stage) error
	// 导入项目配置
	ImportProjectConfig(project *models.Project) error
}

// 配置快照服务接口
type TaskService interface {
	// 创建任务
	CreateTask(task models.Task) (*models.Task, error)
	// 查询任务详情
	SelectTask(id int) (*models.Task, error)
	// 查询任务列表
	SelectTaskList(params models.TaskQueryVO) (int64, []*models.Task, error)
	// 查询所有已激活的task
	SelectActiveTasks() ([]*models.Task, error)
	// 加载任务的配置快照
	LoadConfigSnapshot(snapshotId int) (*models.Project, error)
	// 完成任务
	FinishTask(id int) error
	// 关闭任务
	ShutdownTask(id int) error
	// 暂停任务
	PauseTask(id int) error
	// 开启任务
	StartTask(id int) error
	// 查询任务数量
	TaskCount(projectIds ...int) ([]*models.TaskCountCO, error)
	// 查询所有任务数量
	AllTaskCount() (int, error)
	// 检查task是否完成
	CheckTaskFinish(taskId int) (bool, error)
	// 删除task
	DeleteTask(taskId int) error
}

// 配置快照服务接口
type DBService interface {
	Check() error
}

// 配置快照服务接口
type SnapshotService interface {
	// 根据项目id查询阶段列表
	SelectSnapshot(id int) (models.ConfigSnapshot, error)
}

// 结果服务接口
type ResultService interface {
	// 查询任务的结果列表
	SelectResults(params models.ResultQueryVO) (int64, []*models.Result, error)
	// 导出结果
	ExportResults(params models.ResultQueryVO) ([]*models.Result, error)
	// 插入结果
	InsertResults(taskId int, results []models.Result) error
	// 保存一个check内的数据 v1
	SaveCheckData(checkData *models.QueueCallbackRequestVO, exceedMaxRetryQueueIds []int64) error
	// 保存处理结果的数据 v2
	SaveProcessResultData(result *models.QueueProcessResult, exceedMaxRetry bool) error
	// 查询任务成功结果数量
	ResultCount(taskId ...int) ([]*models.ResultCountCO, error)
	// 查询从某个id起后续结果总数
	ResultCountSince(id int64) (int, int64, error)
}

// 调度任务服务
type QueueService interface {
	// 查询队列爬虫任务
	SelectQueues(params models.QueueQueryVO) ([]*models.Queue, error)
	// 插入队列
	InsertQueue(queue models.Queue) error
	// 重置处理中的queue状态为未处理
	ResetQueuesStatus() error
	// 查询任务未完成的queue数量
	GetUnFinishedCount(taskId int) (int, error)
	// 根据任务id删除queue
	DeleteQueues(taskId int) error
	// 查询任务失败queue数量
	ErrorCount(taskId ...int) ([]*models.ResultCountCO, error)
	// 结束时统计最终错误数
	StatisticFinal(taskId int) error
}

// 调度任务服务
type PluginService interface {
	// 根据名称查询插件
	SelectPluginByName(name string) (*models.Plugin, error)
	// 根据项目查询插件列表
	SelectPluginsByProject(projectId int) ([]models.Plugin, error)
	// 新增插件
	InsertPlugin(plugin models.Plugin) error
	// 批量保存插件
	SavePlugins(projectId int, plugins []*models.Plugin) error
}

// 配置服务
type ConfigService interface {
	ListConfigs() (map[string]string, error)
	UpdateConfig(key, value string) error
}

// 配置服务
type StatisticService interface {
	Save(data map[string]interface{}) error
	List(start time.Time, end time.Time) ([]*models.StatisticVO, error)
}

// 代理配置服务
type ProxyService interface {
	Save(proxy models.Proxy) error
	Delete(idList []int) error
	List(params *models.ProxyQueryVO) (int64, []*models.Proxy, error)
	SelectByProject(projectId int) ([]*models.Proxy, error)
}
