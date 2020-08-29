package utils

import (
	"crypto/tls"
	"digger/models"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/logger"
	"gopkg.in/gomail.v2"
	"time"
)

func EmailNotify(task *models.Task, config *models.EmailConfig) error {
	logger.Info("发送邮件通知...")
	now := time.Now()
	d := gomail.NewDialer(config.Host, config.Port, config.Username, config.Password)
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Username)
	m.SetHeader("To", config.Username)
	m.SetHeader("Subject", "【Digger爬虫任务通知】ID "+convert.IntToStr(task.Id))
	m.SetBody("text/html", `
<style type="text/css">
table {
	-webkit-border-horizontal-spacing: 0px;
	-webkit-border-vertical-spacing: 0px;
	border-top-width: 0px;
	border-right-width: 0px;
	border-bottom-width: 0px;
	border-left-width: 0px;
	display: table;
    border-collapse: separate;
    box-sizing: border-box;
    border-spacing: 2px;
    border-color: grey;
}

th {
	position: relative;
    height: 100%;
    padding: 8px 0;
	border-right: 1px solid #e8eaec;
    height: 40px;
    white-space: nowrap;
    overflow: hidden;
    background-color: #f8f8f9;
    min-width: 0;
    height: 48px;
    box-sizing: border-box;
    text-align: left;
    text-overflow: ellipsis;
    vertical-align: middle;
    border-bottom: 1px solid #e8eaec;
	text-align: center;
}
td {
    border-right: 1px solid #e8eaec;
    background-color: #fff;
    transition: background-color .2s ease-in-out;
    min-width: 0;
    height: 48px;
    box-sizing: border-box;
    text-align: left;
    text-overflow: ellipsis;
    vertical-align: middle;
    border-bottom: 1px solid #e8eaec;
	text-align: center;
}

</style>
<h2>任务已结束</h2>
<table>
	<tr>
	<th style="width: 80px;">ID</th>
	<th style="width: 200px;">开始时间</th>
	<th style="width: 200px;">结束时间</th>
	</tr>
	<tr>
	<td>`+convert.IntToStr(task.Id)+`</td>
	<td>`+gox.GetLongDateString(task.CreateTime)+`</td>
	<td>`+gox.GetLongDateString(now)+`</td>
	</tr>
</table>
`)

	var e error
	for i := 0; i < 3; i++ {
		if err := d.DialAndSend(m); err != nil {
			logger.Error("无法发送邮件通知: ", err)
			e = err
			continue
		}
		e = nil
		break
	}
	if e == nil {
		logger.Info("邮件通知发送成功")
	}
	return e
}
