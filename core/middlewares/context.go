///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package middlewares

import (
	"digger/common"
	"digger/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Context struct {
	*gin.Context
}

func (c *Context) User() *models.User {
	userIfe, exists := c.Get(common.LOGIN_USER)
	if !exists {
		return nil
	}
	user, ok := userIfe.(*models.User)
	if !ok {
		return nil
	}
	return user
}
func (c *Context) Success(data interface{}, metas ...interface{}) {
	var meta interface{}
	if len(metas) == 0 {
		meta = gin.H{}
	} else {
		meta = metas[0]
	}
	if data == nil {
		data = gin.H{}
	}
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "success",
		"data":    data,
		"meta":    meta,
		"error":   "",
	})
}

func WithGinContext(context *gin.Context) *Context {
	return &Context{Context: context}
}
