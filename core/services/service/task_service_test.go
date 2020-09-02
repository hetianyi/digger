///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service_test

import (
	"digger/models"
	"digger/services/service"
	"encoding/json"
	"fmt"
	"github.com/hetianyi/gox/logger"
	"testing"
	"time"
)

func TestTaskServiceImp_CreateTask(t *testing.T) {

	service.TaskService().CreateTask(models.Task{
		ProjectId:      18,
		Status:         0,
		ResultCount:    0,
		IOIn:           0,
		IOOut:          0,
		SuccessRequest: 0,
		ErrorRequest:   0,
		BindNodeMode:   0,
		CreateTime:     time.Now(),
	})
}

func TestTaskServiceImp_SelectTask(t *testing.T) {
	task, _ := service.TaskService().SelectTask(1)
	fmt.Println(task)
}

func TestTaskServiceImp_SelectTaskList(t *testing.T) {
	total, list, _ := service.TaskService().SelectTaskList(models.TaskQueryVO{
		ProjectId: 0,
		Status:    -1,
		PageQueryVO: models.PageQueryVO{
			Page:     1,
			PageSize: 10,
		},
	})

	fmt.Println(len(list))
	logger.Info(fmt.Sprintf("总数：%d", total))
}

func TestTaskServiceImp_LoadConfigSnapshot(t *testing.T) {
	p, _ := service.TaskService().LoadConfigSnapshot(3)
	fmt.Println(p)
}

func TestTaskServiceImp_PauseTask(t *testing.T) {
	fmt.Println(service.TaskService().PauseTask(1))
}

func TestTaskServiceImp_StartTask(t *testing.T) {
	fmt.Println(service.TaskService().StartTask(1))
}

func TestTaskServiceImp_ShutdownTask(t *testing.T) {
	fmt.Println(service.TaskService().ShutdownTask(1))
}

func TestTaskServiceImp_TaskCount(t *testing.T) {
	cos, _ := service.TaskService().TaskCount(1, 2)
	bytes, _ := json.Marshal(cos)
	fmt.Println(string(bytes))
}

func TestTaskServiceImp_AllTaskCount(t *testing.T) {
	fmt.Println(service.TaskService().AllTaskCount())
}