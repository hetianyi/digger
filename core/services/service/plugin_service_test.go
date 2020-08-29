///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service_test

import (
	"digger/crawler"
	"digger/models"
	"digger/services/service"
	"fmt"
	"github.com/hetianyi/gox/logger"
	"os"
	"testing"
)

func TestPluginServiceImp_InsertPlugin(t *testing.T) {
	fmt.Println(service.PluginService().InsertPlugin(models.Plugin{
		Id:        0,
		Name:      "rest",
		Script:    "(func(){return 1;})()",
		Slot:      "s1",
		ProjectId: 1,
	}))
}

func TestPluginServiceImp_SelectPluginByName(t *testing.T) {
	fmt.Println(service.PluginService().SelectPluginByName("xxx"))
}

func TestPluginServiceImp_SelectPluginsByProject(t *testing.T) {
	fmt.Println(service.PluginService().SelectPluginsByProject(1))
}

func TestData(t *testing.T) {
	projectId := insertPlugins(buildRaspberryPiLabsProject())
	fmt.Println("projectId=", projectId)
}

func TestPlugin1(t *testing.T) {

	project, err := service.ProjectService().SelectFullProjectInfo(65)
	if err != nil {
		logger.Fatal(err)
	}

	crawler.Play(&models.Queue{
		Id:         0,
		TaskId:     8,
		StageName:  "page_list",
		Url:        "https://www.taptap.com/ajax/top/download?page=2&total=30",
		MiddleData: "",
	}, project, os.Stdout, func(oldQueue *models.Queue, newQueue []*models.Queue, results []*models.Result, err error) {
		fmt.Println()
	})
}

func TestPlugin2(t *testing.T) {

	project, err := service.ProjectService().SelectFullProjectInfo(65)
	if err != nil {
		logger.Fatal(err)
	}

	crawler.Play(&models.Queue{
		Id:         0,
		TaskId:     8,
		StageName:  "detail",
		Url:        "https://www.taptap.com/app/130651",
		MiddleData: "{\"detail_link\":\"https://www.taptap.com/app/130651\"}",
	}, project, os.Stdout, func(oldQueue *models.Queue, newQueue []*models.Queue, results []*models.Result, err error) {
		fmt.Println()
	})
}

func insertPlugins(projectId int) int {
	service.PluginService().InsertPlugin(models.Plugin{
		Id:   0,
		Name: "demo",
		Script: `
(function(){
	console.log(CRAW_DATA());
	var data = JSON.parse(CRAW_DATA()).data;
	console.log(data.next);
	ADD_QUEUE({
		stage: QUEUE().StageName,
		url: data.next,
	});
	return data.html;
})()
`,
		Slot:      "s2",
		ProjectId: projectId,
	})
	return projectId
}

func buildRaspberryPiLabsProject() int {

	project, err := service.ProjectService().CreateProject(models.Project{
		Id:          0,
		Name:        "taptap",
		DisplayName: "taptap",
		Remark:      "taptap",
		Settings: map[string]string{
			"CONCURRENT_REQUESTS":     "5",
			"DEFAULT_REQUEST_HEADERS": "",
			"DOWNLOAD_DELAY":          "",
		},
		Headers: map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML,like Gecko) Chrome/78.0.3904.108 Safari/537.36",
		},
		Tags:       "",
		StartUrl:   "https://www.taptap.com/ajax/top/download?page=1&total=30",
		StartStage: "page_list",
		Stages:     nil,
	})
	if err != nil {
		logger.Fatal(err)
	}

	stages := []models.Stage{
		{
			ProjectId: project.Id,
			Name:      "page_list",
			IsList:    true,
			ListCss:   "div.taptap-top-card",
			PageCss:   "",
			PageAttr:  "",
			Plugins:   nil,
			PluginsDB: "demo@s2",
			Fields: []models.Field{
				{
					ProjectId: project.Id,
					Name:      "detail_link",
					IsArray:   false,
					IsHtml:    false,
					Css:       "div.top-card-middle>a.card-middle-title",
					Attr:      "href",
					Plugin:    nil,
					PluginDB:  "",
					Remark:    "",
					NextStage: "detail",
				},
			},
		},
		{
			ProjectId: project.Id,
			Name:      "detail",
			IsList:    false,
			ListCss:   "",
			PageCss:   "",
			PageAttr:  "",
			Plugins:   nil,
			PluginsDB: "",
			Fields: []models.Field{
				{
					ProjectId: project.Id,
					Name:      "gameName",
					IsArray:   false,
					IsHtml:    false,
					Css:       "div.base-info-wrap>h1",
					Attr:      "",
					Plugin:    nil,
					PluginDB:  "",
					Remark:    "",
					NextStage: "",
				},
			},
		},
	}

	err = service.ProjectConfigService().SaveStagesAndFields(project.Id, stages)
	if err != nil {
		logger.Fatal(err)
	}

	return project.Id
}

func TestPluginServiceImp_SavePlugins(t *testing.T) {
	fmt.Println(service.PluginService().SavePlugins(1, []*models.Plugin{
		{
			Name:      "p1",
			Script:    "gggggggggg",
			Slot:      "gggg",
			ProjectId: 1,
		},
		{
			Name:      "p2",
			Script:    "hhhhhhhhhhhhhh",
			Slot:      "hhhh",
			ProjectId: 1,
		},
	}))
}
