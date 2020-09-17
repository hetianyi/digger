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

type ProxyQueryResultVO struct {
	Page     int             `json:"page"`
	PageSize int             `json:"pageSize"`
	Total    int64           `json:"total"`
	Data     []*models.Proxy `json:"data"`
}

func QueryProxy(c *gin.Context) {
	page := GetIntParameter(c, "page", 1)
	pageSize := GetIntParameter(c, "pageSize", 20)
	key := GetStrParameter(c, "key", "")

	var reqBody = models.ProxyQueryVO{
		PageQueryVO: models.PageQueryVO{
			PageSize: pageSize,
			Page:     page,
		},
		Key: key,
	}

	total, arr, err := service.ProxyService().List(&reqBody)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(&ProxyQueryResultVO{
		Page:     reqBody.Page,
		PageSize: reqBody.PageSize,
		Total:    total,
		Data:     arr,
	}))
}

func SaveProxy(c *gin.Context) {
	var proxy models.Proxy
	if err := c.ShouldBindJSON(&proxy); err != nil {
		logger.Error(err)
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	proxy.Address = strings.TrimSpace(proxy.Address)
	if proxy.Address == "" {
		c.JSON(http.StatusInternalServerError, ErrorMsg("address is empty"))
		return
	}

	err := service.ProxyService().Save(proxy)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(nil))
}

func DeleteProxy(c *gin.Context) {
	var proxyIds []int
	if err := c.ShouldBindJSON(&proxyIds); err != nil {
		logger.Error(err)
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	if len(proxyIds) == 0 {
		c.JSON(http.StatusOK, Success(nil))
		return
	}

	err := service.ProxyService().Delete(proxyIds)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(nil))
}
