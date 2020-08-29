///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service_test

import (
	"digger/services/service"
	"github.com/hetianyi/gox/logger"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	// initialize logger
	logConfig := &logger.Config{
		Level:              logger.InfoLevel,
		Write2File:         false,
		AlwaysWriteConsole: true,
	}
	logger.Init(logConfig)

	service.InitDb("postgres://postgres:123456@localhost:5432/digger?sslmode=disable")
}
