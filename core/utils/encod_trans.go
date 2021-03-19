///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils

import (
	"github.com/axgle/mahonia"
	"github.com/hetianyi/gox/logger"
	"strings"
)

var (
	targetEncoding      = mahonia.NewDecoder("UTF-8")
	srcEncodingAliasMap = map[string]string{
		"GB2312": "GBK",
		"GBK":    "GBK",
	}
	srcEncodingMap = make(map[string]mahonia.Decoder)
)

func getSourceEncoding(charset string) mahonia.Decoder {
	charset = strings.ToUpper(charset)
	srcEncoding := srcEncodingMap[charset]
	if srcEncoding == nil {
		srcEncoding = srcEncodingMap[srcEncodingAliasMap[charset]]
		if srcEncoding == nil {
			srcEncoding = mahonia.NewDecoder(srcEncodingAliasMap[charset])
		}
		srcEncodingMap[charset] = srcEncoding
	}
	return srcEncoding
}

// Trans2UTF8 将源编码转换为UTF-8编码，
// 如果不支持该编码，将返回原文
func Trans2UTF8(fromCharset, srcContent string) string {
	fromCharset = strings.ToLower(fromCharset)
	if fromCharset == "" || fromCharset == "utf-8" {
		return srcContent
	}
	srcEncoding := getSourceEncoding(fromCharset)
	if srcEncoding == nil {
		logger.Warn("charset \"", fromCharset, "\" is not supported")
		return srcContent
	}
	srcResult := srcEncoding.ConvertString(srcContent)
	_, cdata, err := targetEncoding.Translate([]byte(srcResult), true)
	if err != nil {
		logger.Error("cannot transform encoding from ", fromCharset, " to UTF-8: ", err)
		return srcContent
	}
	return string(cdata)
}
