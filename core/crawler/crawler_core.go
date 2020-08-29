///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package crawler

import (
	"bytes"
	"digger/models"
	"digger/plugins"
	"digger/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"github.com/json-iterator/go"
	"io"
	"strings"
	"time"
)

var (
	httpClient = resty.New()
	json       = jsoniter.ConfigFastest
)

func init() {
	httpClient.SetTimeout(time.Second * 30)
}

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
		NewQueues:  make(map[string]*models.Queue),
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
		NewQueues:  make(map[string]*models.Queue),
		MiddleData: make(map[string]string),
	}
	plugins.InitVM(cxt)

	// slot s1
	err := handleS1(cxt)
	if err != nil {
		logger.Error("error run plugin at slot s1: ", err.Error())
	}
	// slot sr
	err = handleSR(cxt)
	if err != nil {
		callback(queue, nil, nil, err)
		return err
	}
	// slot s2
	handleS2(cxt)
	handlerEngineRoute(cxt, callback)
	return nil
}

func request(queue *models.Queue, project *models.Project) (*resty.Response, error) {
	return httpClient.R().SetHeaders(project.Headers).Get(queue.Url)
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

	queue := cxt.Queue

	doc, err := parseDocument([]byte(cxt.ResponseData))
	if err != nil {
		callback(queue, nil, nil, err)
		return
	}

	extendsData(cxt)

	stage := cxt.Stage

	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
	// 一个stage，可能由pagecss产生一个queue，可能由field产生若干queue，                                                //
	// 如果field都没有next_stage，则是final stage，产生爬虫最终结果，否则产生的结果均保存在中间结果queue的middle_data中   //
	////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

	if nextPage := processPageCss(cxt, doc.Selection); nextPage != "" {
		cxt.NewQueues[nextPage] = &models.Queue{
			TaskId:     queue.TaskId,
			StageName:  stage.Name,
			Url:        nextPage,
			MiddleData: queue.MiddleData,
		}
		cxt.Log.Write([]byte(fmt.Sprintf("Next page: %s", nextPage)))
	}

	cxt.Log.Write([]byte(fmt.Sprintf("\n<span style=\"color:#47F6E0\">================== Process Result of %s</span>\n", queue.Url)))

	// 如果stage是list类型，则循环list
	if stage.IsList {
		if len(stage.Fields) > 0 {
			doc.Find(stage.ListCss).Each(func(i int, s *goquery.Selection) {

				var itemMiddleData []*models.Queue

				for i := range stage.Fields {
					f := stage.Fields[i]
					ret := processField(&f, s)
					// slot s4
					ret = handleS4(cxt, &f, f.Name, ret)

					cxt.Log.Write([]byte(fmt.Sprintf("%s: %s", f.Name, ret)))

					if f.NextStage != "" {
						nextStageUrl := ""
						if f.Plugin != nil {
							nextStageUrl = ret
						} else {
							nextStageUrl, _ = utils.AbsoluteURL(cxt.Queue.Url, ret)
						}
						//cxt.Log.Write([]byte(fmt.Sprintf("Next stage: %s", nextStageUrl)))
						itemMiddleData = append(itemMiddleData, &models.Queue{
							TaskId:    queue.TaskId,
							StageName: f.NextStage,
							Url:       nextStageUrl,
						})
					}
					cxt.MiddleData[f.Name] = ret
				}

				temp, _ := json.MarshalToString(cxt.MiddleData)
				for _, i := range itemMiddleData {
					i.MiddleData = temp
					cxt.NewQueues[i.Url] = i
				}

				if !stage.HasNextStage {
					cxt.AddResult(&models.Result{
						TaskId: queue.TaskId,
						Result: temp,
					})
				}
			})
		}
	} else {
		var itemMiddleData []*models.Queue
		for i := range stage.Fields {
			f := stage.Fields[i]
			ret := processField(&f, doc.Selection)
			// slot s4
			ret = handleS4(cxt, &f, f.Name, ret)

			cxt.Log.Write([]byte(fmt.Sprintf("%s: %s", f.Name, ret)))

			if f.NextStage != "" {
				nextStageUrl := ""
				if f.Plugin != nil {
					nextStageUrl = ret
				} else {
					nextStageUrl, _ = utils.AbsoluteURL(cxt.Queue.Url, ret)
				}
				cxt.Log.Write([]byte(fmt.Sprintf("Next stage: %s", nextStageUrl)))
				itemMiddleData = append(itemMiddleData, &models.Queue{
					TaskId:    queue.TaskId,
					StageName: f.NextStage,
					Url:       nextStageUrl,
				})
			}
			cxt.MiddleData[f.Name] = ret
		}
		temp, _ := json.MarshalToString(cxt.MiddleData)
		fmt.Println(temp)
		for _, i := range itemMiddleData {
			i.MiddleData = temp
			cxt.NewQueues[i.Url] = i
		}
		if !stage.HasNextStage {
			cxt.AddResult(&models.Result{
				TaskId: queue.TaskId,
				Result: temp,
			})
		}
	}
	cxt.Log.Write([]byte(fmt.Sprintf("\n==================\n")))

	var newQueues []*models.Queue
	for _, v := range cxt.NewQueues {
		newQueues = append(newQueues, v)
	}
	callback(queue, newQueues, cxt.Results, nil)
}

// 处理stage的分页
func processPageCss(cxt *models.Context, s *goquery.Selection) string {
	stage := cxt.Stage
	if stage.PageCss != "" {
		nextPage := ""
		sel := s.Find(stage.PageCss)
		if stage.PageAttr == "" {
			nextPage = sel.Text()
		} else {
			if val, exists := sel.Attr(stage.PageAttr); exists {
				nextPage = strings.TrimSpace(val)
			}
		}
		if nextPage != "" {
			plugin := stage.FindPlugins("s4")
			if plugin != nil {
				// slot s4
				nextPage = handleStageS4(cxt, stage, nextPage)
			} else {
				nextPage, _ = utils.AbsoluteURL(cxt.Queue.Url, nextPage)
			}
			fmt.Println(fmt.Sprintf("下一页: %s", nextPage))
		}
		return nextPage
	}
	return ""
}

// 处理stage的字段
func processField(field *models.Field, s *goquery.Selection) string {
	ret := ""
	if field.IsArray {
		var arrayFieldValue []string
		var sel = s
		if field.Css != "" {
			sel = s.Find(field.Css)
		}
		// 如果不是list类型，则直接匹配fields
		// 循环fields，对于list的每个element进行处理
		sel.Each(func(i int, selection *goquery.Selection) {
			v := ""
			if field.Attr == "" {
				if field.IsHtml {
					v, _ = selection.Html()
				} else {
					v = strings.TrimSpace(selection.Text())
				}
			} else {
				if val, exists := selection.Attr(field.Attr); exists {
					v = strings.TrimSpace(val)
				}
			}
			if v != "" {
				arrayFieldValue = append(arrayFieldValue, v)
			}
		})
		ret, _ = json.MarshalToString(arrayFieldValue)
	} else {
		var sel = s
		if field.Css != "" {
			sel = s.Find(field.Css)
		}
		v := ""
		if field.Attr == "" {
			if field.IsHtml {
				v, _ = sel.Html()
			} else {
				v = strings.TrimSpace(sel.Text())
			}
		} else {
			if val, exists := sel.Attr(field.Attr); exists {
				v = strings.TrimSpace(val)
			}
		}
		if v != "" {
			ret = v
		}
	}
	return ret
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
func parseDocument(content []byte) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(bytes.NewBuffer(content))
}
