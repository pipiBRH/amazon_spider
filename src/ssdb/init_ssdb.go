package ssdbtool

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/seefan/gossdb"
	"github.com/seefan/gossdb/conf"
)

var SSDBPool ConnectionPool

type ConnectionPool struct {
	pool *gossdb.Connectors
}

func InitSSDB() {
	var err error
	SSDBPool.pool, err = gossdb.NewPool(&conf.Config{
		Host:             "127.0.0.1",
		Port:             8888,
		ReadWriteTimeout: 180,
		MinPoolSize:      5,
		MaxPoolSize:      20,
		AcquireIncrement: 5,
	})
	if err != nil {
		glog.Fatalf("SSDB pool init error : %+v", err)
	}

	SSDBPool.resetEnableCategory()
}

func (this *ConnectionPool) ClosePool() {
	this.pool.Close()
}

func (this *ConnectionPool) GetSSDBClient() *gossdb.Client {
	c, err := this.pool.NewClient()
	if err != nil {
		glog.Fatalf("SSDB get client error : %+v", err)
	}
	return c
}

func (this *ConnectionPool) resetEnableCategory() {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	res, err := c.HgetAll("all_level")
	if err != nil {
		glog.Error("SSDB get all level error : %+v", err)
	}

	for level := range res {
		err := c.Hclear(fmt.Sprintf("category_enable_lv%v", level))
		if err != nil {
			glog.Error("SSDB clear enable error : %+v", err)
		}
	}
}
