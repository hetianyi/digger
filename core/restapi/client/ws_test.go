///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package client_test

import (
	"digger/models"
	"digger/restapi/client"
	"github.com/hetianyi/gox/logger"
	"testing"
)

func init() {
	// initialize logger
	logConfig := &logger.Config{
		Level:              logger.InfoLevel,
		Write2File:         false,
		AlwaysWriteConsole: true,
	}
	logger.Init(logConfig)
}

func TestInitWsClient(t *testing.T) {
	client.InitWsClient(&models.BootstrapConfig{
		ManagerUrl: "localhost:9012",
		Labels: map[string]string{
			"out": "true",
		},
	})
	<-make(chan int)
}
