///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"bytes"
	"digger/common"
	"digger/models"
	"digger/services/service"
	"encoding/csv"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/httpx"
	"github.com/hetianyi/gox/logger"
	json "github.com/json-iterator/go"
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
		TaskId:       taskId,
		LastResultId: 0,
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

type exportContext struct {
	project          *models.Project
	result           *models.Result
	csvWriter        *csv.Writer
	fieldTypeMapping map[string]bool
	format           string
	writeColumn      bool
	isFirstRecord    bool
	isLastRecord     bool
	buffer           *bytes.Buffer
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

	project, err := service.CacheService().GetSnapshotConfig(task.Id)
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

	fieldTypeMapping := getFieldTypeMapping(project)

	exportPageSize := project.GetIntSetting(common.SETTINGS_EXPORT_PAGE_SIZE, 1000)

	tempFileName := fmt.Sprintf("%s%s-%d-%d.%s",
		gox.TValue(format == "sql", "t_", "").(string),
		strings.ToLower(project.Name),
		taskId,
		gox.GetTimestamp(time.Now()),
		format)

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
	writeCol := true
	var buf bytes.Buffer
	var csvWriter *csv.Writer
	if format == "csv" {
		csvWriter = csv.NewWriter(resultFile)
		// 写入UTF-8 BOM，避免乱码
		resultFile.WriteString("\xEF\xBB\xBF")
	}

	cxt := &exportContext{
		project:          project,
		csvWriter:        csvWriter,
		fieldTypeMapping: fieldTypeMapping,
		isFirstRecord:    true,
		isLastRecord:     false,
		writeColumn:      writeCol,
		format:           format,
		buffer:           &buf,
	}

	var lastResultId int64 = 0
	for {
		page++
		buf.Reset()
		logger.Info(fmt.Sprintf("正在导出第%d页", page))
		trs, err := service.ResultService().ExportResults(models.ResultQueryVO{
			PageQueryVO: models.PageQueryVO{
				Page:     page,
				PageSize: exportPageSize,
			},
			TaskId:       taskId,
			LastResultId: lastResultId,
		})
		if err != nil {
			c.JSON(http.StatusOK, ErrorMsg(err.Error()))
			return
		}
		if len(trs) < exportPageSize {
			cxt.isLastRecord = true
		}
		for _, r := range trs {
			cxt.result = r
			if err = buildResult(cxt); err != nil {
				c.JSON(http.StatusOK, ErrorMsg(err.Error()))
				return
			}
			writeCol = false
		}
		// 到达最后一页
		if cxt.isLastRecord && format == "json" {
			cxt.buffer.WriteString("\n]")
		}
		lastResultId = trs[len(trs)-1].Id
		if _, err := resultFile.WriteString(buf.String()); err != nil {
			c.JSON(http.StatusOK, ErrorMsg(err.Error()))
			return
		}
		if len(trs) < exportPageSize {
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
		logger.Error("无法压缩导出文件: ", err)
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

func buildResult(cxt *exportContext) error {
	switch cxt.format {
	case "sql":
		return buildSQLItem(cxt)
	case "json":
		return buildJSONItem(cxt)
	case "csv":
		return buildCSVItem(cxt)
	default:
		return nil
	}
}

func buildSQLItem(cxt *exportContext) error {
	valueMap := make(map[string]interface{})
	err := json.UnmarshalFromString(cxt.result.Result, &valueMap)
	if err != nil {
		return err
	}
	var fs []string
	for k := range cxt.fieldTypeMapping {
		fs = append(fs, k)
	}
	sort.Strings(fs)
	fLen := len(fs)
	if fLen == 0 {
		return nil
	}
	cxt.buffer.WriteString("insert into t_")
	cxt.buffer.WriteString(strings.ToLower(cxt.project.Name))
	cxt.buffer.WriteString("(")
	for i, v := range fs {
		cxt.buffer.WriteString("`")
		cxt.buffer.WriteString(v)
		cxt.buffer.WriteString("`")
		if i != fLen-1 {
			cxt.buffer.WriteString(",")
		}
	}
	cxt.buffer.WriteString(") values (")
	for i, v := range fs {
		fv := ""
		if cxt.fieldTypeMapping[v] {
			fv, _ = json.MarshalToString(valueMap[v])
		} else {
			fv = valueMap[v].(string)
		}
		cxt.buffer.WriteString("'")
		cxt.buffer.WriteString(strings.ReplaceAll(strings.ReplaceAll(fv, "'", "''"), "\\", "\\\\"))
		cxt.buffer.WriteString("'")
		if i != fLen-1 {
			cxt.buffer.WriteString(",")
		}
	}
	cxt.buffer.WriteString(");\n")
	return nil
}

// 获取字段类型
// isArray: true
// other: false
func getFieldTypeMapping(project *models.Project) map[string]bool {
	mapping := make(map[string]bool)
	for _, s := range project.Stages {
		for _, f := range s.Fields {
			if f.IsArray {
				mapping[f.Name] = true
			} else {
				mapping[f.Name] = false
			}
		}
	}
	return mapping
}

func buildCSVItem(cxt *exportContext) error {

	var fs []string
	for k := range cxt.fieldTypeMapping {
		fs = append(fs, k)
	}
	sort.Strings(fs)
	fLen := len(fs)
	if fLen == 0 {
		return nil
	}

	if cxt.writeColumn {
		cxt.writeColumn = false
		if err := cxt.csvWriter.Write(fs); err != nil {
			return err
		}
	}

	valueMap := make(map[string]interface{})
	err := json.UnmarshalFromString(cxt.result.Result, &valueMap)
	if err != nil {
		return err
	}

	var record = make([]string, fLen)
	for i, v := range fs {
		if !cxt.fieldTypeMapping[v] {
			record[i] = valueMap[v].(string)
		} else {
			line, _ := json.MarshalToString(valueMap[v])
			record[i] = line
		}
	}
	if err := cxt.csvWriter.Write(record); err != nil {
		return err
	}
	return nil
}

func buildJSONItem(cxt *exportContext) error {
	if cxt.isFirstRecord {
		cxt.buffer.WriteString("[\n")
		cxt.isFirstRecord = false
	} else {
		cxt.buffer.WriteString(",\n")
	}
	cxt.buffer.WriteString(cxt.result.Result)
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
