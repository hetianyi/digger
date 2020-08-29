///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"bytes"
	"digger/models"
	"digger/services/service"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/httpx"
	"github.com/hetianyi/gox/logger"
	jsoniter "github.com/json-iterator/go"
	"github.com/mholt/archiver"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type ResultQueryResultVO struct {
	Page     int              `json:"page"`
	PageSize int              `json:"pageSize"`
	Total    int64            `json:"total"`
	Data     []*models.Result `json:"data"`
}

func QueryResult(c *gin.Context) {

	page := GetIntParameter(c, "page", 1)
	pageSize := GetIntParameter(c, "pageSize", 20)
	taskId := GetIntParameter(c, "taskId", 0)

	var reqBody = models.ResultQueryVO{
		PageQueryVO: models.PageQueryVO{
			PageSize: pageSize,
			Page:     page,
		},
		TaskId: taskId,
	}

	if reqBody.TaskId == 0 {
		c.JSON(http.StatusOK, Success(&ResultQueryResultVO{
			Page:     reqBody.Page,
			PageSize: reqBody.PageSize,
			Total:    0,
		}))
		return
	}

	total, arr, err := service.ResultService().SelectResults(reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(&ResultQueryResultVO{
		Page:     reqBody.Page,
		PageSize: reqBody.PageSize,
		Total:    total,
		Data:     arr,
	}))
}

// 为项目开始一个新任务
func SaveResult(c *gin.Context) {

	// 绑定数据
	/*var reqBody models.QueueCallbackRequestVO
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
		return
	}

	if reqBody.Check == "" {
		c.JSON(http.StatusOK, Success(nil))
		return
	}

	// 已经超过最大重试次数的queue
	var exceedsIds []int64
	counts, err := service.CacheService().IncreQueueErrorCount(reqBody.ErrorQueueTaskIds, reqBody.ErrorQueueIds)
	for i, v := range counts {
		if v > common.MAX_RETRY {
			exceedsIds = append(exceedsIds, reqBody.ErrorQueueIds[i])
		}
	}*/

	/*m := groupByTaskId(&reqBody)

	var finalSuccessQueueIds []int64
	for taskId, ids := range m {
		// 过滤已完成的任务
		doneTable, err := service.CacheService().ExistMembers(taskId, ids)
		if err == nil {
			for i, v := range doneTable {
				if !v {
					finalSuccessQueueIds = append(finalSuccessQueueIds, ids[i].(int64))
				}
			}
		}
	}
	reqBody.SuccessQueueIds = finalSuccessQueueIds*/
	// TODO 过期不存在的check则丢弃此次结果
	/*if !scheduler.ExistCheck(reqBody.Check) {
		logger.Error("check不存在，丢弃结果")
		c.JSON(http.StatusOK, Success(nil))
		return
	}

	scheduler.LockChecks()
	defer scheduler.UnLockChecks()

	err = service.ResultService().SaveCheckData(&reqBody, exceedsIds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorMsg(err.Error()))
	} else {
		scheduler.RemoveCheck(&reqBody)
		service.CacheService().SaveSuccessQueueIds(&reqBody)
		c.JSON(http.StatusOK, Success(nil))
	}*/
}

func ExportResult(c *gin.Context) {

	format := strings.ToLower(GetStrParameter(c, "format", "sql"))
	taskId := GetIntParameter(c, "taskId", 0)

	task, err := service.TaskService().SelectTask(taskId)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	if task == nil {
		c.JSON(http.StatusOK, ErrorMsg("task not exists"))
		return
	}

	project, err := service.ProjectService().SelectProjectById(task.ProjectId)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	if project == nil {
		c.JSON(http.StatusOK, ErrorMsg("project not exists"))
		return
	}

	if format != "sql" && format != "json" && format != "csv" {
		c.JSON(http.StatusOK, ErrorMsg("not supported format"))
		return
	}

	tempFileName := fmt.Sprintf("%s%s-%d-%d.%s", gox.TValue(format == "sql", "t_", "").(string), strings.ToLower(project.Name), taskId, gox.GetTimestamp(time.Now()), format)

	resultFile, err := file.CreateFile(os.TempDir() + "/" + tempFileName)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	defer func() {
		resultFile.Close()
		os.Remove(resultFile.Name())
	}()

	page := 0
	pageSize := 10000
	writeCol := true
	var buf bytes.Buffer
	var csvWriter *csv.Writer
	if format == "csv" {
		csvWriter = csv.NewWriter(resultFile)
		// 写入UTF-8 BOM，避免乱码
		resultFile.WriteString("\xEF\xBB\xBF")
	}

	for {
		page++
		buf.Reset()
		logger.Info(fmt.Sprintf("正在导出第%d页...", page))
		_, trs, err := service.ResultService().SelectResults(models.ResultQueryVO{
			PageQueryVO: models.PageQueryVO{
				Page:     page,
				PageSize: pageSize,
			},
			TaskId: taskId,
		})
		if err != nil {
			c.JSON(http.StatusOK, ErrorMsg(err.Error()))
			return
		}
		for _, r := range trs {
			if err = buildResult(format, &buf, csvWriter, writeCol, project, r); err != nil {
				c.JSON(http.StatusOK, ErrorMsg(err.Error()))
				return
			}
			writeCol = false
		}
		if _, err := resultFile.WriteString(buf.String()); err != nil {
			c.JSON(http.StatusOK, ErrorMsg(err.Error()))
			return
		}
		if len(trs) < pageSize {
			break
		}
	}
	if csvWriter != nil {
		csvWriter.Flush()
	}

	resultFile.Close()

	compressFile := tempFileName + ".tar.gz"
	// 压缩文件
	err = archiver.Archive([]string{resultFile.Name()}, compressFile)
	if err != nil {
		logger.Error("err compressing file:", err)
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	compResultFile, _ := file.GetFile(compressFile)
	info, _ := os.Stat(compressFile)

	defer func() {
		compResultFile.Close()
		os.Remove(compressFile)
	}()

	downloadName := fmt.Sprintf("attachment;filename=\"%s\"", filepath.Base(compressFile))

	c.Writer.Header().Set("Content-Disposition", downloadName)
	httpx.ServeContent(c.Writer, c.Request, "", time.Now(), compResultFile, info.Size())
}

func buildResult(format string,
	buff *bytes.Buffer,
	csvWriter *csv.Writer,
	writeCol bool,
	project *models.Project,
	r *models.Result) error {
	switch format {
	case "sql":
		return buildSQLItem(buff, project, r)
	case "json":
		return buildJSONItem(buff, project, r)
	case "csv":
		return buildCSVItem(csvWriter, writeCol, r)
	default:
		return nil
	}
}

func buildSQLItem(buff *bytes.Buffer, project *models.Project, r *models.Result) error {
	m := make(map[string]string)
	err := jsoniter.UnmarshalFromString(r.Result, &m)
	if err != nil {
		return err
	}
	var fs []string
	for k := range m {
		fs = append(fs, k)
	}
	sort.Strings(fs)
	fLen := len(fs)
	if fLen == 0 {
		return nil
	}
	buff.WriteString("insert into t_")
	buff.WriteString(strings.ToLower(project.Name))
	buff.WriteString("(")
	for i, v := range fs {
		buff.WriteString(v)
		if i != fLen-1 {
			buff.WriteString(",")
		}
	}
	buff.WriteString(") values (")
	for i, v := range fs {
		buff.WriteString("'")
		buff.WriteString(strings.ReplaceAll(m[v], "'", "''"))
		buff.WriteString("'")
		if i != fLen-1 {
			buff.WriteString(",")
		}
	}
	buff.WriteString(");\n")
	return nil
}

func buildCSVItem(csvWriter *csv.Writer, writeCol bool, r *models.Result) error {
	m := make(map[string]string)
	err := jsoniter.UnmarshalFromString(r.Result, &m)
	if err != nil {
		return err
	}
	var fs []string
	for k := range m {
		fs = append(fs, k)
	}
	sort.Strings(fs)
	fLen := len(fs)
	if fLen == 0 {
		return nil
	}

	if writeCol {
		if err := csvWriter.Write(fs); err != nil {
			return err
		}
	}

	var record = make([]string, fLen)
	for i, v := range fs {
		record[i] = m[v]
	}
	if err := csvWriter.Write(record); err != nil {
		return err
	}
	return nil
}

func buildJSONItem(buff *bytes.Buffer, project *models.Project, r *models.Result) error {
	buff.WriteString(r.Result)
	buff.WriteString("\n")
	return nil
}

func groupByTaskId(reqBody *models.QueueCallbackRequestVO) map[int][]interface{} {
	ret := make(map[int][]interface{})
	for i, v := range reqBody.SuccessQueueIds {
		t := ret[reqBody.SuccessQueueTaskIds[i]]
		t = append(t, v)
		ret[reqBody.SuccessQueueTaskIds[i]] = t
	}
	return ret
}
