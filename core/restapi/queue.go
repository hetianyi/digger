///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"digger/utils"
	"github.com/gin-gonic/gin"
	"github.com/hetianyi/gox/convert"
	jsoniter "github.com/json-iterator/go"
)

// 为项目开始一个新任务
func GetQueues(c *gin.Context) {
	query := c.Query("limit")
	limit := 10
	if query != "" {
		l, err := convert.StrToInt(query)
		if err == nil {
			limit = l
		}
	}
	if limit <= 0 {
		limit = 1
	}
	if limit > 50 {
		limit = 50
	}

	labelsStr := c.GetHeader("Labels")
	labels := make(map[string]string)
	if labelsStr != "" {
		s, _ := utils.DecodeBase64(labelsStr)
		if s != "" {
			jsoniter.UnmarshalFromString(s, &labels)
		}
	}

	/*check, queues := scheduler.GetQueues(limit, labels)
	c.JSON(http.StatusOK, Success(&models.FetchQueueResponseVO{
		Check:  check,
		Queues: queues,
	}))*/
}
