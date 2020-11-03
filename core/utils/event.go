///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
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
		logger.Error("error parse event payload: ", err)
		return nil
	}
	return &ret
}
