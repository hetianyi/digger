///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package common

type ROLE int
type SLOT string
type EVENT int

const (
	VERSION           = "0.0.1"
	ROLE_NONE    ROLE = 0
	ROLE_MANAGER ROLE = 1
	ROLE_WORKER  ROLE = 2
	// slot
	s1            SLOT = "s1" // 请求之前，输入context等，输出新的url
	sr            SLOT = "sr" // 自定义请求，返回请求结果（例如post请求）
	s2            SLOT = "s2" // 请求之后，处理之前，输入为：http请求返回内容，返回值：处理后的内容和处理引擎指向
	s3            SLOT = "s3" // TODO 处理中，处理使用的引擎，默认为goquery，否则为自定义的插件处理引擎，自定义引擎需要自行匹配field值以及next stage
	s4            SLOT = "s4" // 引擎为goquery时，解析得到字段值之后，可以用来修正数据，如去空格，剪切等
	DefaultSecret      = "123456"
	LOGIN_USER         = "LOGIN_USER"
	MAX_RETRY          = 3

	PAUSE   = 0
	RUNNING = 1
	STOP    = 2

	EV_TASK_CREATED     = 1
	EV_TASK_PAUSE       = 2
	EV_TASK_CONTINUE    = 3
	EV_TASK_STOP        = 4
	EV_ONE_QUEUE_FINISH = 5

	// 并发请求数
	SETTINGS_CONCURRENT_REQUESTS  = "CONCURRENT_REQUESTS"
	SETTINGS_QUEUE_EXPIRE_SECONDS = "QUEUE_EXPIRE_SECONDS"
	// 是否跟随重定向
	SETTINGS_FOLLOW_REDIRECT = "FOLLOW_REDIRECT"
	// 请求超时时间(s)
	SETTINGS_REQUEST_TIMEOUT = "REQUEST_TIMEOUT"
	// 重试次数（单节点，非全局）
	SETTINGS_RETRY_COUNT = "RETRY_COUNT"
	// 重试间隔时间(s)
	SETTINGS_RETRY_WAIT = "RETRY_WAIT"
	// 是否跳过tls验证，解决自谦证书问题
	SETTINGS_SKIP_TLS_VERIFY = "SKIP_TLS_VERIFY"
	// 导出时每次从数据库查询的分页大小，影响导出速度和内存占用
	SETTINGS_EXPORT_PAGE_SIZE = "EXPORT_PAGE_SIZE"
	// 是否遵循robots指令
	SETTINGS_FOLLOW_ROBOTS_TXT = "FOLLOW_ROBOTS_TXT"

	NAME_REGEXP = "^[a-zA-Z0-9_]+$"
	//REPO_DIR    = "/var/digger/repo"

	//
	EMAIL_CONFIG = "email_config"
)

var (
	LogDir = ""
)
