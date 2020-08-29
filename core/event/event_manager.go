///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package event

import (
	"digger/common"
	"digger/models"
	"digger/services/service"
	"digger/utils"
	"github.com/go-redis/redis/v7"
	"github.com/hetianyi/gox/logger"
)

const (
	TASK_EVENT_CHANNEL = "TASK_EVENT_CHANNEL"
)

var (
	sub *redis.PubSub
)

func StartEventManager() {
	logger.Info("开始订阅事件消息")
	sub = service.RedisClient.Subscribe(TASK_EVENT_CHANNEL)
	ch := sub.Channel()
	var eve *redis.Message
	for {
		eve = <-ch
		//logger.Info("收到事件消息")
		event := utils.DecodeEvent(eve.Payload)
		eventRoute(event)
	}
}

func Publish(event *models.RedisEvent) error {
	_, err := service.RedisClient.Publish(TASK_EVENT_CHANNEL, utils.EncodeEvent(event)).Result()
	return err
}

func eventRoute(event *models.RedisEvent) {
	if event == nil {
		logger.Error("null event")
		return
	}
	switch event.Event {
	case common.EV_TASK_CREATED:
		handleTaskCreatedEvent(event)
		break
	case common.EV_TASK_PAUSE:
		handleTaskPauseEvent(event)
		break
	case common.EV_TASK_STOP:
		handleTaskStopEvent(event)
		break
	case common.EV_TASK_CONTINUE:
		handleTaskContinueEvent(event)
		break
	case common.EV_ONE_QUEUE_FINISH:
		handleQueueFinishEvent(event)
		break
	default:
		logger.Error("unrecognized event type: ", event.Event)
	}
}
