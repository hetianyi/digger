///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"digger/common"
	"digger/dispatcher"
	"digger/models"
	"digger/scheduler"
	"digger/services/service"
	"digger/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
	"github.com/hpcloud/tail"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os"
	"time"
)

type TaskQueryResultVO struct {
	Page     int            `json:"page"`
	PageSize int            `json:"pageSize"`
	Total    int64          `json:"total"`
	Data     []*models.Task `json:"data"`
}

// 为项目开始一个新任务
func GetProjectTasks(c *gin.Context) {

	page := GetIntParameter(c, "page", 1)
	pageSize := GetIntParameter(c, "pageSize", 20)
	projectId := GetIntParameter(c, "projectId", 0)
	status := GetIntParameter(c, "status", -1)

	total, tasks, err := service.TaskService().SelectTaskList(models.TaskQueryVO{
		PageQueryVO: models.PageQueryVO{
			Page:     page,
			PageSize: pageSize,
		},
		ProjectId: projectId,
		Status:    status,
	})
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	// 查询统计数据
	if len(tasks) > 0 {
		var taskIds []int
		for _, t := range tasks {
			taskIds = append(taskIds, t.Id)
		}
		sc, err := service.ResultService().ResultCount(taskIds...)
		if err != nil {
			c.JSON(http.StatusOK, ErrorMsg(err.Error()))
			return
		}
		for _, c := range sc {
			for _, t := range tasks {
				if c.TaskId == t.Id {
					t.ResultCount = c.Count
					break
				}
			}
		}

		ec, err := service.QueueService().ErrorCount(taskIds...)
		if err != nil {
			c.JSON(http.StatusOK, ErrorMsg(err.Error()))
			return
		}
		for _, c := range ec {
			for _, t := range tasks {
				if c.TaskId == t.Id {
					t.ErrorRequest = c.Count
					break
				}
			}
		}
	}

	c.JSON(http.StatusOK, Success(&TaskQueryResultVO{
		Page:     page,
		PageSize: pageSize,
		Total:    total,
		Data:     tasks,
	}))
}

// 为项目开始一个新任务
func GetTasks(c *gin.Context) {
	tasks, err := service.TaskService().SelectActiveTasks()
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	c.JSON(http.StatusOK, Success(tasks))
}

// 为项目开始一个新任务
func NewTasks(c *gin.Context) {
	logger.Info("创建新任务")
	_projectId := c.Query("projectId")
	projectId, err := convert.StrToInt(_projectId)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg("Invalid Parameter"))
		return
	}

	project, err := service.ProjectService().SelectFullProjectInfo(projectId)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	if project == nil {
		c.JSON(http.StatusOK, ErrorMsg("project not found"))
		return
	}

	success, err := project.Validate()
	if !success {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	task, err := service.TaskService().CreateTask(models.Task{
		ProjectId:  projectId,
		Status:     1,
		CreateTime: time.Now(),
	})

	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	err = scheduler.Schedule(task)

	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(task))
}

//
func StopTasks(c *gin.Context) {
	id := c.Param("id")
	tid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	err = service.TaskService().ShutdownTask(tid)
	scheduler.Stop(tid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	// 清理redis缓存
	cleanCache(tid)

	c.JSON(http.StatusOK, Success(nil))
}

//
func PauseTasks(c *gin.Context) {
	id := c.Param("id")
	tid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	err = service.TaskService().PauseTask(tid)
	dispatcher.Pause(tid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(nil))
}

func ContinueTasks(c *gin.Context) {
	id := c.Param("id")
	tid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	err = service.TaskService().StartTask(tid)
	dispatcher.Continue(tid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, Success(nil))
}

func TaskLogs(c *gin.Context) {
	id := c.Param("id")
	tid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	task, err := service.TaskService().SelectTask(tid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	if task == nil {
		c.JSON(http.StatusOK, ErrorMsg("task not exists"))
		return
	}
	if task.Status == 2 || task.Status == 3 {
		c.JSON(http.StatusOK, ErrorMsg("task was stopped or finished"))
		return
	}

	fileName := common.LogDir + "/" + convert.IntToStr(tid) + ".log"
	if !file.Exists(fileName) {
		f, err := file.CreateFile(fileName)
		if err != nil {
			c.JSON(http.StatusOK, ErrorMsg(err.Error()))
			return
		}
		f.Close()
	}

	info, err := os.Stat(fileName)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	var offset int64 = 0
	if info.Size() > 5*1024 {
		offset = info.Size() - int64(5*1024)
	}

	ta, err := tail.TailFile(fileName, tail.Config{
		Location: &tail.SeekInfo{
			Offset: offset,
			Whence: 0,
		},
		Follow: true,
		Logger: logrus.StandardLogger(),
	})
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	} // use default options
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	defer func() {
		logger.Info("清理...")
		conn.Close()
		ta.Stop()
		ta.Cleanup()
	}()

	logger.Info("开始读取日志")

	for line := range ta.Lines {
		if line == nil {
			continue
		}
		//decodeurl := url.QueryEscape(line.Text)
		if err := conn.WriteMessage(1, []byte(utils.EncodeBase64(url.QueryEscape(line.Text)))); err != nil {
			conn.Close()
			break
		}
	}
	logger.Info("查看日志结束")
}

//
func DeleteTask(c *gin.Context) {
	id := c.Param("id")
	tid, err := convert.StrToInt(id)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	task, err := service.TaskService().SelectTask(tid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	if task == nil {
		c.JSON(http.StatusOK, ErrorMsg("task not exists"))
		return
	}
	if task.Status != 2 && task.Status != 3 {
		c.JSON(http.StatusOK, ErrorMsg("task cannot be deleted this time"))
		return
	}

	err = service.TaskService().DeleteTask(tid)
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}

	// 清理redis缓存
	cleanCache(tid)

	c.JSON(http.StatusOK, Success(nil))
}

func cleanCache(taskId int) {
	logger.Info("清理redis缓存")
	service.RedisClient.Del(fmt.Sprintf("DONE_QUEUE:%d", taskId))
	service.RedisClient.Del(fmt.Sprintf("FINISH_RES:%d", taskId))
	service.RedisClient.Del(fmt.Sprintf("ERR_QUEUE:%d", taskId))
}
