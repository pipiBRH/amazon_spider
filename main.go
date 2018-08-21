package main

import (
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
	// drivercategory.GetCategoryLevel1()
	drivercategory.GetCategoryLevel(3)
}
