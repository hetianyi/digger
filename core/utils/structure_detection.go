///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func DetectStructure(content string) {

}

// 用goquery解析html文档
func parseCssDocument(content []byte) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(bytes.NewBuffer(content))
}

// 用goquery解析html文档
func parseXpathDocument(content []byte) (*html.Node, error) {
	return htmlquery.Parse(bytes.NewBuffer(content))
}
