///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package worker

import (
	"digger/models"
	"digger/restapi/client"
	"github.com/hetianyi/gox/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var (
	config *models.BootstrapConfig
	lock   = new(sync.Mutex)
)

// 启动manager
func StartAgentWorker(_config *models.BootstrapConfig) {

	config = _config
	client.ManagerUrl([]string{config.ManagerUrl})
	client.InitProcessors()
	startClient()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("worker is down!")
}

func startClient() {
	client.InitWsClient(config)
}
