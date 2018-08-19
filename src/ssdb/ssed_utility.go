package ssdbtool

import (
	"fmt"

	"github.com/golang/glog"
)

func (this *ConnectionPool) SetCate(level int, data map[string]interface{}, parent string) error {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	err := c.Hset("all_level", fmt.Sprint(level), interface{}(1))
	if err != nil {
		glog.Errorf("SSDB set level error : %+v", err)
		return err
	}

	err = c.MultiHset(fmt.Sprintf("category_lv%v", level), data)
	if err != nil {
		glog.Errorf("SSDB category MultiHset error : %+v", err)
		return err
	}

	relation := make(map[string]interface{})
	enable := make(map[string]interface{})
	for key := range data {
		relation[key] = parent
		enable[key] = 1
	}

	err = c.MultiHset(fmt.Sprintf("category_relation_lv%v", level), relation)
	if err != nil {
		glog.Errorf("SSDB category relation MultiHset error : %+v", err)
		return err
	}

	err = c.MultiHset(fmt.Sprintf("category_enable_lv%v", level), enable)
	if err != nil {
		glog.Errorf("SSDB category enable MultiHset error : %+v", err)
		return err
	}

	return nil
}

func (this *ConnectionPool) SetTailCate(data map[string]interface{}) error {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	err := c.MultiHset("category_tail", data)
	if err != nil {
		glog.Errorf("SSDB category tail MultiHset error : %+v", err)
		return err
	}

	enable := make(map[string]interface{})
	for key := range data {
		enable[key] = 1
	}

	err = c.MultiHset("category_enable_tail", enable)
	if err != nil {
		glog.Errorf("SSDB MultiHset error : %+v", err)
		return err
	}

	return nil
}
