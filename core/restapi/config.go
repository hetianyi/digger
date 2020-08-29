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
	"net/http"
)

// 获取配置
func GetConfigs(c *gin.Context) {
	congis, err := service.ConfigService().ListConfigs()
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(congis))
}

// 获取配置
func UpdateConfig(c *gin.Context) {
	// 绑定请求数据
	reqData := &models.Config{}
	if err := c.ShouldBindJSON(reqData); err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	err := service.ConfigService().UpdateConfig(reqData.Key, reqData.Value)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(nil))
}
