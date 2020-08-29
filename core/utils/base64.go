///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils

import "encoding/base64"

func EncodeBase64(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

func DecodeBase64(input string) (string, error) {
	bs, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
