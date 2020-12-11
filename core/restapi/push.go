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
	"github.com/hetianyi/gox/logger"
	"net/http"
	"strings"
)

type PushQueryResultVO struct {
	Page     int                  `json:"page"`
	PageSize int                  `json:"pageSize"`
	Total    int64                `json:"total"`
	Data     []*models.PushSource `json:"data"`
}

func QueryPushSource(c *gin.Context) {
	page := GetIntParameter(c, "page", 1)
	pageSize := GetIntParameter(c, "pageSize", 20)

	var reqBody = models.PushQueryVO{
		PageQueryVO: models.PageQueryVO{
			PageSize: pageSize,
			Page:     page,
		},
	}

	total, arr, err := service.PushService().List(&reqBody)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(&PushQueryResultVO{
		Page:     reqBody.Page,
		PageSize: reqBody.PageSize,
		Total:    total,
		Data:     arr,
	}))
}

func SavePushSource(c *gin.Context) {
	var source models.PushSource
	if err := c.ShouldBindJSON(&source); err != nil {
		logger.Error(err)
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	source.Url = strings.TrimSpace(source.Url)
	if source.Url == "" {
		c.JSON(http.StatusInternalServerError, ErrorMsg("url is empty"))
		return
	}

	source.Method = strings.TrimSpace(source.Method)
	if source.Method == "" {
		source.Method = "POST"
	}

	if source.PushSize <= 0 {
		source.PushSize = 50
	}

	if source.PushInterval <= 0 {
		source.PushInterval = 0
	}

	err := service.PushService().Save(source)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(nil))
}

func DeletePush(c *gin.Context) {
	var pushIds []int
	if err := c.ShouldBindJSON(&pushIds); err != nil {
		logger.Error(err)
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	if len(pushIds) == 0 {
		c.JSON(http.StatusOK, Success(nil))
		return
	}

	err := service.PushService().Delete(pushIds)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(nil))
}
