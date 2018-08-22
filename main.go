package main

import (
	"driver/links"
	"sync"
	// "driver/category"
	"flag"

	"ssdb"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	glog.Info("Spider GO!")

	ssdbtool.InitSSDB()
	defer ssdbtool.SSDBPool.ClosePool()

	// drivercategory.GetCategoryLevel(6)

	var wg sync.WaitGroup
	znum := 1
	wg.Add(znum)
	driverlink.GetProductLinks(znum, &wg)
	wg.Wait()
}
