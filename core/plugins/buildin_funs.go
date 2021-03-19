///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package plugins

import (
	"bytes"
	"crypto/tls"
	"digger/httpclient"
	"digger/models"
	"digger/utils"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/go-resty/resty/v2"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
	uuid "github.com/hetianyi/gox/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/robertkrimen/otto"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	cachedRegexp = make(map[string]*regexp.Regexp)
	httpClient   = resty.New()
	json         = jsoniter.ConfigFastest
	uploadClient *resty.Client
)

func init() {
	uploadClient = resty.New().
		SetTLSClientConfig(&tls.Config{
			InsecureSkipVerify: true,
		})
}

func InitVM(cxt *models.Context) {
	cxt.VM = otto.New()
	initBuildInFunctions(cxt)
}

func initBuildInFunctions(cxt *models.Context) {

	cxt.VM.Set("LEN", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			errMsg := "script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(0)
			return result
		}
		result, _ := cxt.VM.ToValue(len([]rune(call.ArgumentList[0].String())))
		return result
	})

	cxt.VM.Set("STARTS_WITH", func(call otto.FunctionCall) otto.Value {
		boo := false
		if len(call.ArgumentList) != 2 {
			errMsg := "script Err: invalid arg number, expect 2, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(boo)
			return result
		}
		result, _ := cxt.VM.ToValue(strings.HasPrefix(call.ArgumentList[0].String(), call.ArgumentList[1].String()))
		return result
	})

	cxt.VM.Set("ENDS_WITH", func(call otto.FunctionCall) otto.Value {
		boo := false
		if len(call.ArgumentList) != 2 {
			errMsg := "script Err: invalid arg number, expect 2, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(boo)
			return result
		}
		result, _ := cxt.VM.ToValue(strings.HasSuffix(call.ArgumentList[0].String(), call.ArgumentList[1].String()))
		return result
	})

	cxt.VM.Set("SUBSTR", func(call otto.FunctionCall) otto.Value {
		boo := false
		if len(call.ArgumentList) != 3 {
			errMsg := "script Err: invalid arg number, expect 3, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
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
			errMsg := "script Err: invalid arg number, expect 2, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(boo)
			return result
		}
		result, _ := cxt.VM.ToValue(strings.Contains(call.ArgumentList[0].String(), call.ArgumentList[1].String()))
		return result
	})

	cxt.VM.Set("REPLACE", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 3 {
			errMsg := "script Err: invalid arg number, expect 3, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
		}
		result, _ := cxt.VM.ToValue(strings.Replace(call.ArgumentList[0].String(), call.ArgumentList[1].String(), call.ArgumentList[2].String(), -1))
		return result
	})

	cxt.VM.Set("REGEXP_GROUP_FIND", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 3 {
			errMsg := "script Err: invalid arg number, expect 3, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
		}
		reg := cachedRegexp[call.ArgumentList[0].String()]
		if reg == nil {
			cachedRegexp[call.ArgumentList[0].String()] = regexp.MustCompile(call.ArgumentList[0].String())
			reg = cachedRegexp[call.ArgumentList[0].String()]
		}

		result, _ := cxt.VM.ToValue("")
		gox.Try(func() {
			result, _ = cxt.VM.ToValue(reg.ReplaceAllString(reg.FindAllString(call.ArgumentList[1].String(), 1)[0], call.ArgumentList[2].String()))
		}, func(err interface{}) {
			cxt.Log.Write([]byte(err.(error).Error()))
		})
		return result
	})

	cxt.VM.Set("MD5", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			errMsg := "script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
		}
		result, _ := cxt.VM.ToValue(gox.Md5Sum(call.ArgumentList[0].String()))
		return result
	})

	cxt.VM.Set("TRIM", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			errMsg := "script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
		}
		result, _ := cxt.VM.ToValue(strings.TrimSpace(call.ArgumentList[0].String()))
		return result
	})

	cxt.VM.Set("ENV", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			errMsg := "script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
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
			errMsg := "script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
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
			errMsg := "script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
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
			errMsg := "script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
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
			errMsg := "script Err: invalid arg number, expect 2, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
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
		middleDataValue, err := k.Get("middle_data")
		if err != nil {
			v, _ := cxt.VM.ToValue(false)
			return v
		}
		middleDataObj := middleDataValue.Object()
		middleData := make(map[string]string)
		if middleDataObj != nil {
			for _, k := range middleDataObj.Keys() {
				value, _ := middleDataObj.Get(k)
				if value.String() != "" {
					middleData[k] = value.String()
				}
			}
		}
		middleDataString, _ := json.Marshal(middleData)

		url := urlValue.String()
		stage := stageValue.String()
		if url == "" || stage == "" {
			v, _ := cxt.VM.ToValue(false)
			return v
		}
		taskId, _ := convert.StrToInt(cxt.ENV["taskId"])
		cxt.NewQueues = append(cxt.NewQueues, &models.Queue{
			TaskId:     taskId,
			StageName:  stage,
			Url:        url,
			MiddleData: string(middleDataString),
		})
		v, _ := cxt.VM.ToValue(true)
		return v
	})

	// AJAX(method, "http://url", headers, params, body)
	cxt.VM.Set("AJAX", func(call otto.FunctionCall) otto.Value {
		result := make(map[string]interface{})
		if len(call.ArgumentList) != 5 {
			errMsg := "script Err: invalid arg number, expect 5, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result["status"] = 0
			result["error"] = errMsg
			ret, _ := cxt.VM.ToValue(result)
			return ret
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

		parsedUrl, err := utils.Parse(url)
		if err != nil {
			result["error"] = err.Error()
			result["status"] = 0
			ret, _ := cxt.VM.ToValue(result)
			return ret
		}
		client := httpclient.GetClient(0, cxt.Project)
		feedback := utils.TryProxy(parsedUrl.Scheme, client, cxt.Queue.TaskId, cxt)
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
		if resp != nil && cxt.PlayResult != nil {
			cxt.PlayResult.HttpStatus = resp.StatusCode()
			cxt.PlayResult.HttpResult = string(resp.Body())
			cxt.HttpResponseHeaders = resp.Header()
		}

		// feedback
		if feedback != nil {
			if err != nil || resp.StatusCode() != http.StatusOK {
				feedback.Fail()
			} else {
				feedback.Success()
			}
		}
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

	cxt.VM.Set("FROM_JSON", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			errMsg := "script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(nil)
			return result
		}
		var r interface{}
		err := json.UnmarshalFromString(call.ArgumentList[0].String(), &r)
		if err != nil {
			errMsg := fmt.Sprintf("[%d:%d]", cxt.Queue.TaskId, cxt.Queue.Id) + "cannot parse json: " + err.Error()
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(nil)
			return result
		}
		result, _ := cxt.VM.ToValue(r)
		return result
	})

	cxt.VM.Set("TO_JSON", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 1 {
			errMsg := "script Err: invalid arg number, expect 1, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(nil)
			return result
		}
		obj, err := call.ArgumentList[0].Export()
		if err != nil {
			errMsg := "cannot format json: " + err.Error()
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(nil)
			return result
		}
		str, err := json.MarshalToString(obj)
		if err != nil {
			errMsg := "cannot format json: " + err.Error()
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(nil)
			return result
		}
		result, _ := cxt.VM.ToValue(str)
		return result
	})

	cxt.VM.Set("LOG", func(call otto.FunctionCall) otto.Value {
		var buff bytes.Buffer
		for _, v := range call.ArgumentList {
			buff.WriteString(v.String())
		}
		cxt.Log.Write(buff.Bytes())
		return otto.Value{}
	})

	cxt.VM.Set("LOGF", func(call otto.FunctionCall) otto.Value {
		var format string
		var args []interface{}
		for i, v := range call.ArgumentList {
			if i == 0 {
				format = v.String()
			} else {
				args = append(args, v.String())
			}
		}
		log := fmt.Sprintf(format, args...)
		cxt.Log.Write([]byte(log))
		return otto.Value{}
	})

	cxt.VM.Set("XPATH_FIND", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 {
			errMsg := "script Err: invalid arg number, expect 2, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(nil)
			return result
		}

		content := fmt.Sprintf("<html><body><div>%s</div></body></html>", call.ArgumentList[0].String())
		xpath := call.ArgumentList[1].String()

		doc, err := parseXpathDocument(content)
		if err != nil {
			errMsg := "cannot parse xpath document: " + err.Error()
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(nil)
			return result
		}

		list, err := htmlquery.QueryAll(doc, xpath)
		if err != nil {
			errMsg := "cannot query xpath document: " + err.Error()
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(nil)
			return result
		}
		var ret []string
		for _, item := range list {
			ret = append(ret, htmlquery.InnerText(item))
		}

		result, _ := cxt.VM.ToValue(ret)
		return result
	})

	// POST https://xxx.com headers params fileUrl
	cxt.VM.Set("UPLOAD", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 5 {
			errMsg := "script Err: invalid arg number, expect 5, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
		}

		method := strings.ToLower(call.ArgumentList[0].String())
		if method == "" {
			method = "post"
		}
		uploadUrl := call.ArgumentList[1].String()
		headerValue := call.ArgumentList[2].Object()
		paramValue := call.ArgumentList[3].Object()
		fileUrl := call.ArgumentList[4].String()
		headers := make(map[string]string)
		if headerValue != nil {
			for _, k := range headerValue.Keys() {
				value, _ := headerValue.Get(k)
				if value.String() != "" {
					headers[k] = value.String()
				}
			}
		}
		params := make(map[string]string)
		if paramValue != nil {
			for _, k := range paramValue.Keys() {
				value, _ := paramValue.Get(k)
				if value.String() != "" {
					params[k] = value.String()
				}
			}
		}
		params["projectName"] = cxt.Project.Name

		parsedUrl, err := utils.Parse(fileUrl)
		if err != nil {
			errMsg := "cannot parse fileUrl: " + err.Error()
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
		}
		/*upUrl, err := utils.Parse(uploadUrl)
		if err != nil {
			ret, _ := cxt.VM.ToValue(nil)
			return ret
		}*/
		//pageUrl, _ := utils.Parse(cxt.Queue.Url)
		downHost := parsedUrl.Host
		/*if pageUrl != nil {
			downHost = pageUrl.Host
		} else {
			downHost = parsedUrl.Host
		}*/
		if len(headers) == 0 {
			headers = map[string]string{
				"Host":       downHost,
				"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
				"Refer":      cxt.Queue.Url,
			}
		}
		client := httpclient.GetClient(0, cxt.Project)
		feedback := utils.TryProxy(parsedUrl.Scheme, client, cxt.Queue.TaskId, cxt)
		req := httpclient.GetClient(0, cxt.Project).
			R().
			SetDoNotParseResponse(true).
			SetHeaders(headers)

		resp, err := req.Get(fileUrl)
		if err != nil {
			errMsg := "request error for url: " + err.Error()
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
		}

		// feedback
		if feedback != nil {
			if err != nil || resp.StatusCode() != http.StatusOK {
				feedback.Fail()
			} else {
				feedback.Success()
			}
		}
		if resp.StatusCode() != http.StatusOK {
			errMsg := "error download resource: server response http status " + convert.IntToStr(resp.StatusCode())
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
		}
		// 下载文件
		tempFile, err := file.CreateFile(os.TempDir() + "/" + uuid.UUID())
		if err != nil {
			errMsg := fmt.Sprintf("临时文件创建失败: %s", err.Error())
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
		}
		defer func() {
			tempFile.Close()
			if err = os.Remove(tempFile.Name()); err != nil {
				cxt.Log.Write([]byte(fmt.Sprintf("临时文件删除失败: %s", err.Error())))
			}
		}()
		downloadSize, err := io.Copy(tempFile, resp.RawBody())
		if err != nil {
			errMsg := fmt.Sprintf("资源下载失败: %s", err.Error())
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue("")
			return result
		}
		tempFile.Seek(0, 0)
		md5, err := file.GetFileMd5(tempFile.Name())
		if err != nil {
			cxt.Log.Write([]byte(fmt.Sprintf("计算md5失败: %s", err.Error())))
		}
		if md5 != "" {
			params["md5"] = md5
		}

		tempFile.Seek(0, 0)
		cxt.Log.Write([]byte(fmt.Sprintf("资源下载成功，大小: %dkb", downloadSize/1024)))

		/*tempReadFile, err := file.GetFile(tempFile.Name())
		if err != nil {
			cxt.Log.Write([]byte(fmt.Sprintf("临时文件打开失败: %s", err.Error())))
			return otto.Value{}
		}*/

		upReq := uploadClient.
			R().
			SetHeaders(map[string]string{
				"Host":       downHost,
				"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
				"Refer":      cxt.Queue.Url,
			}).
			SetContentLength(true).
			SetQueryParams(params).
			SetBody(tempFile)

		var upResp *resty.Response
		switch method {
		case "get":
			upResp, err = upReq.Get(uploadUrl)
		case "post":
			upResp, err = upReq.Post(uploadUrl)
		case "put":
			upResp, err = upReq.Put(uploadUrl)
		case "delete":
			upResp, err = upReq.Delete(uploadUrl)
		case "options":
			upResp, err = upReq.Options(uploadUrl)
		case "patch":
			upResp, err = upReq.Patch(uploadUrl)
		case "head":
			upResp, err = upReq.Head(uploadUrl)
		default:
			upResp, err = upReq.Get(uploadUrl)
		}
		result, _ := cxt.VM.ToValue(string(upResp.Body()))
		return result
	})

	cxt.VM.Set("FORMAT", func(call otto.FunctionCall) otto.Value {
		if len(call.ArgumentList) != 2 {
			errMsg := "script Err: invalid arg number, expect 2, got " + convert.IntToStr(len(call.ArgumentList))
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
			result, _ := cxt.VM.ToValue(nil)
			return result
		}
		formatString := call.ArgumentList[0].String()
		var argsString []interface{}
		if len(call.ArgumentList) > 1 {
			for i, s := range call.ArgumentList {
				if i == 0 {
					continue
				}
				argsString = append(argsString, s.String())
			}
		}

		var ret = ""
		gox.Try(func() {
			ret = fmt.Sprintf(formatString, argsString...)
		}, func(e interface{}) {
			errMsg := e.(error).Error()
			logger.Error(errMsg)
			cxt.Log.Write([]byte(errMsg))
		})

		result, _ := cxt.VM.ToValue(ret)
		return result
	})
}

// 用goquery解析html文档
// reParse: 丢弃旧的重新使用cxt.ResponseData解析
func parseXpathDocument(content string) (*html.Node, error) {
	return htmlquery.Parse(bytes.NewBuffer([]byte(content)))
}
