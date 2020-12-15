///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package crontask

import (
	"digger/models"
	"digger/scheduler"
	"digger/services/service"
	"fmt"
	"github.com/hetianyi/gox/logger"
	"github.com/hetianyi/gox/timer"
	"github.com/robfig/cron/v3"
	"sync"
	"time"
)

var (
	cronExecutor   *cron.Cron
	cronProjectMap = make(map[int]*cron.EntryID)
	cronLock       = new(sync.Mutex)
)

func init() {
	// Seconds field, required
	cron.New(cron.WithSeconds())
	// Seconds field, optional
	cronExecutor = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	cronExecutor.Start()
}

func StartCronScanner() {
	timer.Start(0, time.Second*15, 0, func(t *timer.Timer) {
		projects, err := service.ProjectService().SelectCronProjectList()
		if err != nil {
			logger.Error(err)
			return
		}
		removingNotLongerExistsProjects(projects)
		for _, p := range projects {
			cronTask(p)
		}
	})
}

func removingNotLongerExistsProjects(projects []*models.Project) {
	cronLock.Lock()
	defer cronLock.Unlock()
	var excludes []int
	for k := range cronProjectMap {
		es := false
		for _, p := range projects {
			if p.Id == k {
				es = true
				break
			}
		}
		if !es {
			excludes = append(excludes, k)
		}
	}
	for _, v := range excludes {
		//logger.Info("no longer schedule task for project ", v)
		cronExecutor.Remove(*cronProjectMap[v])
		delete(cronProjectMap, v)
	}
}

func RemoveCron(projectId int) {
	logger.Info(fmt.Sprintf("项目%d定时任务移除", projectId))
	cronLock.Lock()
	defer cronLock.Unlock()

	if cronProjectMap[projectId] != nil {
		//logger.Info("no longer schedule task for project ", projectId)
		cronExecutor.Remove(*cronProjectMap[projectId])
		delete(cronProjectMap, projectId)
	}
}

func cronTask(project *models.Project) {
	if project == nil {
		return
	}
	cronLock.Lock()
	defer cronLock.Unlock()

	if cronProjectMap[project.Id] != nil {
		return
	}

	logger.Info(fmt.Sprintf("项目%s开启定时任务：%s", project.Name, project.Cron))

	id, err := cronExecutor.AddFunc(project.Cron, func() {
		scheduleStartTask(project)
	})
	if err != nil {
		logger.Error(fmt.Sprintf("cannot schedule timer task for project %s: %s", project.Name, err.Error()))
		return
	}
	cronProjectMap[project.Id] = &id
}

func scheduleStartTask(project *models.Project) {

	logger.Info(fmt.Sprintf("项目%s启动新任务", project.Name))

	project, err := service.ProjectService().SelectFullProjectInfo(project.Id)
	if err != nil {
		logger.Error(fmt.Sprintf("项目%s无法启动任务: 查询项目配置错误: %s", project.Name, err.Error()))
		return
	}
	if project == nil {
		logger.Error(fmt.Sprintf("项目%s无法启动任务: 找不到项目配置", project.Name))
		return
	}

	_, err = project.Validate()
	if err != nil {
		logger.Error(fmt.Sprintf("项目%s无法启动任务: 配置验证失败", project.Name))
		return
	}

	task, err := service.TaskService().CreateTask(models.Task{
		ProjectId:  project.Id,
		Status:     1,
		CreateTime: time.Now(),
	})

	if err != nil {
		logger.Error(fmt.Sprintf("项目%s任务启动错误: %s", project.Name, err.Error()))
		return
	}

	err = scheduler.Schedule(task)
	if err != nil {
		logger.Error(fmt.Sprintf("项目%s任务启动错误: %s", project.Name, err.Error()))
		return
	}
	logger.Info(fmt.Sprintf("项目%s任务启动成功", project.Name))
}
