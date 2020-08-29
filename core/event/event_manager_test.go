///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package event_test

import (
	"digger/event"
	"digger/models"
	"digger/services/service"
	"testing"
)

func init() {
	service.InitRedis("123456@192.168.0.100:20021#0")
}

func TestStartEventManager(t *testing.T) {
	event.StartEventManager()
}

func TestPublish(t *testing.T) {
	//event.Publish("你好")
	event.Publish(&models.RedisEvent{
		Event: 1,
		Body:  nil,
	})
}
