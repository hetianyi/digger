package utils

import (
	"digger/models"
	"github.com/hetianyi/gox/logger"
	jsoniter "github.com/json-iterator/go"
)

var encodedEvent struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

func EncodeEvent(event *models.RedisEvent) string {
	s, _ := jsoniter.MarshalToString(event)
	return s
}

func DecodeEvent(event string) *models.RedisEvent {
	var ret models.RedisEvent
	err := jsoniter.UnmarshalFromString(event, &ret)
	if err != nil {
		logger.Error("无法解析事件: ", err)
		return nil
	}
	return &ret
}
