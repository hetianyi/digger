package crawler

import (
	"digger/models"
	"digger/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/hetianyi/gox/logger"
)

// 分页选择器处理
func processList(cxt *models.Context) error {
	stage := cxt.Stage
	if stage.ListCss != "" {
		return processListByCssSelector(cxt)
	}
	if stage.ListXpath != "" {
		return processListByXpathSelector(cxt)
	}
	logger.Debug("no list selector")
	return nil
}

// css选择器
func processListByCssSelector(cxt *models.Context) error {

	queue := cxt.Queue
	stage := cxt.Stage

	doc, err := parseCssDocument(cxt)
	if err != nil {
		return err
	}

	if len(stage.Fields) == 0 {
		return nil
	}
	doc.Find(stage.ListCss).Each(func(i int, s *goquery.Selection) {
		var itemMiddleData []*models.Queue
		for i := range stage.Fields {
			f := stage.Fields[i]
			ret := processCssField(cxt, &f, s)

			cxt.Log.Write([]byte(fmt.Sprintf("%s: %s", f.Name, ret)))

			if f.NextStage != "" {
				nextStageUrl := ""
				if f.Plugin != nil {
					nextStageUrl = ret.(string)
				} else {
					nextStageUrl, _ = utils.AbsoluteURL(cxt.Queue.Url, ret.(string))
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
			cxt.NewQueues = append(cxt.NewQueues, i)
		}

		if !stage.HasNextStage {
			cxt.AddResult(&models.Result{
				TaskId: queue.TaskId,
				Result: temp,
			})
		}
	})
	return nil
}

// xpath选择器
func processListByXpathSelector(cxt *models.Context) error {

	queue := cxt.Queue
	stage := cxt.Stage
	doc, err := parseXpathDocument(cxt)
	if err != nil {
		return err
	}
	list, err := htmlquery.QueryAll(doc, stage.ListXpath)
	if err != nil {
		return err
	}
	for _, node := range list {
		var itemMiddleData []*models.Queue
		for i := range stage.Fields {
			f := stage.Fields[i]
			ret := processXpathField(cxt, &f, node)

			cxt.Log.Write([]byte(fmt.Sprintf("%s: %s", f.Name, ret)))

			if f.NextStage != "" {
				nextStageUrl := ""
				if f.Plugin != nil {
					nextStageUrl = ret.(string)
				} else {
					nextStageUrl, _ = utils.AbsoluteURL(cxt.Queue.Url, ret.(string))
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
			cxt.NewQueues = append(cxt.NewQueues, i)
		}

		if !stage.HasNextStage {
			cxt.AddResult(&models.Result{
				TaskId: queue.TaskId,
				Result: temp,
			})
		}
	}
	return nil
}
