package driverlink

import (
	"curl"
	"fmt"
	"net/url"

	"github.com/anaskhan96/soup"
	"github.com/golang/glog"
)

func tamp() {
	for index := 399; index < 400; index++ {
		target :=
			fmt.Sprintf("https://www.amazon.co.jp/b/ref=s9_acss_bw_cg_finnbdt_3b1_w?node=393994011&page=%v",
				index)
		rdata, err := curl.GetURLData(target)
		if err != nil {
			glog.Errorf("Curl Error : %+v", err)
		}

		doc := soup.HTMLParse(rdata)
		root := doc.FindAll("div", "class", "s-item-container")
		if len(root) == 0 {
			glog.Errorf("Nil Product Page : %s", target)
		}
		for _, sub := range root {
			sub_root := sub.Find("a", "class", "s-access-detail-page")
			if sub_root.Pointer != nil {
				href := sub_root.Attrs()["href"]
				res, err := url.Parse(href)
				if err != nil {
					glog.Warningf("Url Parse Error : %+v", err)
				}
				fmt.Println(res.Hostname(), res.EscapedPath())
			}
		}
		fmt.Printf("page : %v\n", index)
	}
}
