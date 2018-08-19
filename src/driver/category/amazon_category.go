package drivercategory

import (
	"curl"
	"net/url"
	"schema"
	"ssdb"

	"github.com/anaskhan96/soup"
	"github.com/golang/glog"
)

func GetCategoryLevel() {
	target := "https://www.amazon.co.jp/gp/site-directory"
	rdata, err := curl.GetURLData(target)
	if err != nil {
		glog.Errorf("Curl Error : %+v", err)
	}

	doc := soup.HTMLParse(rdata)
	root := doc.FindAll("div", "class", "popover-grouping")
	if len(root) == 0 {
		glog.Errorf("Nil Product Page : %s", target)
	}

	data := make(map[string]interface{})
	for _, sub := range root {
		sub_root := sub.Find("h2")
		if sub_root.Pointer != nil {
			h2 := sub_root.Text()
			if _, ok := schema.BlockCategory[h2]; !ok {
				links := sub.FindAll("a")
				for _, element := range links {
					res, err := url.Parse(element.Attrs()["href"])
					if err != nil {
						glog.Warningf("Url Parse Error : %+v", err)
					}
					if len(res.Query()["node"]) > 0 {
						data[res.Query()["node"][0]] = res.EscapedPath()
					}
				}
			}
		}
	}
	ssdbtool.SSDBPool.SetCate(1, data, "")
}
