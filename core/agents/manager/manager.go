///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package manager

import (
	"context"
	"digger/common"
	"digger/crontask"
	"digger/dispatcher"
	"digger/middlewares"
	"digger/models"
	"digger/restapi"
	"digger/scheduler"
	"digger/services/service"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/file"
	"github.com/hetianyi/gox/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	config *models.BootstrapConfig
)

// 启动manager
func StartAgentManager(_config *models.BootstrapConfig) {

	config = _config

	if !file.Exists(common.REPO_DIR) {
		if err := file.CreateDirs(common.REPO_DIR); err != nil {
			logger.Fatal(err)
		}
	}

	// 初始化数据库
	service.InitDb(config.DBString)

	// 初始化redis
	service.InitRedis(config.RedisString)

	// 启动之前已激活的任务
	scheduler.StartScheduler()

	dispatcher.StartDispatcher(config)

	crontask.StartCronScanner()

	// 订阅事件消息
	// go event.StartEventManager()

	gin.SetMode(gin.ReleaseMode)

	app := gin.Default()
	app.Use(middlewares.CORSMiddleware())

	anonymousGroup := app.Group("/")
	{
		anonymousGroup.POST("/api/v1/login", restapi.Login)
		anonymousGroup.POST("/api/v1/logout", restapi.Login)
		anonymousGroup.GET("/api/ws", restapi.Ws)
		anonymousGroup.GET("/api/tasks/:id/logs/ws", restapi.TaskLogs)
		anonymousGroup.GET("/api/v1/snapshot", restapi.GetProjectConfigSnapshot)
	}

	authGroup := app.Group("/api", middlewares.AuthorizationMiddleware())
	{
		authGroup.GET("/v1/user", restapi.GetUserInfo)
		authGroup.GET("/v1/projects", restapi.GetProjectList)
		authGroup.POST("/v1/projects", restapi.CreateProject)
		authGroup.POST("/v1/projects/:id", restapi.StartNewTask)
		authGroup.PUT("/v1/projects/:id", restapi.UpdateProject)
		authGroup.DELETE("/v1/projects/:id", restapi.DeleteProject)
		authGroup.GET("/v1/projects/:id", restapi.GetProject)

		authGroup.PUT("/v1/projects/:id/config", restapi.SaveProjectConfig)

		authGroup.GET("/v1/queues", restapi.GetQueues)

		authGroup.PUT("/v1/projects/:id/plugins", restapi.SavePlugins)
		authGroup.GET("/v1/projects/:id/export", restapi.ExportProjectConfig)
		authGroup.POST("/v1/projects/:id/import", restapi.ImportProjectConfig)

		authGroup.GET("/v1/tasks", restapi.GetProjectTasks)
		authGroup.POST("/v1/tasks", restapi.NewTasks)
		authGroup.PUT("/v1/tasks/:id/stop", restapi.StopTasks)
		authGroup.PUT("/v1/tasks/:id/pause", restapi.PauseTasks)
		authGroup.PUT("/v1/tasks/:id/continue", restapi.ContinueTasks)
		authGroup.DELETE("/v1/tasks/:id", restapi.DeleteTask)

		authGroup.GET("/v1/results", restapi.QueryResult)
		authGroup.GET("/v1/results/export", restapi.ExportResult)

		authGroup.GET("/v1/configs", restapi.GetConfigs)
		authGroup.PUT("/v1/configs", restapi.UpdateConfig)

		authGroup.GET("/v1/nodes", restapi.GetNodes)

		authGroup.GET("/v1/statistics", restapi.GetStatistic)

		authGroup.POST("/v1/play", restapi.PlayExistStage)
		authGroup.POST("/v2/play", restapi.PlayFromTempStage)
		authGroup.POST("/v2/play/parse", restapi.ParseConfigFile)

		authGroup.GET("/v1/proxies", restapi.QueryProxy)
		authGroup.POST("/v1/proxies", restapi.SaveProxy)
		authGroup.DELETE("/v1/proxies", restapi.DeleteProxy)

		authGroup.GET("/v1/pushes", restapi.QueryPushSource)
		authGroup.POST("/v1/pushes", restapi.SavePushSource)
		authGroup.DELETE("/v1/pushes", restapi.DeletePush)
	}

	app.Use(static.Serve("/", static.LocalFile(config.UIDir, true)))

	// TODO
	srv := &http.Server{
		Handler: app,
		Addr:    "0.0.0.0:" + convert.IntToStr(config.Port),
	}

	logger.Info("server is listening on 0.0.0.0:", config.Port)
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				logger.Error("run server error:" + err.Error())
			} else {
				logger.Info("server graceful down")
			}
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx2, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx2); err != nil {
		logger.Error("run server error:" + err.Error())
	}
}
