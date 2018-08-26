package main

import (
	_ "amazon_spider/src/cmdline"
	"amazon_spider/src/driver/category"
	"amazon_spider/src/driver/links"
	"amazon_spider/src/schema"
	"amazon_spider/src/ssdb"
	"sync"

	"github.com/golang/glog"
)

func main() {
	defer glog.Flush()
	glog.Info("Spider GO!")

	defer ssdbtool.SSDBPool.ClosePool()

	if schema.Config.Spider.EnableCategory {
		glog.Info("Category Spider start")
		startCategory(schema.Config.Spider.CategoryLevel)
	}

	if schema.Config.Spider.EnableProduct {
		glog.Info("Product Spider start")
		startProduct(schema.Config.Spider.Znum)
	}
}

func startCategory(level int) {
	for drivercategory.GetCategoryLevel(level) {
		level++
	}
}

func startProduct(znum int) {
	var wg sync.WaitGroup
	wg.Add(znum)
	driverlink.GetProductLinks(znum, &wg)
	wg.Wait()
}
