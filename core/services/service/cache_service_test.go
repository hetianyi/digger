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

func init() {
	//service.InitRedis("123456@192.168.0.100:20021#0")
}

func TestCacheServiceImp_IncreQueueErrorCount(t *testing.T) {
	fmt.Println(service.CacheService().IncreQueueErrorCount([]int{1, 2, 3}, []int64{7, 8, 9}))
}

func TestCacheServiceImp_ExistMembers(t *testing.T) {
	service.CacheService().ExistMembers(1, []interface{}{1, 2, 3, 4})
}

func TestCacheServiceImp_SaveSuccessQueueIds(t *testing.T) {
	service.CacheService().SaveSuccessQueueIds(&models.QueueCallbackRequestVO{
		SuccessQueueIds:     []int64{1, 2, 3, 4},
		SuccessQueueTaskIds: []int{5, 5, 6, 6},
	})
}

func TestCacheServiceImp_IncreConcurrentTaskCount(t *testing.T) {
	fmt.Println(service.CacheService().IncreConcurrentTaskCount("xxx", 1, 1))
	fmt.Println(service.CacheService().IncreConcurrentTaskCount("aaa", 1, 1))
}

func TestCacheServiceImp_DecreConcurrentTaskCount(t *testing.T) {
	fmt.Println(service.CacheService().DecreConcurrentTaskCount("xxx", 1))
}
