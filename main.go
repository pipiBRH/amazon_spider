package main

import (
	"cmdline"
	"driver/category"
	"driver/links"
	"schema"
	"ssdb"
	"sync"

	"github.com/golang/glog"
)

func main() {
	cmd.InitCmd()

	defer glog.Flush()
	glog.Info("Spider GO!")

	ssdbtool.InitSSDB()
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
