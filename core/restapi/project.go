///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"bytes"
	"digger/crontask"
	"digger/models"
	"digger/services/service"
	"digger/utils"
	"github.com/gin-gonic/gin"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/httpx"
	"github.com/hetianyi/gox/logger"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"strings"
	"time"
)

type ProjectQueryResultVO struct {
	Page     int               `json:"page"`
	PageSize int               `json:"pageSize"`
	Total    int64             `json:"total"`
	Data     []*models.Project `json:"data"`
}

// 为项目开始一个新任务
func GetProjectList(c *gin.Context) {
	page := GetIntParameter(c, "page", 1)
	pageSize := GetIntParameter(c, "pageSize", 10)
	order := GetIntParameter(c, "order", 1)

	total, projects, err := service.ProjectService().SelectProjectList(models.ProjectQueryVO{
		PageQueryVO: models.PageQueryVO{
			Page:     page,
			PageSize: pageSize,
		},
		Order: order,
	})

	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	// 查询统计数据
	var ids []int
	for _, p := range projects {
		ids = append(ids, p.Id)
		p.Extras = map[string]interface{}{}
	}
	if len(ids) > 0 {
		ts, _ := service.TaskService().TaskCount(ids...)
		if len(ts) > 0 {
			for _, c := range ts {
				for _, p := range projects {
					if p.Id == c.ProjectId {
						p.Extras["tasks"] = c
						break
					}
				}
			}
		}
	}

	c.JSON(http.StatusOK, Success(&ProjectQueryResultVO{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Data:     projects,
	}))
}

func CreateProject(c *gin.Context) {
	// 绑定请求数据
	var reqData models.Project
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	reqData.Name = strings.TrimSpace(reqData.Name)
	// 参数校验
	if reqData.Name == "" || len(reqData.Name) > 255 {
		c.JSON(http.StatusOK, ErrorMsg("invalid parameter value"))
		return
	}
	if reqData.DisplayName == "" {
		reqData.DisplayName = reqData.Name
	}
	reqData.StartUrl = ""
	reqData.StartStage = ""

	project, err := service.ProjectService().CreateProject(reqData)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(project))
}

func UpdateProject(c *gin.Context) {

	id := c.Param("id")
	pid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	// 绑定请求数据
	var reqData models.Project
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	reqData.Name = strings.TrimSpace(reqData.Name)
	// 参数校验
	if reqData.Name == "" || len(reqData.Name) > 255 {
		c.JSON(http.StatusOK, ErrorMsg("invalid parameter value"))
		return
	}
	if reqData.DisplayName == "" {
		reqData.DisplayName = reqData.Name
	}

	oldProject, err := service.ProjectService().SelectProjectById(pid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	oldProject.Name = reqData.Name
	oldProject.DisplayName = reqData.DisplayName
	oldProject.Remark = reqData.Remark
	oldProject.Tags = reqData.Tags
	oldProject.Cron = reqData.Cron
	oldProject.EnableCron = reqData.EnableCron

	if !oldProject.EnableCron {
		crontask.RemoveCron(pid)
	}

	project, err := service.ProjectService().UpdateProject(*oldProject)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(project))
}

// 为项目开始一个新任务
func StartNewTask(c *gin.Context) {
	//projectId := c.Param("id")
	c.JSON(http.StatusOK, Success(nil))
}

// 为项目开始一个新任务
func DeleteProject(c *gin.Context) {

	id := c.Param("id")
	pid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	success, err := service.ProjectService().DeleteProject(pid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	if !success {
		c.JSON(http.StatusOK, ErrorMsg("未知原因"))
		return
	}

	crontask.RemoveCron(pid)

	c.JSON(http.StatusOK, Success(nil))
}

// 为项目开始一个新任务
func GetProject(c *gin.Context) {

	id := c.Param("id")
	pid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	project, err := service.ProjectService().SelectFullProjectInfo(pid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	if project == nil {
		c.JSON(http.StatusOK, ErrorMsg("not found"))
		return
	}

	plugins, err := service.PluginService().SelectPluginsByProject(pid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	if project.Headers == nil || len(project.Headers) == 0 {
		project.Headers = map[string]string{
			"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML,like Gecko) Chrome/78.0.3904.108 Safari/537.36",
		}
	}
	if project.Settings == nil || len(project.Settings) == 0 {
		project.Settings = map[string]string{
			"CONCURRENT_REQUESTS": "5",
		}
	}
	if project.NodeAffinity == nil || len(project.NodeAffinity) == 0 {
		project.NodeAffinity = []string{""}
	}
	if project.Stages == nil || len(project.Stages) == 0 {
		project.Stages = []models.Stage{
			{
				Fields: []models.Field{{}},
			},
		}
	}

	yml, err := utils.Convert2Yaml(project)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(map[string]interface{}{
		"project": project,
		"yaml":    yml,
		"plugins": plugins,
	}))
}

func SaveProjectConfig(c *gin.Context) {

	id := c.Param("id")
	pid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	// 绑定请求数据
	var reqData models.PlayInputVO2
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	// 解析配置文件
	var project models.Project
	err = utils.ParseYamlFromString(reqData.Project, &project)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	// TODO validate config

	oldProject, err := service.ProjectService().SelectProjectById(pid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	oldProject.Headers = project.Headers
	oldProject.Settings = project.Settings
	oldProject.StartUrl = project.StartUrl
	oldProject.StartStage = project.StartStage
	oldProject.NodeAffinity = project.NodeAffinity

	// 默认设置
	if oldProject.Settings == nil {
		oldProject.Settings = make(map[string]string)
	}
	if oldProject.Headers == nil {
		oldProject.Headers = make(map[string]string)
	}
	if len(oldProject.Settings) == 0 || oldProject.Settings["CONCURRENT_REQUESTS"] == "" {
		oldProject.Settings["CONCURRENT_REQUESTS"] = "5"
	}
	if len(oldProject.Headers) == 0 || oldProject.Headers["User-Agent"] == "" {
		oldProject.Headers["User-Agent"] = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML,likeGecko) Chrome/78.0.3904.108 Safari/537.36"
	}

	err = service.ProjectConfigService().SaveProjectConfig(oldProject, project.Stages)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(nil))
}

func ExportProjectConfig(c *gin.Context) {
	id := c.Param("id")
	pid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	project, err := service.ProjectService().SelectFullProjectInfo(pid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	marshalIndentBytes, err := jsoniter.MarshalIndent(project, "", "  ")
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	reader := bytes.NewReader(marshalIndentBytes)
	c.Writer.Header().Set("Content-Disposition", "attachment;filename=\""+project.Name+".cfg.json\"")
	httpx.ServeContent(c.Writer, c.Request, "", time.Now(), reader, int64(len(marshalIndentBytes)))
}

func ImportProjectConfig(c *gin.Context) {

	id := c.Param("id")
	pid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	fileHeaders := form.File["config"]
	if len(fileHeaders) == 0 {
		c.JSON(http.StatusOK, ErrorMsg("no content"))
		return
	}
	file := fileHeaders[0]
	if file.Size > 1024*1024 {
		c.JSON(http.StatusOK, ErrorMsg("file exceeds max size"))
		return
	}
	reader, err := file.Open()
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	buf := make([]byte, file.Size)
	_, err = reader.Read(buf)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	var reqData models.Project
	err = jsoniter.Unmarshal(buf, &reqData)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	project, err := service.ProjectService().SelectFullProjectInfo(pid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	if project == nil {
		c.JSON(http.StatusOK, ErrorMsg("project not exists"))
		return
	}

	logger.Info("开始导入配置")
	reqData.Id = pid
	if err = service.ProjectConfigService().ImportProjectConfig(&reqData); err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	logger.Info("配置导入成功")
	c.JSON(http.StatusOK, Success(nil))
}
