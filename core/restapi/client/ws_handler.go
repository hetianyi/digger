///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package client

import (
	"bytes"
	"digger/crawler"
	"digger/models"
	"errors"
	"fmt"
	"github.com/hetianyi/gox/logger"
)

func processQueue(queue *models.Queue) (*models.QueueProcessResult, error) {
	// 生成的新queue
	var returnNewQueues []*models.Queue
	// 产生的爬虫结果
	var results []*models.Result
	var log = &models.InMemLogWriter{
		Data: new(bytes.Buffer),
	}

	errChan := make(chan error)

	project := GetConfigSnapshot(queue.TaskId) //service.CacheService().GetSnapshotConfig(queue.TaskId)
	if project == nil {
		return nil, errors.New("cannot get project config snapshot")
	}

	Put(&job{
		finishChan: errChan,
		job: func() error {
			return crawler.Process(queue, project, log, func(cxt *models.Context, oldQueue *models.Queue, newQueues []*models.Queue, _results []*models.Result, err error) {
				if err != nil {
					return
				}
				results = append(results, _results...)
				returnNewQueues = append(returnNewQueues, newQueues...)
			})
		},
	})
	err := <-errChan
	if err != nil {
		logger.Error(err)
		log.Write([]byte(fmt.Sprintf("<span style=\"color:#F38F8F\">Err: process queue(%d): %s</span>\n", queue.Id, err.Error())))
		errMsg := ""
		if err != crawler.RobotsBlockErr {
			errMsg = log.Get()
		}
		return &models.QueueProcessResult{
			TaskId:    queue.TaskId,
			QueueId:   queue.Id,
			Expire:    queue.Expire,
			Logs:      log.Get(),
			NewQueues: nil,
			Results:   nil,
			Error:     errMsg,
		}, err
	}
	return &models.QueueProcessResult{
		TaskId:    queue.TaskId,
		QueueId:   queue.Id,
		Expire:    queue.Expire,
		Logs:      log.Get(),
		NewQueues: returnNewQueues,
		Results:   results,
	}, nil
}
