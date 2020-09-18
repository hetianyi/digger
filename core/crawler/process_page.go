package crawler

import (
	"digger/models"
	"digger/utils"
	"fmt"
	"github.com/antchfx/htmlquery"
	"strings"
)

// 分页选择器处理
func processPage(cxt *models.Context) (string, error) {
	stage := cxt.Stage
	if stage.PageCss != "" {
		return processPageByCssSelector(cxt)
	}
	if stage.PageXpath != "" {
		return processPageByXpathSelector(cxt)
	}
	// logger.Debug("no page selector")
	return "", nil
}

// css选择器
func processPageByCssSelector(cxt *models.Context) (string, error) {
	stage := cxt.Stage
	doc, err := parseCssDocument(cxt)
	if err != nil {
		return "", nil
	}
	nextPage := ""
	sel := doc.Selection.Find(stage.PageCss)
	if sel == nil {
		return "", nil
	}
	if stage.PageAttr == "" {
		nextPage = sel.Text()
	} else {
		if val, exists := sel.Attr(stage.PageAttr); exists {
			nextPage = strings.TrimSpace(val)
		}
	}
	return handleNextPage(cxt, nextPage), nil
}

// xpath选择器
func processPageByXpathSelector(cxt *models.Context) (string, error) {
	stage := cxt.Stage
	doc, err := parseXpathDocument(cxt)
	if err != nil {
		return "", err
	}
	node := htmlquery.FindOne(doc, stage.PageXpath)
	if node == nil {
		return "", nil
	}
	nextPage := strings.TrimSpace(htmlquery.InnerText(node))
	return handleNextPage(cxt, nextPage), nil
}

func handleNextPage(cxt *models.Context, nextPage string) string {
	stage := cxt.Stage
	if nextPage != "" {
		plugin := stage.FindPlugins("s4")
		if plugin != nil {
			// slot s4
			nextPage = handleStageS4(cxt, stage, nextPage)
		} else {
			nextPage, _ = utils.AbsoluteURL(cxt.Queue.Url, nextPage)
		}
		fmt.Println(fmt.Sprintf("下一页: %s", nextPage))
	}
	return nextPage
}
