///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package models

import (
	"digger/common"
	"errors"
	"github.com/hetianyi/gox/convert"
	"regexp"
	"time"
)

// 项目
type Project struct {
	Id                 int               `json:"id" yaml:"-" gorm:"column:id;primary_key"`
	Name               string            `json:"name" yaml:"-" gorm:"column:name"`
	DisplayName        string            `json:"display_name" yaml:"-" gorm:"column:display_name"`
	Remark             string            `json:"remark" yaml:"-" gorm:"column:remark"`
	Tags               string            `json:"tags" yaml:"-" gorm:"column:tags"`
	StartUrl           string            `json:"start_url" yaml:"start_url" gorm:"column:start_url"`
	StartStage         string            `json:"start_stage" yaml:"start_stage" gorm:"column:start_stage"`
	Cron               string            `json:"cron" yaml:"-" gorm:"column:cron"`
	EnableCron         bool              `json:"enable_cron" yaml:"-" gorm:"column:enable_cron"`
	CreateTime         time.Time         `json:"create_time" yaml:"-" gorm:"column:create_time"`
	UpdateTime         time.Time         `json:"update_time" yaml:"-" gorm:"column:update_time"`
	Stages             []Stage           `json:"stages" yaml:"stages" gorm:"-"`
	Headers            map[string]string `json:"headers" yaml:"headers" gorm:"-"`
	HeadersDB          string            `json:"-" yaml:"-" gorm:"column:headers"`
	Settings           map[string]string `json:"settings" yaml:"settings" gorm:"-"`
	SettingsDB         string            `json:"-" yaml:"-" gorm:"column:settings"`
	NodeAffinity       []string          `json:"node_affinity" yaml:"node_affinity" gorm:"-"`
	NodeAffinityParsed []KV              `json:"node_affinity_parsed" yaml:"-" gorm:"-"`
	NodeAffinityDB     string            `json:"-" yaml:"-" gorm:"column:node_affinity"`
	// web接口附加额外信息
	Extras map[string]interface{} `json:"extras" yaml:"-" gorm:"-"`
}

func (Project) TableName() string {
	return "t_project"
}

func (p *Project) GetStageByName(stageName string) *Stage {
	le := len(p.Stages)
	if le == 0 {
		return nil
	}
	for _, s := range p.Stages {
		if s.Name == stageName {
			return &s
		}
	}
	return nil
}

func (p *Project) GetPluginByName(plugin string) *Stage {
	le := len(p.Stages)
	if le == 0 {
		return nil
	}
	for _, s := range p.Stages {
		if s.Name == plugin {
			return &s
		}
	}
	return nil
}

func (p *Project) GetIntSetting(name string) int {
	i, _ := convert.StrToInt(p.Settings[name])
	return i
}

func (p *Project) GetStrSetting(name string) string {
	return p.Settings[name]
}

func (p *Project) GetBoolSetting(name string) bool {
	i, _ := convert.StrToBool(p.Settings[name])
	return i
}

func (p *Project) Validate() (bool, error) {
	m, _ := regexp.MatchString(common.NAME_REGEXP, p.Name)
	if p.Name == "" || !m {
		return false, errors.New("invalid name: \"" + p.Name + "\"")
	}
	if p.DisplayName == "" {
		return false, errors.New("invalid display_name: \"" + p.DisplayName + "\"")
	}
	if p.StartUrl == "" {
		return false, errors.New("invalid start_url: \"" + p.StartUrl + "\"")
	}
	if len(p.Stages) == 0 {
		return false, errors.New("no stage")
	}
	for i, stage := range p.Stages {
		m, _ := regexp.MatchString(common.NAME_REGEXP, stage.Name)
		if stage.Name == "" || !m {
			return false, errors.New("invalid stage name: \"" + stage.Name + "\" @ index " + convert.IntToStr(i))
		}
		if len(stage.Fields) == 0 {
			return false, errors.New("stage \"" + stage.Name + "\" has no field")
		}
		for i, field := range stage.Fields {
			m, _ := regexp.MatchString(common.NAME_REGEXP, field.Name)
			if field.Name == "" || !m {
				return false, errors.New("invalid field name: \"" + field.Name + "\" @ index " + convert.IntToStr(i))
			}
			if field.NextStage != "" && p.GetStageByName(field.NextStage) == nil {
				return false, errors.New("next stage \"" + field.NextStage + "\" not found @ field \"" + field.Name + "\"")
			}
			if field.PluginDB != "" && field.Plugin == nil {
				return false, errors.New("plugin \"" + field.PluginDB + "\" not found @ field \"" + field.Name + "\"")
			}
		}
	}
	return true, nil
}

// 阶段
type Stage struct {
	Id           int      `json:"id" yaml:"-" gorm:"column:id;primary_key"`
	ProjectId    int      `json:"project_id" yaml:"-" gorm:"column:project_id"`
	Name         string   `json:"name" yaml:"name" gorm:"column:name"`
	IsList       bool     `json:"is_list" yaml:"is_list" gorm:"column:is_list"`
	IsUnique     bool     `json:"is_unique" yaml:"is_unique" gorm:"column:is_unique"`
	ListXpath    string   `json:"list_xpath" yaml:"list_xpath" gorm:"column:list_xpath"`
	ListCss      string   `json:"list_css" yaml:"list_css" gorm:"column:list_css"`
	PageXpath    string   `json:"page_xpath" yaml:"page_xpath" gorm:"column:page_xpath"`
	PageCss      string   `json:"page_css" yaml:"page_css" gorm:"column:page_css"`
	PageAttr     string   `json:"page_attr" yaml:"page_attr" gorm:"column:page_attr"`
	PluginsDB    string   `json:"plugins" yaml:"plugin" gorm:"column:plugins"`
	Plugins      []Plugin `json:"plugin_array" yaml:"-" gorm:"-"`
	Fields       []Field  `json:"fields" yaml:"fields" gorm:"-"`
	HasNextStage bool     `json:"has_next_stage" yaml:"-" gorm:"-"`
}

func (Stage) TableName() string {
	return "t_stage"
}

func (s Stage) FindPlugins(slot string) *Plugin {
	if len(s.Plugins) == 0 {
		return nil
	}
	for _, p := range s.Plugins {
		if p.Slot == slot {
			return &p
		}
	}
	return nil
}

// 字段
type Field struct {
	Id        int     `json:"id" yaml:"-" gorm:"column:id;primary_key"`
	ProjectId int     `json:"project_id" yaml:"-" gorm:"column:project_id"`
	StageId   int     `json:"stage_id" yaml:"-" gorm:"column:stage_id"`
	Name      string  `json:"name" yaml:"name" gorm:"column:name"`
	IsArray   bool    `json:"is_array" yaml:"is_array" gorm:"column:is_array"`
	IsHtml    bool    `json:"is_html" yaml:"is_html" gorm:"column:is_html"`
	Xpath     string  `json:"xpath" yaml:"xpath" gorm:"column:xpath"`
	Css       string  `json:"css" yaml:"css" gorm:"column:css"`
	Attr      string  `json:"attr" yaml:"attr" gorm:"column:attr"`
	PluginDB  string  `json:"plugin" yaml:"plugin" gorm:"column:plugin"`
	Plugin    *Plugin `json:"_plugin" yaml:"-" gorm:"-"`
	Remark    string  `json:"remark" yaml:"remark" gorm:"column:remark"`
	NextStage string  `json:"next_stage" yaml:"next_stage" gorm:"column:next_stage"`
}

func (Field) TableName() string {
	return "t_field"
}

type Task struct {
	Id               int       `json:"id" gorm:"column:id;primary_key"`
	ProjectId        int       `json:"project_id" gorm:"column:project_id"`
	ConfigSnapShotId int       `json:"config_snapshot_id" gorm:"column:config_snapshot_id"`
	Status           int       `json:"status" gorm:"column:status"`
	ResultCount      int       `json:"result_count" gorm:"column:result_count"`
	IOIn             int64     `json:"io_in" gorm:"column:io_in"`
	IOOut            int64     `json:"io_out" gorm:"column:io_out"`
	SuccessRequest   int       `json:"success_request" gorm:"column:success_request"`
	ErrorRequest     int       `json:"error_request" gorm:"column:error_request"`
	BindNodeMode     int       `json:"bind_node_mode" gorm:"column:bind_node_mode"`
	CreateTime       time.Time `json:"create_time" gorm:"column:create_time"`
}

func (Task) TableName() string {
	return "t_task"
}

type ConfigSnapshot struct {
	Id        int    `json:"id" gorm:"column:id;primary_key"`
	ProjectId int    `json:"project_id" gorm:"column:project_id"`
	Config    string `json:"config" gorm:"column:config"`
}

func (ConfigSnapshot) TableName() string {
	return "t_config_snapshot"
}

type Result struct {
	Id     int64  `json:"id" gorm:"column:id;primary_key"`
	TaskId int    `json:"task_id" gorm:"column:task_id"`
	Result string `json:"result" gorm:"column:result"`
}

func (Result) TableName() string {
	return "t_result"
}

// 调度任务队列
type Queue struct {
	Id         int64  `json:"id" gorm:"column:id;primary_key"`
	TaskId     int    `json:"task_id" gorm:"column:task_id"`
	StageName  string `json:"stage_name" gorm:"column:stage_name"`
	Url        string `json:"url" gorm:"column:url"`
	MiddleData string `json:"middle_data" gorm:"column:middle_data"`
	Expire     int64  `json:"expire" gorm:"column:expire"`
}

func (Queue) TableName() string {
	return "t_queue"
}

// 调度任务队列
type Plugin struct {
	Id        int    `json:"id" gorm:"column:id;primary_key"`
	Name      string `json:"name" gorm:"column:name"`
	Script    string `json:"script" gorm:"column:script"`
	Slot      string `json:"slot" gorm:"-"`
	ProjectId int    `json:"project_id" gorm:"column:project_id"`
}

func (Plugin) TableName() string {
	return "t_plugin"
}

type User struct {
	Id       int    `json:"id" gorm:"column:id;primary_key"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}

func (User) TableName() string {
	return "t_user"
}

type Config struct {
	Key   string `json:"key" gorm:"column:key;primary_key"`
	Value string `json:"value" gorm:"column:value"`
}

func (Config) TableName() string {
	return "t_config"
}

type Statistic struct {
	Id         int       `json:"id" gorm:"column:id;primary_key"`
	Data       string    `json:"data" gorm:"column:data"`
	CreateTime time.Time `json:"create_time" yaml:"-" gorm:"column:create_time"`
}

func (Statistic) TableName() string {
	return "t_statistic"
}
