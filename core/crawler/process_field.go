package crawler

import (
	"digger/models"
	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/hetianyi/gox/logger"
	"golang.org/x/net/html"
	"strings"
)

// 使用css selector处理stage的字段
func processCssField(field *models.Field, s *goquery.Selection) string {
	ret := ""
	if field.IsArray {
		var arrayFieldValue []string
		var sel = s
		if field.Css != "" {
			sel = s.Find(field.Css)
		}
		// 如果不是list类型，则直接匹配fields
		// 循环fields，对于list的每个element进行处理
		sel.Each(func(i int, selection *goquery.Selection) {
			v := ""
			if field.Attr == "" {
				if field.IsHtml {
					v, _ = selection.Html()
				} else {
					v = strings.TrimSpace(selection.Text())
				}
			} else {
				if val, exists := selection.Attr(field.Attr); exists {
					v = strings.TrimSpace(val)
				}
			}
			arrayFieldValue = append(arrayFieldValue, v)
		})
		ret, _ = json.MarshalToString(arrayFieldValue)
	} else {
		var sel = s
		if field.Css != "" {
			sel = s.Find(field.Css)
		}
		v := ""
		if field.Attr == "" {
			if field.IsHtml {
				v, _ = sel.Html()
			} else {
				v = strings.TrimSpace(sel.Text())
			}
		} else {
			if val, exists := sel.Attr(field.Attr); exists {
				v = strings.TrimSpace(val)
			}
		}
		ret = v
	}
	return ret
}

// 使用css selector处理stage的字段
func processXpathField(field *models.Field, node *html.Node) string {
	ret := ""
	if field.IsArray {
		var arrayFieldValue []string
		if field.Xpath != "" {
			list, err := htmlquery.QueryAll(node, field.Xpath)
			if err != nil {
				logger.Error(err)
				return ""
			}
			for _, item := range list {
				if field.IsHtml {
					arrayFieldValue = append(arrayFieldValue, htmlquery.OutputHTML(item, false))
				} else {
					arrayFieldValue = append(arrayFieldValue, htmlquery.InnerText(item))
				}
			}
			ret, _ = json.MarshalToString(arrayFieldValue)
		} else {
			arrayFieldValue = append(arrayFieldValue, strings.TrimSpace(node.Data))
		}
		ret, _ = json.MarshalToString(arrayFieldValue)
	} else {
		if field.Xpath != "" {
			item := htmlquery.FindOne(node, field.Xpath)
			//fmt.Println(htmlquery.OutputHTML(node, false))
			if item != nil {
				if field.IsHtml {
					ret = htmlquery.OutputHTML(item, false)
				} else {
					ret = strings.TrimSpace(htmlquery.InnerText(item))
				}
			}
		} else {
			if field.IsHtml {
				ret = htmlquery.OutputHTML(node, false)
			} else {
				ret = strings.TrimSpace(htmlquery.InnerText(node))
			}
		}
	}
	//fmt.Println(ret)
	return ret
}
