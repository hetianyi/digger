///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"compress/gzip"
	"digger/dispatcher"
	"digger/models"
	"digger/services/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hetianyi/gox"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strings"
	"time"
)

type Item struct {
	Date  string `json:"date"`
	Value int    `json:"value"`
}

type LineType struct {
	Type  string          `json:"type"`
	Name  string          `json:"name"`
	Color string          `json:"color"`
	Data  [][]interface{} `json:"data"`
}

func GetStatistic(c *gin.Context) {

	start := strings.TrimSpace(GetStrParameter(c, "start", ""))
	end := strings.TrimSpace(GetStrParameter(c, "end", ""))

	startTime, _ := time.Parse("2006-01-02 15:04:05", start)
	endTime, _ := time.Parse("2006-01-02 15:04:05", end)

	vos, err := service.StatisticService().List(startTime, endTime)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	configs, err := service.ConfigService().ListConfigs()
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	retMap := make(map[string]interface{})
	var ret []*LineType
	totalRequests := &LineType{
		Name:  "请求数",
		Type:  "line",
		Color: "#00F730",
		Data:  [][]interface{}{},
	}
	errorRequests := &LineType{
		Name:  "错误数",
		Type:  "line",
		Color: "#FF0000",
		Data:  [][]interface{}{},
	}
	result := &LineType{
		Name:  "结果数",
		Color: "#B117AC",
		Type:  "line",
		Data:  [][]interface{}{},
	}
	ret = append(ret, totalRequests, errorRequests, result)

	retMap["series"] = ret

	retMap["inforCardData"] = []map[string]interface{}{
		{
			"title": "项目",
			"icon":  "md-apps",
			"count": configs["project_count"],
			"color": "#2d8cf0",
		},
		{
			"title": "任务",
			"icon":  "ios-stats",
			"count": configs["task_count"],
			"color": "#19be6b",
		},
		{
			"title": "工作节点",
			"icon":  "ios-globe-outline",
			"count": dispatcher.CountClient(),
			"color": "#ff9900",
		},
		{
			"title": "累计请求",
			"icon":  "md-paper-plane",
			"count": configs["total_request_count"],
			"color": "#ed3f14",
		},
		{
			"title": "累计结果",
			"icon":  "md-checkbox",
			"count": configs["result_count"],
			"color": "#E46CBB",
		},
	}

	for _, e := range vos {
		val := 0
		if e.Data["request_count"] != nil {
			val = int(e.Data["request_count"].(float64))
		}
		totalRequests.Data = append(totalRequests.Data, []interface{}{gox.GetLongDateString(e.CreateTime), val})

		if e.Data["result_count"] != nil {
			val = int(e.Data["result_count"].(float64))
		}
		result.Data = append(result.Data, []interface{}{gox.GetLongDateString(e.CreateTime), val})

		if e.Data["error_request_count"] != nil {
			val = int(e.Data["error_request_count"].(float64))
		}
		errorRequests.Data = append(errorRequests.Data, []interface{}{gox.GetLongDateString(e.CreateTime), val})
	}

	success := Success(retMap)
	bytes, _ := jsoniter.Marshal(success)
	fmt.Println(len(bytes))
	c.Writer.Header().Set("Content-Encoding", "gzip")
	gwriter := gzip.NewWriter(c.Writer)
	gwriter.Write(bytes)
	gwriter.Close()
}

func group(vos []*models.StatisticVO, parts int, step time.Duration) {

}
