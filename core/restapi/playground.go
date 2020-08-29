///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"digger/crawler"
	"digger/models"
	"digger/services/service"
	"digger/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

// 为项目开始一个新任务
func PlayExistStage(c *gin.Context) {
	var reqBody models.PlayInputVO1
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	project, err := service.ProjectService().SelectFullProjectInfo(reqBody.ProjectId)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	if project == nil {
		c.JSON(http.StatusOK, ErrorMsg("project not found"))
	}

	stage := project.GetStageByName(reqBody.StageName)
	if stage == nil {
		c.JSON(http.StatusOK, ErrorMsg("stage not found"))
	}

	err = crawler.Play(&models.Queue{
		Id:         0,
		TaskId:     0,
		StageName:  reqBody.StageName,
		Url:        reqBody.Url,
		MiddleData: "",
	}, project, os.Stdout, func(oldQueue *models.Queue, newQueue []*models.Queue, results []*models.Result, err error) {
		if err != nil {
			c.JSON(http.StatusOK, ErrorMsg(err.Error()))
			return
		}
		ret := &models.PlayOutputVO{
			ProjectId: reqBody.ProjectId,
			Url:       reqBody.Url,
			StageName: reqBody.StageName,
			Next:      newQueue,
			Result:    results,
		}
		c.JSON(http.StatusOK, Success(ret))
		return
	})
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
	}
}

// {
//	"stage_name":"detail",
//	"url":"https://shumeipai.nxez.com/2020/07/02/rpi-fan-on-sale.html",
//	"project":"name: shumeipai_labs\ndisplay_name: 树莓派实验室\nremark: 树莓派实验室\nheaders:\n  User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML,like\n    Gecko) Chrome/78.0.3904.108 Safari/537.36\nsettings:\n  CONCURRENT_REQUESTS: \"5\"\n  DEFAULT_REQUEST_HEADERS: \"\"\n  DOWNLOAD_DELAY: \"\"\ntags: \"\"\nstart_url: https://shumeipai.nxez.com/\nstart_stage: page_list\nstages:\n- name: page_list\n  is_list: true\n  list_css: h3.entry-title>a\n  page_css: a.next\n  page_attr: href\n  fields:\n  - name: detail_link\n    is_array: false\n    is_html: false\n    css: \"\"\n    attr: href\n    remark: \"\"\n    next_stage: detail\n- name: detail\n  is_list: false\n  list_css: div#search_nature_rg>ul>li\n  page_css: li.next>a\n  page_attr: href\n  fields:\n  - name: title\n    is_array: false\n    is_html: false\n    css: h1.entry-title\n    attr: \"\"\n    remark: \"\"\n    next_stage: \"\"\n  - name: date\n    is_array: false\n    is_html: false\n    css: span.entry-meta-date>a\n    attr: \"\"\n    remark: \"\"\n    next_stage: \"\"\n  - name: author\n    is_array: false\n    is_html: false\n    css: span.entry-meta-author>a\n    attr: \"\"\n    remark: \"\"\n    next_stage: \"\"\n  - name: tags\n    is_array: true\n    is_html: false\n    css: div.entry-tags>ul>li>a\n    attr: \"\"\n    remark: \"\"\n    next_stage: \"\""
// }
func PlayFromTempStage(c *gin.Context) {
	var reqBody models.PlayInputVO2
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
		return
	}

	var project models.Project
	err := utils.ParseYamlFromString(reqBody.Project, &project)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	for i, s := range project.Stages {
		var fs []models.Field
		hasNextStage := false
		for _, f := range s.Fields {
			if f.StageId == s.Id {
				fs = append(fs, f)
				if f.NextStage != "" {
					hasNextStage = true
				}
			}
		}
		project.Stages[i].Fields = fs
		project.Stages[i].HasNextStage = hasNextStage
	}

	stage := project.GetStageByName(reqBody.StageName)
	if stage == nil {
		c.JSON(http.StatusOK, ErrorMsg("stage not found"))
		return
	}

	plugins, err := service.PluginService().SelectPluginsByProject(reqBody.ProjectId)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	service.AssembleProjectPlugin(&project, plugins)

	write := false
	err = crawler.Play(&models.Queue{
		Id:         0,
		TaskId:     0,
		StageName:  reqBody.StageName,
		Url:        reqBody.Url,
		MiddleData: "",
	}, &project, os.Stdout, func(oldQueue *models.Queue, newQueue []*models.Queue, results []*models.Result, err error) {
		if err != nil {
			write = true
			c.JSON(http.StatusOK, ErrorMsg(err.Error()))
			return
		}
		ret := &models.PlayOutputVO{
			ProjectId: 0,
			Url:       reqBody.Url,
			StageName: reqBody.StageName,
			Next:      newQueue,
			Result:    results,
		}
		c.JSON(http.StatusOK, Success(ret))
		write = true
		return
	})
	if err != nil && !write {
		c.JSON(http.StatusOK, Success(err.Error()))
	}
}

// {
//	"stage_name":"detail",
//	"url":"https://shumeipai.nxez.com/2020/07/02/rpi-fan-on-sale.html",
//	"project":"name: shumeipai_labs\ndisplay_name: 树莓派实验室\nremark: 树莓派实验室\nheaders:\n  User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML,like\n    Gecko) Chrome/78.0.3904.108 Safari/537.36\nsettings:\n  CONCURRENT_REQUESTS: \"5\"\n  DEFAULT_REQUEST_HEADERS: \"\"\n  DOWNLOAD_DELAY: \"\"\ntags: \"\"\nstart_url: https://shumeipai.nxez.com/\nstart_stage: page_list\nstages:\n- name: page_list\n  is_list: true\n  list_css: h3.entry-title>a\n  page_css: a.next\n  page_attr: href\n  fields:\n  - name: detail_link\n    is_array: false\n    is_html: false\n    css: \"\"\n    attr: href\n    remark: \"\"\n    next_stage: detail\n- name: detail\n  is_list: false\n  list_css: div#search_nature_rg>ul>li\n  page_css: li.next>a\n  page_attr: href\n  fields:\n  - name: title\n    is_array: false\n    is_html: false\n    css: h1.entry-title\n    attr: \"\"\n    remark: \"\"\n    next_stage: \"\"\n  - name: date\n    is_array: false\n    is_html: false\n    css: span.entry-meta-date>a\n    attr: \"\"\n    remark: \"\"\n    next_stage: \"\"\n  - name: author\n    is_array: false\n    is_html: false\n    css: span.entry-meta-author>a\n    attr: \"\"\n    remark: \"\"\n    next_stage: \"\"\n  - name: tags\n    is_array: true\n    is_html: false\n    css: div.entry-tags>ul>li>a\n    attr: \"\"\n    remark: \"\"\n    next_stage: \"\""
// }
func ParseConfigFile(c *gin.Context) {
	var reqBody models.PlayInputVO2
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
		return
	}

	var project models.Project
	err := utils.ParseYamlFromString(reqBody.Project, &project)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(&project))
}
