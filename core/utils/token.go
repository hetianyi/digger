///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils

import (
	"digger/common"
	"digger/models"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func MakeToken(user *models.User) (tokenStr string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.Id,
		"username": user.Username,
		"nbf":      time.Now().Unix(),
	})
	return token.SignedString([]byte(common.DefaultSecret))
}
