///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package crawler_test

import (
	"bytes"
	"digger/crawler"
	"digger/models"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"os"
	"testing"
)

func TestProcess(t *testing.T) {
	crawler.Process(&models.Queue{}, nil, os.Stdout, func(cxt *models.Context, oldQueue *models.Queue, newQueue []*models.Queue, results []*models.Result, err error) {

	})
}

func TestProcess1(t *testing.T) {
	c := colly.NewCollector()
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	/*c.OnHTML("div#bd_auto>div.storey_one>div.storey_one_left>div.sidemenu>div.flq_body>div.level_one:nth-child(3)>div.submenu", func(element *colly.HTMLElement) {
		element.ForEach("dd>a", func(i int, element *colly.HTMLElement) {
			val := element.Attr("href")
			fmt.Println(val)
		})
	})*/

	c.OnResponse(func(response *colly.Response) {
		doc, _ := goquery.NewDocumentFromReader(bytes.NewBuffer(response.Body))
		doc.Find("div#bd_auto>div.storey_one>div.storey_one_left>div.sidemenu>div.flq_body>div.level_one:nth-child(3)>div.submenu").Each(func(i int, s *goquery.Selection) {
			val, _ := s.Attr("class")
			fmt.Println(val)
		})
	})
	c.Visit("http://book.dangdang.com/")
}

func TestProcessPageCss(t *testing.T) {

}
