///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"digger/services/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetProjectConfigSnapshot(c *gin.Context) {
	taskId := GetIntParameter(c, "taskId", 0)
	if taskId == 0 {
		c.JSON(http.StatusOK, ErrorMsg("invalid parameter"))
		return
	}
	project, err := service.CacheService().GetSnapshotConfig(taskId)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(project))
}
