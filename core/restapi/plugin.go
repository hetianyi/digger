///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"digger/models"
	"digger/services/service"
	"github.com/gin-gonic/gin"
	"github.com/hetianyi/gox/convert"
	"net/http"
)

type SavePluginRequestVO struct {
	ProjectId int              `json:"projectId"`
	Plugins   []*models.Plugin `json:"plugins"`
}

// 为项目开始一个新任务
func SavePlugins(c *gin.Context) {

	id := c.Param("id")
	pid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	// 绑定请求数据
	reqData := &SavePluginRequestVO{}
	if err := c.ShouldBindJSON(reqData); err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	project, err := service.ProjectService().SelectProjectById(pid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	if project == nil {
		c.JSON(http.StatusOK, ErrorMsg("project not found"))
		return
	}

	err = service.PluginService().SavePlugins(pid, reqData.Plugins)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(nil))
}
