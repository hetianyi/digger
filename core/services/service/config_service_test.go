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

func TestConfigServiceImpl_UpdateConfig(t *testing.T) {
	c := &models.EmailConfig{
		Host:     "192.168.0.101",
		Port:     6379,
		Username: "1129353184@qq.com",
		Password: "xxxxsdasd4dfs",
	}
	fmt.Println(service.ConfigService().UpdateConfig("email_config", c.String()))
}

func TestConfigServiceImpl_ListConfigs(t *testing.T) {
	fmt.Println(service.ConfigService().ListConfigs())
}
