///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils

import (
	"bytes"
	"digger/models"
	"github.com/go-resty/resty/v2"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	jsoniter "github.com/json-iterator/go"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
)

var (
	proxyStageManagerInstance *proxyStageManager
)

func init() {
	proxyStageManagerInstance = &proxyStageManager{
		lock:       new(sync.Mutex),
		proxiesMap: make(map[int][]*proxyState),
	}
}

func ConvertLogLevel(levelString string) logger.Level {
	levelString = strings.ToLower(levelString)
	switch levelString {
	case "trace":
		return logger.TraceLevel
	case "debug":
		return logger.DebugLevel
	case "info":
		return logger.InfoLevel
	case "warn":
		return logger.WarnLevel
	case "error":
		return logger.ErrorLevel
	case "fatal":
		return logger.FatalLevel
	default:
		return logger.InfoLevel
	}
}

func ParseLabels(labels string) map[string]string {
	sp := strings.Split(labels, ",")
	if len(sp) == 0 {
		return map[string]string{}
	}
	_regex := regexp.MustCompile("([^=]+)=(.*)")
	ret := make(map[string]string)
	for _, l := range sp {
		if !_regex.MatchString(l) {
			continue
		}
		name := _regex.ReplaceAllString(l, "$1")
		value := _regex.ReplaceAllString(l, "$2")
		ret[name] = value
	}
	return ret
}

func ReverseParseLabels(labels string) string {
	labelMap := make(map[string]string)
	jsoniter.UnmarshalFromString(labels, &labelMap)
	var keys []string
	for k := range labelMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	for index, v := range keys {
		buf.WriteString(v)
		buf.WriteString("=")
		buf.WriteString(labelMap[v])
		if index != len(keys)-1 {
			buf.WriteString(",")
		}
	}
	return buf.String()
}

func ParseRedisConnStr(connStr string) *models.RedisConfig {
	_regex := regexp.MustCompile("(.*)@([^@]+)#([0-9]{1,2})")
	if !_regex.MatchString(connStr) {
		return nil
	}
	pass := _regex.ReplaceAllString(connStr, "$1")
	address := _regex.ReplaceAllString(connStr, "$2")
	db := _regex.ReplaceAllString(connStr, "$3")
	dbNo, _ := convert.StrToInt(db)

	return &models.RedisConfig{
		Address:  address,
		Password: pass,
		DB:       dbNo,
	}
}

func ParseEmailNotifierStr(emailNotifier string) *models.EmailConfig {
	_regex := regexp.MustCompile("^([^:]+):(.*)@([^@]+):([0-9]+)$")
	if !_regex.MatchString(emailNotifier) {
		return nil
	}
	// cvfbyhzhqtmvbafj
	user := _regex.ReplaceAllString(emailNotifier, "$1")
	pass := _regex.ReplaceAllString(emailNotifier, "$2")
	host := _regex.ReplaceAllString(emailNotifier, "$3")
	_port := _regex.ReplaceAllString(emailNotifier, "$4")
	port, err := convert.StrToInt(_port)
	if err != nil {
		logger.Error(err)
		port = 465
	}

	return &models.EmailConfig{
		Username: user,
		Password: pass,
		Host:     host,
		Port:     port,
	}
}

// if these param exist in system env , then replace it with system env
func ExchangeEnvValue(key string, then func(envValue string)) {
	envVal := strings.TrimSpace(GetEnv(key))
	if envVal != "" {
		then(envVal)
	}
}

func GetEnv(key string) string {
	return os.Getenv(key)
}

func ParseNodeAffinity(label string) *models.KV {
	_regex := regexp.MustCompile("([^=]+)=(.+)")
	if !_regex.MatchString(label) {
		return nil
	}
	name := _regex.ReplaceAllString(label, "$1")
	value := _regex.ReplaceAllString(label, "$2")
	return &models.KV{name, value}
}

func TryProxy(schema string, client *resty.Client, taskId int, cxt *models.Context) *proxyState {
	// select proxy from project config
	if cxt != nil && len(cxt.Project.Proxies) > 0 {
		if proxy, state := proxyStageManagerInstance.selectProxy(taskId, cxt.Project.Proxies); proxy != nil {
			client.SetProxy(checkProxyUrl(proxy.Address))
			return state
		}
	}
	// if project proxy config is not available, take from environment.
	schema = strings.ToLower(schema)
	if schema == "https" {
		proxyUrl := GetEnv("https_proxy")
		if proxyUrl != "" {
			logger.Info("using proxy: ", proxyUrl)
			client.SetProxy(checkProxyUrl(proxyUrl))
		}
	} else {
		proxyUrl := GetEnv("http_proxy")
		if proxyUrl != "" {
			logger.Info("using proxy: ", proxyUrl)
			client.SetProxy(checkProxyUrl(proxyUrl))
		}
	}
	return nil
}

func checkProxyUrl(address string) string {
	if !strings.HasPrefix(address, "http://") && !strings.HasPrefix(address, "https://") {
		return "http://" + address
	}
	return address
}
