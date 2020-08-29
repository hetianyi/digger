///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	context "digger/middlewares"
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

	// 校验密码
	//encPassword := utils.EncryptPassword(reqData.Password)
	if DefaultUser.Username != reqData.Username || DefaultUser.Password != reqData.Password {
		c.JSON(http.StatusOK, ErrorMsg("not authorized"))
		return
	}

	user := DefaultUser

	// 获取token
	tokenStr, err := utils.MakeToken(user)
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
