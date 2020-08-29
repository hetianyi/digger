///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service_test

import (
	"digger/models"
	"digger/services/service"
	"fmt"
	"testing"
)

func TestResultServiceImp_InsertResults(t *testing.T) {
	service.ResultService().InsertResults(1, []models.Result{
		{
			Id:     0,
			TaskId: 1,
			Result: "{\"field1\":\"value1撒打算的阿萨德啊\"}",
		},
		{
			Id:     0,
			TaskId: 1,
			Result: "{\"field2\":\"value1撒打算的阿萨德啊\"}",
		},
		{
			Id:     0,
			TaskId: 1,
			Result: "{\"field3\":\"value1撒打算的阿萨德啊\"}",
		},
	})
}

func TestResultServiceImp_ResultCount(t *testing.T) {
	service.ResultService().SaveProcessResultData(&models.QueueProcessResult{
		TaskId:    19,
		InitUrl:   "",
		QueueId:   1495,
		Expire:    0,
		RequestId: "",
		Error:     "",
		Logs:      "",
		NewQueues: []*models.Queue{
			{
				Id:         0,
				TaskId:     19,
				StageName:  "xxx111",
				Url:        "http://111",
				MiddleData: "{}",
				Expire:     11,
			},
			{
				Id:         0,
				TaskId:     19,
				StageName:  "xxx222",
				Url:        "http://222",
				MiddleData: "{}",
				Expire:     11,
			},
		},
		Results: nil,
	}, false)
}

func TestResultServiceImp_SaveCheckData(t *testing.T) {
	fmt.Println("1"[0:1])
}
