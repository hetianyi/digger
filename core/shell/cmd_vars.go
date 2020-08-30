///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package shell

import (
	"digger/common"
)

// var sets
var (
	bootMode    common.ROLE
	showVersion bool   // show app version
	logLevel    string // log level(trace, debug, info, warn, error, fatal)
	secret      string // secret of this instance
	port        int
	instanceId  int
	logDir      string
	managerUrl  string
	dbConn      string
	labels      string
	uiDir       string
	redisConn   string
)
