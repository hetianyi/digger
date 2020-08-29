///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"digger/dispatcher"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetNodes(c *gin.Context) {
	c.JSON(http.StatusOK, Success(dispatcher.GetNodes()))
}
