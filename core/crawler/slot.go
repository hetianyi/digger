///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package crawler

import (
	"digger/models"
	"errors"
	"github.com/hetianyi/gox/convert"
	"net/http"
)

func handleS1(cxt *models.Context) error {
	plugin := cxt.Stage.FindPlugins("s1")
	if plugin == nil {
		return nil
	}
	s, err := cxt.Exec(plugin.Script)
	if err != nil {
		return err
	}
	if s != "" {
		cxt.Queue.Url = s
	}
	return nil
}

func handleSR(cxt *models.Context) error {
	plugin := cxt.Stage.FindPlugins("sr")
	if plugin == nil {
		resp, err := request(cxt.Queue, cxt.Project)
		if err != nil {
			return err
		}

		cxt.HttpStatusCode = resp.StatusCode()

		if resp.StatusCode() == http.StatusOK {
			body := string(resp.Body())
			cxt.ResponseData = body
			return nil
		} else {

			return errors.New("error http status " + convert.IntToStr(resp.StatusCode()))
		}
	}
	data, err := cxt.Exec(plugin.Script)
	cxt.ResponseData = data
	return err
}

func handleS2(cxt *models.Context) error {
	plugin := cxt.Stage.FindPlugins("s2")
	if plugin == nil {
		return nil
	}
	s, err := cxt.Exec(plugin.Script)
	if err != nil {
		return err
	}
	if s != "" {
		cxt.ResponseData = s
	}
	return nil
}

func hasS3(cxt *models.Context) bool {
	plugin := cxt.Stage.FindPlugins("s3")
	if plugin == nil {
		return false
	}
	return true
}

func handleS3(cxt *models.Context) error {
	plugin := cxt.Stage.FindPlugins("s3")
	if plugin == nil {
		return nil
	}
	_, err := cxt.Exec(plugin.Script)
	if err != nil {
		return err
	}
	return nil
}

func handleS4(cxt *models.Context, field *models.Field, fieldName, fieldValue string) string {
	cxt.ENV["currentFieldName"] = fieldName
	cxt.ENV["currentFieldValue"] = fieldValue
	if field.Plugin == nil {
		return fieldValue
	}
	finalValue, err := cxt.Exec(field.Plugin.Script)
	if err != nil {
		return fieldValue
	}
	return finalValue
}

func handleStageS4(cxt *models.Context, stage *models.Stage, nextPage string) string {
	cxt.ENV["currentFieldName"] = "$" + stage.Name
	cxt.ENV["currentFieldValue"] = nextPage
	plugin := stage.FindPlugins("s4")
	if plugin == nil {
		return nextPage
	}
	finalValue, err := cxt.Exec(plugin.Script)
	if err != nil {
		return nextPage
	}
	return finalValue
}
