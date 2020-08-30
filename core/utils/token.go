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
	type MyCustomClaims struct {
		Id int `json:"id"`
		Username string `json:"username"`
		jwt.StandardClaims
	}

	// Create the Claims
	claims := MyCustomClaims{
		user.Id,
		user.Username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour*12).Unix(),
			Issuer:    user.Username,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(common.DefaultSecret))
}
