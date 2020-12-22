///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils_test

import (
	"digger/utils"
	"fmt"
	"github.com/go-resty/resty/v2"
	"testing"
)

func TestTrans2UTF8(t *testing.T) {
	httpClient := resty.New()
	resp, _ := httpClient.R().
		Get("https://www.tui78.com/yimiao/85738.html")

	fmt.Println(utils.Trans2UTF8("gbk", string(resp.Body())))
}
