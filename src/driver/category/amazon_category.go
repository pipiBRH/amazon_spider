package drivercategory

import (
	"curl"
	"fmt"
	"net/url"
	"schema"
	"ssdb"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
)

func GetCategoryLevel1() {
	target := "https://www.amazon.co.jp/gp/site-directory"
	rdata, err := curl.GetURLData(target)
	if err != nil {
		glog.Errorf("Curl Error : %+v\n", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(rdata))
	if err != nil {
		glog.Errorf("goquery parser error : %+v\n", err)
		return
	}

	root := doc.Find(".popover-grouping")
	if root.Size() == 0 {
		glog.Errorf("Nil Product Page : %s", target)
		return
	}

	data := make(map[string]interface{})

	root.Each(func(index int, sub *goquery.Selection) {
		subRoot := sub.Find("h2")
		if _, ok := schema.BlockCategory[subRoot.Text()]; !ok {
			sub.Find("a").Each(func(index int, sub *goquery.Selection) {
				if value, check := sub.Attr("href"); check {
					res, err := url.Parse(value)
					if err != nil {
						glog.Warningf("Url Parse Error : %+v\n", err)
					} else {
						if len(res.Query()["node"]) > 0 {
							node := res.Query()["node"][0]
							categoryUrl := fmt.Sprintf("https://www.amazon.co.jp%v?node=%v", res.EscapedPath(), node)
							data[node] = categoryUrl
						}
					}
				}
			})
		}
	})

	ssdbtool.SSDBPool.SetCate(1, data, "")
}

func GetCategoryLevel(level int) {
	tail := make(map[string]interface{})
	parentLevel := level - 1
	levelData, err := ssdbtool.SSDBPool.GetCategoryLinks(parentLevel)
	if err != nil {
		glog.Errorf("get level_%v links error  : %+v\n", parentLevel, err)
		return
	}

	for parentNode, links := range levelData {
		data := make(map[string]interface{})
		rdata, err := curl.GetURLData(string(links))
		if err != nil {
			glog.Errorf("Curl links : %v   Error : %+v\n", links, err)
			continue
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(rdata))
		if err != nil {
			glog.Errorf("Parser links : %v   Error : %+v\n", links, err)
			continue
		}
		root := doc.Find(`[id^="categoryTiles"]`)
		if root.Size() < 1 {
			root = doc.Find(`[id^="contentGrid"]`)
		}
		if root.Size() < 1 {
			tail[parentNode] = string(links)
			continue
		}
		root.Find(`a`).Each(func(index int, sub *goquery.Selection) {
			if value, check := sub.Attr("href"); check {
				res, err := url.Parse(value)
				if err != nil {
					glog.Warningf("Url Parse Error : %+v\n", err)
				} else {
					if len(res.Query()["node"]) > 0 {
						node := res.Query()["node"][0]
						categoryUrl := fmt.Sprintf("https://www.amazon.co.jp%v?node=%v", res.EscapedPath(), node)
						data[node] = categoryUrl
					}
				}
			}
		})
		err = ssdbtool.SSDBPool.SetCate(level, data, parentNode)
		if err != nil {
			glog.Warningf("links Warning : %+v\n", links)
			tail[parentNode] = string(links)
		}
	}
	ssdbtool.SSDBPool.SetTailCate(tail)
}
