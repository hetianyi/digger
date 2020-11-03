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
	"github.com/hetianyi/gox/logger"
	"net/http"
	"testing"
)

func TestDetect(t *testing.T) {
	httpClient := resty.New()
	url := "https://music.163.com/discover/playlist"
	resp, err := httpClient.R().Get(url)

	if err != nil {
		logger.Fatal(err)
	}

	if resp.StatusCode() != http.StatusOK {
		logger.Fatal(fmt.Sprintf("error http status: %d", resp.StatusCode()))
	}
	body := resp.Body()

	utils.DetectStructure(string(body))
}
