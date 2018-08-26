package driverlink

import (
	"amazon_spider/src/curl"
	"amazon_spider/src/schema"
	"amazon_spider/src/ssdb"
	"crypto/md5"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/golang/glog"
)

func GetProductLinks(znum int, wg *sync.WaitGroup) {
	// ssdbtool.SSDBPool.SetLinkQueue()
	for index := 0; index < znum; index++ {
		go start(wg)
	}
}

func start(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		targetUrl, tailKey, pageLog, err := ssdbtool.SSDBPool.GetQueueLink()
		if err != nil {
			break
		}
		if targetUrl == "" {
			size, err := ssdbtool.SSDBPool.GetQueueSize()
			if err != nil {
				break
			}
			if size == 0 {
				glog.Warningln("tail queue empty")
				break
			}
		}

		rdata, err := curl.GetURLDataChrome(targetUrl)
		if err != nil {
			glog.Errorf("Curl Error => %+v", err)
		}

		time.Sleep(time.Microsecond * 100 * schema.Config.Spider.Sleep)

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(rdata))
		if err != nil {
			fmt.Println("DOM:", rdata)
			glog.Errorf("Parser links => %v\n   Error => %+v", targetUrl, err)
			continue
		}
		pageText := doc.Find(".pagnDisabled").First().Text()
		fmt.Println("pageText :", pageText)
		totalPage, err := strconv.ParseInt(strings.Trim(pageText, " "), 10, 32)
		if err != nil {
			glog.Warningf("get page links => %v\n   Error => %+v", targetUrl, err)
			continue
		}

		glog.Infof("target => %v | totalpage => %v", targetUrl, totalPage)

		for sp := pageLog; sp <= int(totalPage); sp++ {
			pdata := make(map[string]interface{})
			target := fmt.Sprintf("%v&page=%v", targetUrl, sp)
			fmt.Println(target)

			rdata, err := curl.GetURLDataChrome(target)
			if err != nil {
				glog.Errorf("Curl Error => %+v", err)
				continue
			}

			time.Sleep(time.Microsecond * 100 * schema.Config.Spider.Sleep)

			doc, err := goquery.NewDocumentFromReader(strings.NewReader(rdata))
			if err != nil {
				glog.Errorf("Parser links => %v\n   Error => %+v", targetUrl, err)
				continue
			}

			root := doc.Find(".s-item-container")
			if root.Size() == 0 {
				glog.Warningf("Nil Product Page => %s", target)
				break
			}
			root.Each(func(i int, s *goquery.Selection) {
				result, ok := s.Find(".s-access-detail-page").First().Attr("href")
				if ok {
					res, err := url.Parse(result)
					if err != nil {
						glog.Warningf("Url Parse Error => %+v", err)
					}
					pid := fmt.Sprintf("%x", md5.Sum([]byte(res.EscapedPath())))
					productUrl := fmt.Sprintf("https://%v%v\n\n", res.Hostname(), res.EscapedPath())
					pdata[pid] = productUrl
				}
			})
			err = ssdbtool.SSDBPool.SetProductLink(pdata)
			if err == nil {
				ssdbtool.SSDBPool.SavePageLog(tailKey, sp+1)
			}
		}
	}
}
