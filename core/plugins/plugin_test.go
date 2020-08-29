///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package plugins_test

import (
	"digger/models"
	"digger/plugins"
	"fmt"
	"testing"
)

var cxt *models.Context

func init() {
	cxt = &models.Context{
		ResponseData: "{\"name\":\"张三\"}",
		ENV: map[string]string{
			"stage":     "ccc",
			"fieldName": "f1撒旦",
		},
	}
	plugins.InitVM(cxt)
}

func TestNewVM(t *testing.T) {
	fmt.Println(cxt.Exec(`
(function() {
	return ENV("fieldName");
})()
`))
}

func Test2(t *testing.T) {
	fmt.Println(cxt.Exec(`
(function() {
	return ENV("stage");
})()
`))
}

func Test3(t *testing.T) {
	fmt.Println(cxt.Exec(`
(function() {
	console.log(typeof CONTEXT_DATA("field1"))
	return CONTEXT_DATA("field1");
})()
`))
}

func Test4(t *testing.T) {
	fmt.Println(cxt.Exec(`
(function() {
	console.log(typeof CONTEXT_DATA("arr1"))
	console.log(CONTEXT_DATA("arr1")[0])
	console.log(CONTEXT_DATA("arr1")[1])
	console.log(CONTEXT_DATA("arr1")[2])
})()
`))
}

func Test5(t *testing.T) {
	fmt.Println(cxt.Exec(`
(function() {
	console.log(CONTEXT_DATA("project").Name)
	console.log(CONTEXT_DATA("project").DisplayName)
	console.log(CONTEXT_DATA("project").stages[0].ListCss)
})()
`))
}

func Test6(t *testing.T) {
	fmt.Println(cxt.Exec(`
(function() {
	var data = JSON.parse(CRAW_DATA());
	console.log(data.name)
})()
`))
}

func TestAjax1(t *testing.T) {
	fmt.Println(cxt.Exec(`
(function() {
	// AJAX(method, "http://url", headers, params, body) page_size=25&page=2
	var params = {
		page: 2,
		page_size: 25
	}
	var ret = AJAX("get", "https://hub.docker.com/v2/repositories/library/node/tags/", null, params, null)
	console.log(ret.err)
	console.log(ret.status)
	console.log(ret.data)
})()
`))
}

func TestAjax2(t *testing.T) {
	fmt.Println(cxt.Exec(`
(function() {
	// AJAX(method, "http://url", headers, params, body) page_size=25&page=2
	var params = {
		page: 2,
		page_size: 25
	}
	var ret = AJAX("post", "https://api.segment.io/v1/m", null, null, "{\"series\":[{\"type\":\"Counter\",\"metric\":\"analytics_js.integration.invoke\",\"value\":1,\"tags\":{\"method\":\"page\",\"integration_name\":\"Segment.io\"}}]}")
	console.log(ret.err)
	console.log(ret.status)
	console.log(ret.data)
})()
`))
}

func Test7(t *testing.T) {
	fmt.Println(cxt.Exec(`
(function() {
	console.log(SUBSTR("订了一份外卖", 0, 3))
})()
`))
}
