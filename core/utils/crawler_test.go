///////////////////////////////////////////
// Copyright(C) 2020
// Author : Jason He
// Version: 0.0.1
///////////////////////////////////////////
package utils_test

import (
	"digger/utils"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"testing"
	"time"
)

func TestGoQuery1(t *testing.T) {
	// Request the HTML page.
	res, err := http.Get("http://book.dangdang.com/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("div#bd_auto>div.storey_one>div.storey_one_left>div.sidemenu>div.flq_body>div.level_one:nth-child(3)>div.submenu").Each(func(i int, s *goquery.Selection) {
		s.Find("dd>a").Each(func(i int, selection *goquery.Selection) {
			fmt.Println(selection.Attr("href"))
		})
	})
}

func TestGoQuery2(t *testing.T) {
	// Request the HTML page.
	res, err := http.Get("http://product.dangdang.com/27894120.html")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(doc.Find("ul#main-img-slider li").Children().Size())

	//for s := doc.Find("ul#main-img-slider>li>a"); s != nil; s = s.Next() {
	//	fmt.Println(s.Attr("data-imghref"))
	//}

	// Find the review items
	doc.Find("ul#main-img-slider>li>a").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Attr("data-imghref"))
	})

	time.Sleep(time.Second * 10)
}

func TestColly1(t *testing.T) {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.Visit("http://book.dangdang.com/")
}

func TestGoQuery3(t *testing.T) {
	// Request the HTML page.
	res, err := http.Get("http://book.dangdang.com/")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Find the review items
	doc.Find("div#bd_auto>div.storey_one>div.storey_one_left>div.sidemenu>div.flq_body>div.level_one:nth-child(3)>div.submenu").Each(func(i int, s *goquery.Selection) {
		s.Find("dd>a").Each(func(i int, selection *goquery.Selection) {
			fmt.Println(selection.Attr("href"))
		})
	})
}

func TestParseLabels(t *testing.T) {
	bytes, _ := json.MarshalIndent(utils.ParseLabels("role=worker,sex=,age=123,xxxx,ll=00,,"), "", "  ")
	fmt.Println(string(bytes))
}
