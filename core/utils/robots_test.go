///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils_test

import (
	"digger/models"
	"digger/utils"
	"fmt"
	"github.com/hetianyi/gox/logger"
	"github.com/temoto/robotstxt"
	"net/url"
	"testing"
)

func TestLoadRobotsTxt(t *testing.T) {
	bytes, err := utils.LoadRobotsTxt("http://www.zhanxixi.com/robots.txt", &models.Project{})
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(string(bytes))

	r1, err := robotstxt.FromBytes(bytes)

	u1 := "/forum-41-1.html"
	u2 := "/api/"
	u3 := "/api/v1/upload"
	u4 := "/api/v1/upload?name=zhangsan"
	fmt.Println(r1.TestAgent(u1, "Bot"))
	fmt.Println(r1.TestAgent(u2, "Bot"))
	fmt.Println(r1.TestAgent(u3, "Bot"))
	fmt.Println(r1.TestAgent(u4, "Bot"))

	fmt.Println("\n\n-------------------------------------------------")
	bytes, err = utils.LoadRobotsTxt("https://news.163.com/robots.txt", &models.Project{})
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(string(bytes))
}

func TestUrl(t *testing.T) {
	u1 := "http://www.zhanxixi.com/api/v1/upload?name=zhangsan"
	base, _ := url.Parse(u1)
	fmt.Println(base.Path)
	fmt.Println(base.RawQuery)
	fmt.Println(base.RawPath)
	fmt.Println(base.String())

	fmt.Println(url.Parse("1"))
}
