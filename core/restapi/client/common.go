///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package client

import (
	"github.com/go-resty/resty/v2"
)

type Response struct {
}

var (
	httpClient  *resty.Client
	fetchSize   = 5
	managerUrls []string // only one for now
)

func init() {
	httpClient = resty.New()
}

func ManagerUrl(_managerUrls []string) {
	managerUrls = _managerUrls
}
