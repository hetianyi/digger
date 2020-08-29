///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package client

import (
	"digger/models"
	"digger/utils"
	"fmt"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"sync"
)

var (
	snapshotConfigCache     = make(map[int]*models.Project)
	snapshotConfigCacheLock = new(sync.Mutex)
)

type respBody struct {
	models.Resp
	Data *models.Project `json:"data"`
}

func GetConfigSnapshot(taskId int) *models.Project {
	snapshotConfigCacheLock.Lock()
	defer snapshotConfigCacheLock.Unlock()

	if snapshotConfigCache[taskId] != nil {
		return snapshotConfigCache[taskId]
	}

	managerUrl, err := utils.Parse(config.ManagerUrl)
	if err != nil {
		logger.Fatal("invalid manager url: ", config.ManagerUrl, ": ", err.Error())
	}

	resp, err := httpClient.R().
		SetHeader("Accept", "application/json").
		SetHeader("Content-Type", "application/json").
		SetQueryParam("taskId", convert.IntToStr(taskId)).
		Get(managerUrl.Scheme + "://" + managerUrl.Host + "/api/v1/snapshot")

	if err != nil {
		logger.Error(fmt.Sprintf("error request config snapshot: %s", err.Error()))
		return nil
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Error(fmt.Sprintf("error request config snapshot: %d", resp.StatusCode()))
		return nil
	}
	var respEntity respBody
	err = jsoniter.Unmarshal(resp.Body(), &respEntity)
	if err != nil {
		logger.Error(fmt.Sprintf("error request config snapshot: %s", err.Error()))
		return nil
	}
	if respEntity.Code != 0 {
		logger.Error(fmt.Sprintf("error request config snapshot: invalid response code %d: %s", respEntity.Code, respEntity.Message))
		return nil
	}
	project := respEntity.Data
	if project == nil {
		logger.Error(fmt.Sprintf("cannot get config snapshot"))
		return nil
	}
	if project.Id > 0 {
		snapshotConfigCache[taskId] = project
		return project
	}
	return nil
}
