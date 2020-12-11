///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service_test

import (
	"digger/services/service"
	"fmt"
	"testing"
	"time"
)

func TestStatisticServiceImp_Save(t *testing.T) {
	service.StatisticService().Save(map[string]interface{}{
		"request_count":       120,
		"error_request_count": 5,
		"result_count":        206,
	})
}

func TestStatisticServiceImp_List(t *testing.T) {
	start := "2020-09-01 20:00:46"
	end := "2020-09-01 20:03:46"
	startTime, _ := time.Parse("2006-01-02 15:04:05", start)
	endTime, _ := time.Parse("2006-01-02 15:04:05", end)
	vos, _ := service.StatisticService().List(startTime, endTime)
	fmt.Println(vos)
}

func TestStatisticServiceImp_Save1(t *testing.T) {
	start := time.Now().Add(0 - time.Hour*120)
	now := time.Now()
	for now.Unix()-start.Unix() > 0 {
		service.StatisticService().Save(map[string]interface{}{
			"request_count":       120,
			"error_request_count": 5,
			"result_count":        206,
			"time":                start,
		})
		start = start.Add(time.Minute)
	}
}
