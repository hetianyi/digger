///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package agents

import (
	"digger/agents/manager"
	"digger/agents/worker"
	"digger/common"
	"digger/models"
	"digger/utils"
	"encoding/json"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
)

var (
	config *models.BootstrapConfig
)

func Start(_config models.BootstrapConfig) {
	config = &_config

	initLogging()

	bytes, _ := json.MarshalIndent(_config, "", "  ")

	logger.Debug("\nbootstrap configuration:\n" + string(bytes))

	// 创建日志文件夹
	if !file.Exists(_config.LogDir) {
		if err := file.CreateDirs(_config.LogDir); err != nil {
			logger.Fatal(err)
		}
	}

	if config.BootMode == common.ROLE_MANAGER {
		manager.StartAgentManager(config)
	} else {
		config.Labels["id"] = convert.IntToStr(config.InstanceId)
		worker.StartAgentWorker(config)
	}
}

func initLogging() {
	// init logger
	logConfig := &logger.Config{
		Level:              utils.ConvertLogLevel(config.LogLevel),
		AlwaysWriteConsole: true,
		Write2File:         false,
	}
	logger.Init(logConfig)
}
