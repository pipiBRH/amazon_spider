package main

import (
	"driver/links"
	"sync"

	"driver/category"
	"flag"
	"schema"
	"ssdb"

	"github.com/BurntSushi/toml"
	"github.com/golang/glog"
)

func main() {

	configPath := flag.String("config", "conf/dev.toml", "specific config file")
	flag.Parse()

	if _, err := toml.DecodeFile(*configPath, &schema.Config); err != nil {
		glog.Fatalf("Parser config error : %+v", err)
	}

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
