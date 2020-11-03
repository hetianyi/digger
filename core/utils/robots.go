///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils

import (
	"digger/httpclient"
	"digger/models"
	"errors"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"net/url"
)

func LoadRobotsTxt(_url string, project *models.Project) ([]byte, error) {
	base, err := url.Parse(_url)
	if err != nil {
		logger.Warn("无法获取robots.txt: ", err)
		return nil, err
	}
	robotTxtUrl := base.Scheme + "://" + base.Host + gox.TValue(base.Port() == "", "", ":"+base.Port()).(string) + "/robots.txt"
	url, err := Parse(robotTxtUrl)
	if err != nil {
		return nil, err
	}
	client := httpclient.GetClient(0, project)
	TryProxy(url.Scheme, client, 0, nil) // TODO

	response, err := client.R().Get(robotTxtUrl)
	if err != nil {
		logger.Warn("无法获取robots.txt: ", err)
		return nil, err
	}
	if response.StatusCode() == 200 {
		return response.Body(), nil
	}
	return nil, errors.New("http status " + convert.IntToStr(response.StatusCode()))
}
