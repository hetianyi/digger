///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils_test

import (
	"digger/models"
	"digger/utils"
	"fmt"
	"github.com/hetianyi/gox"
	"testing"
	"time"
)

func TestParseRedisConnStr(t *testing.T) {
	fmt.Println(utils.ParseRedisConnStr("123@456@192.168.0.100:6379#1"))
	fmt.Println(utils.ParseRedisConnStr("@r@@1asd_-./@192.168.0.100:6379#1"))
	fmt.Println(utils.ParseRedisConnStr("@@@@192.168.0.100:6379#1"))
	fmt.Println(utils.ParseRedisConnStr("@192.168.0.100:6379#1"))
	fmt.Println(utils.ParseEmailNotifierStr("123@qq.com:123dsf@192.168.0.100:6379"))
	fmt.Println(gox.GetTimestamp(time.Now()))
	fmt.Println(gox.CreateTime(gox.GetTimestamp(time.Now()) - 2141))
}

func TestEmailNotify(t *testing.T) {

	c := &models.EmailConfig{
		Host:     "smtp.qq.com",
		Port:     465,
		Username: "xxx@qq.com",
		Password: "xxx",
	}

	fmt.Println(utils.EmailNotify(&models.Task{
		Id:         1,
		CreateTime: time.Unix((gox.GetTimestamp(time.Now())-214112000)/1000, 0),
	}, c))
}
