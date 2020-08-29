///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"digger/models"
	"github.com/gin-gonic/gin"
	"github.com/hetianyi/gox/convert"
	"strings"
)

var DefaultUser = &models.User{
	Id:       0,
	Username: "admin",
	Password: "admin",
}

func Success(data interface{}) *models.RestResponse {
	return &models.RestResponse{
		Resp: models.Resp{
			Code:    0,
			Message: "Success",
		},
		Data: data,
	}
}

func Error() *models.RestResponse {
	return &models.RestResponse{
		Resp: models.Resp{
			Code:    1,
			Message: "Server Error",
		},
		Data: nil,
	}
}

func ErrorData(data interface{}) *models.RestResponse {
	return &models.RestResponse{
		Resp: models.Resp{
			Code:    1,
			Message: "Server Error",
		},
		Data: data,
	}
}

func ErrorMsg(msg string) *models.RestResponse {
	return &models.RestResponse{
		Resp: models.Resp{
			Code:    1,
			Message: msg,
		},
		Data: nil,
	}
}

func GetIntParameter(c *gin.Context, key string, defaultValue int) int {
	v, err := convert.StrToInt(c.Query(key))
	if err != nil {
		return defaultValue
	}
	return v
}

func GetInt64Parameter(c *gin.Context, key string, defaultValue int64) int64 {
	v, err := convert.StrToInt64(c.Query(key))
	if err != nil {
		return defaultValue
	}
	return v
}

func GetStrParameter(c *gin.Context, key string, defaultValue string) string {
	v := strings.TrimSpace(c.Query(key))
	if v == "" {
		return defaultValue
	}
	return v
}
