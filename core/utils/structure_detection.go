///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/hetianyi/gox/logger"
	"golang.org/x/net/html"
)

func DetectStructure(content []byte) {
	doc, err := parseXpathDocument(content)
	if err != nil {
		logger.Fatal(err)
	}
	iteratorChildren(doc)
}

func iteratorNext(node *html.Node, brothers []*html.Node) {
	next := node.NextSibling
	if next == nil {
		// 开始判断该批兄弟节点
		return
	}
	defer iteratorNext(next, brothers)
	iteratorChildren(next)
}

// 迭代孩子节点
func iteratorChildren(node *html.Node) {
	child := node.FirstChild
	if child == nil {
		return
	}
	var brothers []*html.Node
	iteratorNext(child, brothers)
}

func analysis(brothers []*html.Node) {
	if len(brothers) == 0 {
		return
	}

	attrNames := make(map[string]int)
	attrValues := make(map[string]int)

	for _, v := range brothers {
		for _, a := range v.Attr {
			key := fmt.Sprintf("%s/%s", a.Key, a.Val)
			attrValues[key] = attrNames[key] + 1
			attrNames[a.Key] = attrNames[a.Key] + 1
		}
	}

}

// 用goquery解析html文档
func parseCssDocument(content []byte) (*goquery.Document, error) {
	return goquery.NewDocumentFromReader(bytes.NewBuffer(content))
}

// 用goquery解析html文档
func parseXpathDocument(content []byte) (*html.Node, error) {
	return htmlquery.Parse(bytes.NewBuffer(content))
}
