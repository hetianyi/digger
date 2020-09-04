package crawler

import (
	"digger/models"
	"digger/utils"
	"fmt"
)

// 分页选择器处理
func processNoneList(cxt *models.Context) error {
	stage := cxt.Stage
	queue := cxt.Queue

	if len(stage.Fields) == 0 {
		return nil
	}
	var itemMiddleData []*models.Queue
	for i := range stage.Fields {
		f := stage.Fields[i]
		ret := ""
		if f.Css != "" {
			doc, err := parseCssDocument(cxt)
			if err != nil {
				return err
			}
			ret = processCssField(&f, doc.Selection)
		} else if f.Xpath != "" {
			doc, err := parseXpathDocument(cxt)
			if err != nil {
				return err
			}
			ret = processXpathField(&f, doc)
		}
		// slot s4
		ret = handleS4(cxt, &f, f.Name, ret)

		cxt.Log.Write([]byte(fmt.Sprintf("%s: %s\n", f.Name, ret)))

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
	return nil
}
