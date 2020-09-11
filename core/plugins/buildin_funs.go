///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package plugins

import (
	"digger/httpclient"
	"digger/models"
	"digger/utils"
	"github.com/go-resty/resty/v2"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"github.com/robertkrimen/otto"
	"regexp"
	"strings"
)

var (
	cachedRegexp = make(map[string]*regexp.Regexp)
	httpClient   = resty.New()
)

func InitVM(cxt *models.Context) {
	cxt.VM = otto.New()
	initBuildInFunctions(cxt)
}

func initBuildInFunctions(cxt *models.Context) {

	cxt.VM.Set("LEN", func(call otto.FunctionCall) otto.Value {
		boo := false
		if len(call.ArgumentList) != 1 {
			logger.Error("script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList)))
			result, _ := cxt.VM.ToValue(boo)
			return result
		}
		result, _ := cxt.VM.ToValue(len([]rune(call.ArgumentList[0].String())))
		return result
	})

	cxt.VM.Set("STARTS_WITH", func(call otto.FunctionCall) otto.Value {
		boo := false
		if len(call.ArgumentList) != 2 {
			logger.Error("script Err: invalid arg number, expect 2, got " + convert.IntToStr(len(call.ArgumentList)))
			result, _ := cxt.VM.ToValue(boo)
			return result
		}
		result, _ := cxt.VM.ToValue(strings.HasPrefix(call.ArgumentList[0].String(), call.ArgumentList[1].String()))
		return result
	})

	cxt.VM.Set("END_WITH", func(call otto.FunctionCall) otto.Value {
		boo := false
		if len(call.ArgumentList) != 2 {
			logger.Error("script Err: invalid arg number, expect 2, got " + convert.IntToStr(len(call.ArgumentList)))
			result, _ := cxt.VM.ToValue(boo)
			return result
		}
		result, _ := cxt.VM.ToValue(strings.HasSuffix(call.ArgumentList[0].String(), call.ArgumentList[1].String()))
		return result
	})

	cxt.VM.Set("SUBSTR", func(call otto.FunctionCall) otto.Value {
		boo := false
		if len(call.ArgumentList) != 3 {
			logger.Error("script Err: invalid arg number, expect 3, got " + convert.IntToStr(len(call.ArgumentList)))
			result, _ := cxt.VM.ToValue(boo)
			return result
		}
		s := []rune(call.ArgumentList[0].String())
		start, _ := call.ArgumentList[1].ToInteger()
		end, _ := call.ArgumentList[2].ToInteger()

		if start >= int64(len(s)) {
			ret, _ := cxt.VM.ToValue("")
			return ret
		}
		if end > int64(len(s)) {
			end = int64(len(s))
		}
		result, _ := cxt.VM.ToValue(string(s[start:end]))
		return result
	})

	cxt.VM.Set("CONTAINS", func(call otto.FunctionCall) otto.Value {
		boo := false
		if len(call.ArgumentList) != 2 {
			logger.Error("script Err: invalid arg number, expect 2, got " + convert.IntToStr(len(call.ArgumentList)))
			result, _ := cxt.VM.ToValue(boo)
			return result
		}
		result, _ := cxt.VM.ToValue(strings.Contains(call.ArgumentList[0].String(), call.ArgumentList[1].String()))
		return result
	})

	cxt.VM.Set("REPLACE", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 3 {
			logger.Error("script Err: invalid arg number, expect 3, got " + convert.IntToStr(len(call.ArgumentList)))
			return otto.Value{}
		}
		result, _ := cxt.VM.ToValue(strings.Replace(call.ArgumentList[0].String(), call.ArgumentList[1].String(), call.ArgumentList[2].String(), -1))
		return result
	})

	cxt.VM.Set("REGEXP_GROUP_FIND", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 3 {
			logger.Error("script Err: invalid arg number, expect 3, got " + convert.IntToStr(len(call.ArgumentList)))
			return otto.Value{}
		}
		reg := cachedRegexp[call.ArgumentList[0].String()]
		if reg == nil {
			cachedRegexp[call.ArgumentList[0].String()] = regexp.MustCompile(call.ArgumentList[0].String())
			reg = cachedRegexp[call.ArgumentList[0].String()]
		}
		result, _ := cxt.VM.ToValue(reg.ReplaceAllString(reg.FindAllString(call.ArgumentList[1].String(), 1)[0], call.ArgumentList[2].String()))
		return result
	})

	cxt.VM.Set("MD5", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			logger.Error("script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList)))
			return otto.Value{}
		}
		result, _ := cxt.VM.ToValue(gox.Md5Sum(call.ArgumentList[0].String()))
		return result
	})

	cxt.VM.Set("TRIM", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			logger.Error("script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList)))
			return otto.Value{}
		}
		result, _ := cxt.VM.ToValue(strings.TrimSpace(call.ArgumentList[0].String()))
		return result
	})

	cxt.VM.Set("ENV", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			logger.Error("script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList)))
			return otto.Value{}
		}
		result, _ := cxt.VM.ToValue(cxt.ENV[call.ArgumentList[0].String()])
		return result
	})

	cxt.VM.Set("RESPONSE_DATA", func(call otto.FunctionCall) otto.Value {
		result, _ := cxt.VM.ToValue(cxt.ResponseData)
		return result
	})

	cxt.VM.Set("SET_RESPONSE_DATA", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			logger.Error("script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList)))
			return otto.Value{}
		}
		cxt.ResponseData = call.ArgumentList[0].String()
		return otto.Value{}
	})

	cxt.VM.Set("QUEUE", func(call otto.FunctionCall) otto.Value {
		result, _ := cxt.VM.ToValue(cxt.Queue)
		return result
	})

	cxt.VM.Set("MIDDLE_DATA", func(call otto.FunctionCall) otto.Value {
		result, _ := cxt.VM.ToValue(cxt.MiddleData)
		return result
	})

	cxt.VM.Set("ABS", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			logger.Error("script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList)))
			return otto.Value{}
		}
		relativeUrl := call.ArgumentList[0].String()
		ret, err := utils.AbsoluteURL(cxt.Queue.Url, relativeUrl)
		if err != nil {
			v, _ := cxt.VM.ToValue(relativeUrl)
			return v
		}
		v, _ := cxt.VM.ToValue(ret)
		return v
	})

	cxt.VM.Set("ADD_RESULT", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			logger.Error("script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList)))
			return otto.Value{}
		}
		k := call.ArgumentList[0].String()
		cxt.AddResult(&models.Result{
			TaskId: cxt.Queue.TaskId,
			Result: k,
		})
		return otto.Value{}
	})

	cxt.VM.Set("ADD_QUEUE", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			logger.Error("script Err: invalid arg number, expect 2, got " + convert.IntToStr(len(call.ArgumentList)))
			return otto.Value{}
		}
		k := call.ArgumentList[0].Object()
		if k == nil {
			v, _ := cxt.VM.ToValue(false)
			return v
		}
		urlValue, err := k.Get("url")
		if err != nil {
			v, _ := cxt.VM.ToValue(false)
			return v
		}
		stageValue, err := k.Get("stage")
		if err != nil {
			v, _ := cxt.VM.ToValue(false)
			return v
		}
		url := urlValue.String()
		stage := stageValue.String()
		if url == "" || stage == "" {
			v, _ := cxt.VM.ToValue(false)
			return v
		}
		taskId, _ := convert.StrToInt(cxt.ENV["taskId"])
		cxt.NewQueues[url] = &models.Queue{
			TaskId:    taskId,
			StageName: stage,
			Url:       url,
		}
		v, _ := cxt.VM.ToValue(true)
		return v
	})

	// AJAX(method, "http://url", headers, params, body)
	cxt.VM.Set("AJAX", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 5 {
			logger.Error("script Err: invalid arg number, expect 5, got " + convert.IntToStr(len(call.ArgumentList)))
			return otto.Value{}
		}

		method := strings.ToLower(call.ArgumentList[0].String())
		if method == "" {
			method = "get"
		}
		url := call.ArgumentList[1].String()
		body := call.ArgumentList[4].String()
		headerValue := call.ArgumentList[2].Object()
		headers := make(map[string]string)
		if headerValue != nil {
			for _, k := range headerValue.Keys() {
				value, _ := headerValue.Get(k)
				if value.String() != "" {
					headers[k] = value.String()
				}
			}
		}

		paramsValue := call.ArgumentList[3].Object()
		params := make(map[string]string)
		if paramsValue != nil {
			for _, k := range paramsValue.Keys() {
				value, _ := paramsValue.Get(k)
				if value.String() != "" {
					params[k] = value.String()
				}
			}
		}

		var resp *resty.Response
		var err error

		req := httpclient.GetClient(0, cxt.Project).
			R().
			SetHeaders(headers).
			SetQueryParams(params).
			SetBody(body)

		switch method {
		case "get":
			resp, err = req.Get(url)
		case "post":
			resp, err = req.Post(url)
		case "put":
			resp, err = req.Put(url)
		case "delete":
			resp, err = req.Delete(url)
		case "options":
			resp, err = req.Options(url)
		case "patch":
			resp, err = req.Patch(url)
		case "head":
			resp, err = req.Head(url)
		default:
			resp, err = req.Get(url)
		}

		result := make(map[string]interface{})
		if err != nil {
			result["error"] = err.Error()
		} else {
			cxt.HttpStatusCode = resp.StatusCode()
			result["status"] = resp.StatusCode()
			result["data"] = string(resp.Body())
		}
		ret, _ := cxt.VM.ToValue(result)
		return ret
	})

}
