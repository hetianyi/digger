///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package crawler

import (
	"bytes"
	"digger/httpclient"
	"digger/models"
	"digger/plugins"
	"digger/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/go-resty/resty/v2"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"github.com/json-iterator/go"
	"golang.org/x/net/html"
	"io"
	"net/http"
)

var (
	json = jsoniter.ConfigFastest
)

// 根据爬虫任务和队列数据开始处理
func Process(
	queue *models.Queue,
	project *models.Project,
	log io.Writer,
	callback func(oldQueue *models.Queue, newQueue []*models.Queue, results []*models.Result, err error)) error {

	stage := project.GetStageByName(queue.StageName)

	cxt := &models.Context{
		Project: project,
		Queue:   queue,
		Stage:   stage,
		Log:     log,
		ENV: map[string]string{
			"stage":  queue.StageName,
			"taskId": convert.IntToStr(queue.TaskId),
		},
		NewQueues:  []*models.Queue{},
		MiddleData: make(map[string]string),
	}
	plugins.InitVM(cxt)

	// slot s1
	err := handleS1(cxt)
	if err != nil {
		log.Write([]byte(fmt.Sprintf("<span style=\"color:#F38F8F\">Err: run plugin at slot s1: %s</span>", err.Error())))
		logger.Error("error run plugin at slot s1: ", err.Error())
	}

	log.Write([]byte(fmt.Sprintf("crawling => %s", queue.Url)))

	// slot sr
	err = handleSR(cxt)
	if err != nil {
		log.Write([]byte(fmt.Sprintf("<span style=\"color:#F38F8F\">Err: process url: %s: %s</span>", queue.Url, err.Error())))
		callback(queue, nil, nil, err)
		return err
	}
	// slot s2
	handleS2(cxt)
	handlerEngineRoute(cxt, callback)
	return nil
}

func Play(
	queue *models.Queue,
	project *models.Project,
	log io.Writer,
	callback func(oldQueue *models.Queue, newQueue []*models.Queue, results []*models.Result, err error)) error {

	stage := project.GetStageByName(queue.StageName)
	cxt := &models.Context{
		Project: project,
		Queue:   queue,
		Stage:   stage,
		Log:     log,
		ENV: map[string]string{
			"stage":  queue.StageName,
			"taskId": convert.IntToStr(queue.TaskId),
		},
		NewQueues:  []*models.Queue{},
		MiddleData: make(map[string]string),
	}
	plugins.InitVM(cxt)

	// slot s1 请求之前的url插槽
	err := handleS1(cxt)
	if err != nil {
		logger.Error("error run plugin at slot s1: ", err.Error())
	}
	// slot sr: http请求插槽
	err = handleSR(cxt)
	if err != nil {
		callback(queue, nil, nil, err)
		return err
	}
	// slot s2 请求之后结果预处理插槽
	handleS2(cxt)
	// 处理引擎结果路由
	handlerEngineRoute(cxt, callback)
	return nil
}

func request(queue *models.Queue, cxt *models.Context) (*resty.Response, error) {
	url, err := utils.Parse(queue.Url)
	if err != nil {
		return nil, err
	}
	// TODO bugs
	client := httpclient.GetClient(queue.TaskId, cxt.Project)
	feedback := utils.TryProxy(url.Scheme, client, queue.TaskId, cxt)
	response, err := client.R().
		SetHeaders(cxt.Project.Headers).
		Get(queue.Url)
	// feedback
	if feedback != nil {
		if err != nil || response.StatusCode() != http.StatusOK {
			feedback.Fail()
		} else {
			feedback.Success()
		}
	}
	return response, err
}

func handlerEngineRoute(
	cxt *models.Context,
	callback func(oldQueue *models.Queue, newQueue []*models.Queue, results []*models.Result, err error)) {

	if hasS3(cxt) {
		extendsData(cxt)
		err := handleS3(cxt)
		processAfterS3(err, cxt, callback)
		return
	}
	processDefaultStage(cxt, callback)
}

func processAfterS3(
	err error,
	cxt *models.Context,
	callback func(oldQueue *models.Queue, newQueue []*models.Queue, results []*models.Result, err error)) {

	oldQueue := cxt.Queue

	if err != nil {
		callback(oldQueue, nil, nil, err)
		return
	}
	if err != nil {
		callback(oldQueue, nil, nil, err)
		return
	}

	newQueueMap := cxt.NewQueues
	var newQueue []*models.Queue
	for _, q := range newQueueMap {
		q.TaskId = oldQueue.TaskId
		newQueue = append(newQueue, q)
	}
	callback(oldQueue, newQueue, cxt.Results, nil)
}

// 处理stage入口
func processDefaultStage(
	cxt *models.Context,
	callback func(oldQueue *models.Queue, newQueue []*models.Queue, results []*models.Result, err error)) {

	extendsData(cxt)
	queue := cxt.Queue

	stage := cxt.Stage

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 一个stage，可能由pagecss产生一个queue，可能由field产生若干queue，                                                //
	// 如果field都没有next_stage，则是final stage，产生爬虫最终结果，否则产生的结果均保存在中间结果queue的middle_data中   //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	nextPage, err := processPage(cxt)
	if err != nil {
		callback(queue, nil, nil, err)
		return
	}

	if nextPage != "" {
		cxt.NewQueues = append(cxt.NewQueues, &models.Queue{
			TaskId:     queue.TaskId,
			StageName:  stage.Name,
			Url:        nextPage,
			MiddleData: queue.MiddleData,
		})
		cxt.Log.Write([]byte(fmt.Sprintf("Next page: %s", nextPage)))
	}

	cxt.Log.Write([]byte(fmt.Sprintf("\n<span style=\"color:#47F6E0\">================== Process Result of %s</span>\n", queue.Url)))

	// 如果stage是list类型，则循环list
	if stage.IsList {
		if err := processList(cxt); err != nil {
			callback(queue, nil, nil, err)
		}
	} else {
		if err := processNoneList(cxt); err != nil {
			callback(queue, nil, nil, err)
		}
	}
	cxt.Log.Write([]byte(fmt.Sprintf("\n==================\n")))

	var newQueues []*models.Queue
	for _, v := range cxt.NewQueues {
		newQueues = append(newQueues, v)
	}
	callback(queue, newQueues, cxt.Results, nil)
}

func extendsData(cxt *models.Context) {
	cxt.MiddleData = make(map[string]string)
	queue := cxt.Queue

	oldDataMap := make(map[string]string)
	if cxt.Queue.MiddleData != "" {
		if err := json.UnmarshalFromString(queue.MiddleData, &oldDataMap); err != nil {
			return
		}
	}

	for k, v := range oldDataMap {
		cxt.MiddleData[k] = v
	}
}

// 用goquery解析html文档
func parseCssDocument(cxt *models.Context) (*goquery.Document, error) {
	if cxt.CssQueryDoc != nil {
		return cxt.CssQueryDoc, nil
	}
	d, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(cxt.ResponseData)))
	cxt.CssQueryDoc = d
	return d, err
}

// 用goquery解析html文档
func parseXpathDocument(cxt *models.Context) (*html.Node, error) {
	if cxt.XpathQueryDoc != nil {
		return cxt.XpathQueryDoc, nil
	}
	doc, err := htmlquery.Parse(bytes.NewBuffer([]byte(cxt.ResponseData)))
	cxt.XpathQueryDoc = doc
	return doc, err
}
