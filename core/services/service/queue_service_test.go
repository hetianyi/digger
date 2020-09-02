package service_test

import (
	"digger/models"
	"digger/services/service"
	"fmt"
	"testing"
)

func TestQueueServiceImp_SelectQueues(t *testing.T) {
	list, _ := service.QueueService().SelectQueues(models.QueueQueryVO{
		TaskId: 1,
		Status: 0,
		Limit:  10,
	})
	fmt.Println(list)
}

func TestQueueServiceImp_InsertQueue(t *testing.T) {
	fmt.Println(service.QueueService().InsertQueue(models.Queue{
		Id:         0,
		TaskId:     1,
		StageName:  "xxx",
		Url:        "http://",
		MiddleData: "{}",
	}))
}

func TestQueueServiceImpl_GetUnFinishedCount(t *testing.T) {
	fmt.Println(service.QueueService().GetUnFinishedCount(1))
}


func TestQueueServiceImpl_StatisticFinal(t *testing.T) {
	service.QueueService().StatisticFinal(52)
}