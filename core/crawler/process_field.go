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
func processCssField(cxt *models.Context, field *models.Field, s *goquery.Selection) interface{} {
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
			// slot s4
			v = handleS4(cxt, field, field.Name, v)
			arrayFieldValue = append(arrayFieldValue, v)
		})
		return arrayFieldValue
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
		// slot s4
		return handleS4(cxt, field, field.Name, v)
	}
}

// 使用css selector处理stage的字段
func processXpathField(cxt *models.Context, field *models.Field, node *html.Node) interface{} {
	if field.IsArray {
		var arrayFieldValue []string
		if field.Xpath != "" {
			list, err := htmlquery.QueryAll(node, field.Xpath)
			if err != nil {
				logger.Error(err)
				return ""
			}
			for _, item := range list {
				v := ""
				if field.IsHtml {
					v = htmlquery.OutputHTML(item, false)
				} else {
					v = htmlquery.InnerText(item)
				}
				// slot s4
				v = handleS4(cxt, field, field.Name, v)
				arrayFieldValue = append(arrayFieldValue, v)
			}
			return arrayFieldValue
		} else {
			// TODO 提取公共代码
			v := ""
			if field.IsHtml {
				v = htmlquery.OutputHTML(node, false)
			} else {
				v = strings.TrimSpace(htmlquery.InnerText(node))
			}
			// slot s4
			v = handleS4(cxt, field, field.Name, v)
			arrayFieldValue = append(arrayFieldValue, strings.TrimSpace(v))
		}
		return arrayFieldValue
	} else {
		v := ""
		if field.Xpath != "" {
			item := htmlquery.FindOne(node, field.Xpath)
			if item != nil {
				if field.IsHtml {
					v = htmlquery.OutputHTML(item, false)
				} else {
					v = strings.TrimSpace(htmlquery.InnerText(item))
				}
			}
		} else {
			if field.IsHtml {
				v = htmlquery.OutputHTML(node, false)
			} else {
				v = strings.TrimSpace(htmlquery.InnerText(node))
			}
		}
		// slot s4
		return handleS4(cxt, field, field.Name, v)
	}
}
