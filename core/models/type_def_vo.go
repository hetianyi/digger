///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package models

import (
	"digger/common"
	"github.com/PuerkitoBio/goquery"
	"github.com/dgrijalva/jwt-go"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/robertkrimen/otto"
	"golang.org/x/net/html"
	"io"
	"time"
)

type BootstrapConfig struct {
	BootMode    common.ROLE
	Port        int
	InstanceId  int
	LogDir      string
	LogLevel    string
	Secret      string
	ManagerUrl  string
	DBString    string
	RedisString string
	Labels      map[string]string
	UIDir       string
}

type RedisConfig struct {
	Address  string
	Password string
	DB       int
}

type EmailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (e *EmailConfig) String() string {
	return e.Username + ":" + e.Password + "@" + e.Host + ":" + convert.IntToStr(e.Port)
}

type PageQueryVO struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type Resp struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type ProjectQueryVO struct {
	PageQueryVO
	Name        string   `json:"name"`
	DisplayName string   `json:"display_name"`
	Tags        []string `json:"tags"`
	Order       int      `json:"order"`
}

type ResultQueryVO struct {
	PageQueryVO
	TaskId       int   `json:"task_id"`
	LastResultId int64 `json:"last_result_id"`
}

type TaskQueryVO struct {
	PageQueryVO
	ProjectId int `json:"project_id"`
	Status    int `json:"status"`
}

type QueueQueryVO struct {
	TaskId             int  `json:"task_id"`
	Status             int  `json:"status"`
	Limit              int  `json:"limit"`
	QueueExpireSeconds int  `json:"queue_expire_seconds"`
	LockStatus         bool `json:"lock"` // 查询之后是否立即锁定状态
}

type FetchQueueResponseVO struct {
	Check  string   `json:"check"`
	Queues []*Queue `json:"queues"`
}

type QueueCallbackRequestVO struct {
	Check               string    `json:"check"`
	SuccessQueueIds     []int64   `json:"success_queue_ids"`
	SuccessQueueTaskIds []int     `json:"success_queue_task_ids"`
	ErrorQueueIds       []int64   `json:"error_queue_ids"`      // 错误的queue，会重试
	ErrorQueueTaskIds   []int     `json:"error_queue_task_ids"` // 错误的queue，会重试
	NewQueues           []*Queue  `json:"new_queues"`
	Results             []*Result `json:"results"`
}

type PlayInputVO1 struct {
	ProjectId int    `json:"project_id"`
	StageName string `json:"stage_name"`
	Url       string `json:"url"`
}

type PlayInputVO2 struct {
	StageName string `json:"stage_name"`
	Url       string `json:"url"`
	Project   string `json:"project"`
	ProjectId int    `json:"project_id"`
}

type PlayOutputVO struct {
	ProjectId int       `json:"projectId"`
	StageName string    `json:"stage_name"`
	Url       string    `json:"url"`
	Next      []*Queue  `json:"next"`
	Result    []*Result `json:"result"`
}

// TODO
type Context struct {
	Project        *Project
	Stage          *Stage
	Queue          *Queue
	HttpStatusCode int
	ResponseData   string
	Log            io.Writer
	// stage
	// taskId
	//
	ENV           map[string]string
	NewQueues     map[string]*Queue
	MiddleData    map[string]string
	Results       []*Result
	VM            *otto.Otto
	CssQueryDoc   *goquery.Document
	XpathQueryDoc *html.Node
}

func (c *Context) Exec(script string) (string, error) {
	//return v.vm.Run()
	value, err := c.VM.Run(script)
	if err != nil {
		return "", err
	}
	return value.String(), nil
}

func (c *Context) AddQueue(q *Queue) {
	c.NewQueues[q.Url] = q
}

func (c *Context) AddResult(r *Result) {
	c.Results = append(c.Results, r)
}

type Slot interface {
}

type ResultCountCO struct {
	TaskId int `json:"task_id" gorm:"column:task_id"`
	Count  int `json:"count" gorm:"column:count"`
}

type TaskCountCO struct {
	ProjectId   int `json:"projectId" gorm:"column:project_id"`
	ActiveCount int `json:"active_count" gorm:"column:active_count"`
	PauseCount  int `json:"pause_count" gorm:"column:pause_count"`
	StopCount   int `json:"stop_count" gorm:"column:stop_count"`
	FinishCount int `json:"finish_count" gorm:"column:finish_count"`
}

type RedisEvent struct {
	Event int               `json:"event"`
	Body  map[string]string `json:"body"`
}

type DispatchWork struct {
	RequestId string `json:"requestId"`
	Queue     *Queue `json:"queue"`
}

type QueueProcessResult struct {
	TaskId    int       `json:"taskId"`
	InitUrl   string    `json:"init_url"`
	QueueId   int64     `json:"queueId"`
	Expire    int64     `json:"expire"`
	RequestId string    `json:"requestId"`
	Error     string    `json:"error"`
	Logs      string    `json:"logs"`
	NewQueues []*Queue  `json:"newQueues"`
	Results   []*Result `json:"results"`
}

type RestResponse struct {
	Resp
	Data interface{} `json:"data"`
}

type Node struct {
	InstanceId int               `json:"instance_id"`
	RegisterAt string            `json:"register_at"`
	Address    string            `json:"address"`
	Status     int               `json:"status"`
	Down       int               `json:"down"`
	Assign     int64             `json:"assign"`
	Success    int64             `json:"success"`
	Error      int64             `json:"error"`
	Labels     map[string]string `json:"labels"`
}

type DynamicQueryVO struct {
	PageQueryVO
	Tag      string `json:"tag"`
	Username string `json:"username"`
	Key      string `json:"key"`
}

type MyCustomClaims struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

type KV struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type StatisticVO struct {
	Id         int                    `json:"id"`
	Data       map[string]interface{} `json:"data"`
	CreateTime time.Time              `json:"create_time"`
}

func (s StatisticVO) From(d *Statistic) *StatisticVO {
	data := make(map[string]interface{})
	err := jsoniter.UnmarshalFromString(d.Data, &data)
	if err != nil {
		logger.Error(err)
		return &s
	}
	s.Id = d.Id
	s.Data = data
	s.CreateTime = d.CreateTime
	return &s
}

type ProxyQueryVO struct {
	PageQueryVO
	Key string `json:"key"`
}
