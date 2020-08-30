///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"digger/common"
	context "digger/middlewares"
	"digger/models"
	"digger/services/service"
	"digger/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserRequestData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c *gin.Context) {
	// 绑定请求数据
	var reqData UserRequestData
	if err := c.ShouldBindJSON(&reqData); err != nil {
		c.JSON(http.StatusOK, ErrorMsg("not authorized"))
		return
	}

	configs, err := service.ConfigService().ListConfigs()
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg(err.Error()))
		return
	}
	if configs["admin_user"] == "" {
		configs["admin_user"] = DefaultUser.Username
	}
	if configs["admin_password"] == "" {
		configs["admin_password"] = DefaultUser.Password
	}
	if configs["secret"] == "" {
		configs["secret"] = common.DefaultSecret
	}

	u := &models.User{
		Id:       DefaultUser.Id,
		Username: configs["admin_user"],
		Password: configs["admin_password"],
	}

	// 校验密码
	//encPassword := utils.EncryptPassword(reqData.Password)
	if u.Username != reqData.Username || u.Password != reqData.Password {
		c.JSON(http.StatusOK, ErrorMsg("not authorized"))
		return
	}

	// 获取token
	tokenStr, err := utils.MakeToken(u, configs["secret"])
	if err != nil {
		c.JSON(http.StatusOK, ErrorMsg("not authorized"))
		return
	}

	c.JSON(http.StatusOK, Success(tokenStr))
}

func GetUserInfo(c *gin.Context) {
	ctx := context.WithGinContext(c)
	user := ctx.User()
	if user == nil {
		c.JSON(http.StatusOK, ErrorMsg("not authorized"))
		return
	}
	c.JSON(http.StatusOK, Success(user))
}
