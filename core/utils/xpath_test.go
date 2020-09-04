///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils_test

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"strings"
	"testing"
)

func TestXPath1(t *testing.T) {
	node, _ := htmlquery.LoadURL("http://www.shumeipai.net/resource.php?mod=view&rid=466")
	nodes, err := htmlquery.QueryAll(node, "//div[@class='single']/a/text()")
	fmt.Println(err)
	for _, n := range nodes {
		fmt.Println(n.Data)
	}
}

func TestXPath2(t *testing.T) {
	node, _ := htmlquery.LoadURL("http://www.shumeipai.net/resource.php?mod=view&rid=466")
	nodes, err := htmlquery.QueryAll(node, "//div/span[text()='\n30']/../a/text()")
	fmt.Println(err)
	for _, n := range nodes {
		fmt.Println(n.Data)
		fmt.Println(htmlquery.InnerText(n))
		fmt.Println(htmlquery.OutputHTML(n, false))
	}
}

func TestXPath3(t *testing.T) {
	node, _ := htmlquery.LoadDoc("C:\\Users\\Jason\\Downloads\\123.html")
	fmt.Println(htmlquery.OutputHTML(node, false))
	nodes, _ := htmlquery.QueryAll(node, "//table[@class='table']/tbody/tr/th[text()='医保地区']/../td/span/child::text()")
	for _, n := range nodes {
		fmt.Println(strings.TrimSpace(htmlquery.InnerText(n)))
	}
}

func TestXPath5(t *testing.T) {
	node, _ := htmlquery.LoadDoc("C:\\Users\\Jason\\Downloads\\123.html")
	fmt.Println(htmlquery.OutputHTML(node, false))
	nodes := htmlquery.FindOne(node, "//table[@class='table']/tbody/tr/th[text()='医保地区']/../td/span/child::text()[2]")
	fmt.Println(strings.TrimSpace(htmlquery.InnerText(nodes)))
}

func TestXPath6(t *testing.T) {
	node, _ := htmlquery.LoadDoc("C:\\Users\\Jason\\Downloads\\6.html")
	fmt.Println(htmlquery.OutputHTML(node, false))
	nodes, err := htmlquery.QueryAll(node, "//div[contains(@class, \"u-cover\")]")
	fmt.Println(err)
	for _, n := range nodes {
		fmt.Println(strings.TrimSpace(htmlquery.InnerText(n)))
	}
}

func TestXPath7(t *testing.T) {
	node, _ := htmlquery.LoadDoc("C:\\Users\\Jason\\Downloads\\7.html")
	//fmt.Println(htmlquery.OutputHTML(node, false))
	nodes := htmlquery.FindOne(node, "//p[@id='album-desc-more']")
	fmt.Println(strings.TrimSpace(htmlquery.OutputHTML(nodes, false)))
	fmt.Println(strings.TrimSpace(htmlquery.InnerText(nodes)))
}
