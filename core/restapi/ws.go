///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package restapi

import (
	"digger/dispatcher"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/hetianyi/gox/logger"
	"net/http"
)

func Ws(c *gin.Context) {

	host := c.Request.Header.Get("X-Real-Ip")
	if host == "" {
		host = c.Request.Header.Get("Host")
		if host == "" {
			host = c.Request.RemoteAddr
		}
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

	if err := dispatcher.RegisterWsConnection(conn, host); err != nil {
		logger.Error(err)
		conn.Close()
	}
}
