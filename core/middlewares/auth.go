///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package middlewares

import (
	"digger/common"
	"digger/models"
	"digger/services/service"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func AuthorizationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token string
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			tokenStr = c.Query("token")
		}
		if strings.HasPrefix(tokenStr, "Bearer ") {
			tokenStr = tokenStr[7:]
		}

		// 校验token
		user, err := CheckToken(tokenStr)

		// 校验失败，返回错误响应
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &models.RestResponse{
				Resp: models.Resp{
					Code:    1,
					Message: err.Error(),
				},
				Data: nil,
			})
			return
		}

		// 设置用户
		c.Set(common.LOGIN_USER, user)

		// 校验成功
		c.Next()
	}
}

func SecretFunc() jwt.Keyfunc {
	configs, _ := service.ConfigService().ListConfigs()
	if configs["secret"] == "" {
		configs["secret"] = common.DefaultSecret
	}
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(configs["secret"]), nil
	}
}

func CheckToken(tokenStr string) (*models.User, error) {
	token, err := jwt.Parse(tokenStr, SecretFunc())
	if err != nil {
		return nil, err
	}

	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return nil, err
	}

	//验证token，如果token被修改过则为false
	if !token.Valid {
		err = errors.New("token is invalid")
		return nil, err
	}

	id := int(claim["id"].(float64))
	username := claim["username"].(string)

	return &models.User{
		Id:       id,
		Username: username,
	}, nil
}
