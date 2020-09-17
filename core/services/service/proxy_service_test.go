///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package service_test

import (
	"digger/models"
	"digger/services/service"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"testing"
)

func TestProxyServiceImp_Save(t *testing.T) {
	fmt.Println(service.ProxyService().Save(models.Proxy{
		Address:        "127.0.0.1:1080",
		Remark:         "this is a local proxy1",
		ProxyGenScript: "",
	}))
	fmt.Println(service.ProxyService().Save(models.Proxy{
		Address:        "127.0.0.1:1081",
		Remark:         "this is a local proxy2",
		ProxyGenScript: "",
	}))
}

func TestProxyServiceImp_List(t *testing.T) {
	fmt.Println(service.ProxyService().List(&models.ProxyQueryVO{
		PageQueryVO: models.PageQueryVO{
			Page:     1,
			PageSize: 10,
		},
		Key: "proxy1",
	}))
}

func TestProxyServiceImp_Delete(t *testing.T) {
	service.ProxyService().Delete([]int{5})
}

func TestProxyServiceImp_SelectByProject(t *testing.T) {
	proxies, _ := service.ProxyService().SelectByProject(14)
	fmt.Println(jsoniter.MarshalToString(proxies))
}
