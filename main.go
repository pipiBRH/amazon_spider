package main

import (
	"driver/links"
	"sync"

	"driver/category"
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

	// startCategory(1)
	startProduct(2)

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
