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
	response, err := httpclient.GetClient(0, project).R().Get(robotTxtUrl)
	if err != nil {
		logger.Warn("无法获取robots.txt: ", err)
		return nil, err
	}
	if response.StatusCode() == 200 {
		return response.Body(), nil
	}
	return nil, errors.New("http status " + convert.IntToStr(response.StatusCode()))
}
