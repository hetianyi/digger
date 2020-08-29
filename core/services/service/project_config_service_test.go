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
	"github.com/hetianyi/gox/logger"
	json "github.com/json-iterator/go"
	"testing"
)

func TestProjectConfigServiceImp_SelectStages(t *testing.T) {
	array, err := service.ProjectConfigService().SelectStages(1)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(len(array))
	s, _ := json.MarshalToString(array)
	fmt.Println(s)
}

func TestProjectConfigServiceImp_SelectFields(t *testing.T) {
	array, err := service.ProjectConfigService().SelectFields(1)
	if err != nil {
		logger.Fatal(err)
	}
	fmt.Println(len(array))
	s, _ := json.MarshalToString(array)
	fmt.Println(s)
}

func TestProjectConfigServiceImp_SaveStagesAndFields(t *testing.T) {
	var stages = []models.Stage{
		{
			Name:     "阶段3",
			IsList:   true,
			ListCss:  "xxx",
			PageCss:  "xxxxx",
			PageAttr: "asdasdad",
			Fields: []models.Field{
				{
					Name:      "撒旦1",
					IsArray:   false,
					IsHtml:    true,
					Css:       "s速度",
					Attr:      "阿萨德岁的",
					Remark:    "asd速度",
					NextStage: "阶段2",
				},
				{
					Name:      "撒旦2",
					IsArray:   false,
					IsHtml:    true,
					Css:       "s速度aasd",
					Attr:      "阿萨德岁as的",
					Remark:    "asdasdsa速度",
					NextStage: "",
				},
			},
		},
		{
			Name:     "阶段4",
			IsList:   false,
			ListCss:  "xxx",
			PageCss:  "xxxxx",
			PageAttr: "asdasdad",
			Fields: []models.Field{
				{
					Name:      "撒旦3",
					IsArray:   false,
					IsHtml:    true,
					Css:       "s速度",
					Attr:      "阿萨德岁的",
					Remark:    "asd速度",
					NextStage: "阶段2",
				},
				{
					Name:      "撒旦4",
					IsArray:   false,
					IsHtml:    true,
					Css:       "s速度aasd",
					Attr:      "阿萨德岁as的",
					Remark:    "asdasdsa速度",
					NextStage: "",
				},
			},
		},
	}
	err := service.ProjectConfigService().SaveStagesAndFields(18, stages)
	if err != nil {
		logger.Fatal(err)
	}
}
