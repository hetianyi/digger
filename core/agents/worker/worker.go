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
	"time"
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

func watchAge(ageSec int) {
	if ageSec > 3600 {
		ageSec = 3600
	}
	if ageSec < 60 {
		ageSec = 60
	}
	logger.Info("开始生命倒计时：", ageSec, "s")
	time.Sleep(time.Second * time.Duration(ageSec))
	logger.Info("==================================================")
	logger.Info("!!!Exit due to exceed max age!!!")
	logger.Info("!!!Worker正常退出：已达到最大生命!!!")
	logger.Info("==================================================")
	os.Exit(0)
}
