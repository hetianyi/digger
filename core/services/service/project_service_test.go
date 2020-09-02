///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service_test

import (
	"digger/models"
	"digger/services/service"
	"digger/utils"
	"fmt"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/logger"
	json "github.com/json-iterator/go"
	"testing"
)

func TestProjectServiceImp_SelectProjectById(t *testing.T) {
	project, err := service.ProjectService().SelectProjectById(11)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(json.MarshalToString(project))
}

func TestProjectServiceImp_SelectFullProjectInfo(t *testing.T) {
	project, err := service.ProjectService().SelectFullProjectInfo(67)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(json.MarshalToString(project))
	s, _ := utils.Convert2Yaml(project)
	fmt.Println(s)
}

func TestProjectServiceImp_SelectProjectList(t *testing.T) {
	total, list, err := service.ProjectService().SelectProjectList(models.ProjectQueryVO{
		PageQueryVO: models.PageQueryVO{
			Page:     1,
			PageSize: 10,
		},
	})
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(len(list))
	logger.Info(fmt.Sprintf("总数：%d", total))
}

func TestProjectServiceImp_InsertProject(t *testing.T) {
	project, err := service.ProjectService().CreateProject(models.Project{
		Id:          0,
		Name:        "GO222",
		DisplayName: "GOGOGO",
		Remark:      "阿萨德",
		Settings:    nil,
		Tags:        "是大三大四的",
		StartUrl:    "http://xxx",
		StartStage:  "StartStage",
		Stages:      nil,
	})
	if err != nil {
		logger.Fatal(err)
	}
	s, _ := json.MarshalToString(project)
	fmt.Println(s)
}

func TestProjectServiceImp_UpdateProject(t *testing.T) {
	success, err := service.ProjectService().UpdateProject(models.Project{
		Id:          11,
		Name:        "GO222",
		DisplayName: "----",
		Remark:      "----------",
		Settings:    nil,
		Tags:        "--------------------",
		Stages:      nil,
	})
	if err != nil {
		logger.Fatal(err)
	}
	if success {
		logger.Info("更新成功")
	}

	// 验证
	project, err := service.ProjectService().SelectProjectById(9)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(json.MarshalToString(project))
}

func TestProjectServiceImp_DeleteProject(t *testing.T) {
	success, err := service.ProjectService().DeleteProject(11)

	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(gox.TValue(success, "删除成功", "删除失败"))
}

func TestProjectServiceImp_SelectCronProjectList(t *testing.T) {
	projects, err := service.ProjectService().SelectCronProjectList()
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(projects)
}

func TestProjectServiceImp_AllProjectCount(t *testing.T) {
	fmt.Println(service.ProjectService().AllProjectCount())
}