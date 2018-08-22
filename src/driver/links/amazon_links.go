package driverlink

import (
	"curl"
	"fmt"
	"net/url"
	"ssdb"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
)

func GetProductLinks(znum int, wg *sync.WaitGroup) {
	ssdbtool.SSDBPool.SetLinkQueue()
	for index := 0; index < znum; index++ {
		go start(wg)
	}
}

func start(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		targetUrl, err := ssdbtool.SSDBPool.GetQueueLink()

		if err != nil {
			glog.Warningf("get target links error  : %+v", err)
			break
		}
		if targetUrl == "" {
			size, err := ssdbtool.SSDBPool.GetQueueSize()
			if err != nil {
				glog.Errorf("get size error")
			}
			if size == 0 {
				glog.Warningln("tail queue empty")
				break
			}
		}

		rdata, err := curl.GetURLData(targetUrl)
		if err != nil {
			glog.Errorf("Curl Error : %+v", err)
		}
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(rdata))
		if err != nil {
			glog.Errorf("Parser links : %v\n   Error : %+v", targetUrl, err)
			continue
		}

		totalPage, err := strconv.ParseInt(strings.Trim(doc.Find(".pagnDisabled").First().Text(), " "), 10, 32)
		if err != nil {
			glog.Warningf("get page links : %v\n   Error : %+v", targetUrl, err)
			continue
		}
		for sp := 1; sp < int(totalPage); sp++ {
			target := fmt.Sprintf("%v&page=%v", targetUrl, sp)
			rdata, err := curl.GetURLData(target)
			if err != nil {
				glog.Errorf("Curl Error : %+v", err)
				continue
			}

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(rdata))
			if err != nil {
				glog.Errorf("Parser links : %v\n   Error : %+v", targetUrl, err)
				continue
			}

			root := doc.Find(".s-item-container")
			if root.Size() == 0 {
				glog.Warningf("Nil Product Page : %s", target)
				break
			}
			root.Each(func(i int, s *goquery.Selection) {
				s.Find("s-access-detail-page").Each(func(subI int, sub *goquery.Selection) {
					result, ok := sub.Attr("href")
					if ok {
						res, err := url.Parse(result)
						if err != nil {
							glog.Warningf("Url Parse Error : %+v", err)
						}
						fmt.Println(res.Hostname(), res.EscapedPath())
					}
				})
			})
		}
	}
}
