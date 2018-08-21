package ssdbtool

import (
	"errors"
	"fmt"

	"github.com/seefan/gossdb"

	"github.com/golang/glog"
)

func (this *ConnectionPool) SetCate(level int, data map[string]interface{}, parent string) error {
	if len(data) < 1 {
		glog.Errorln("SSDB set data nil")
		return errors.New("set data nil")
	}
	c := SSDBPool.GetSSDBClient()
	defer c.Close()
	SSDBPool.resetEnableCategory()

	err := c.Hset("all_level", fmt.Sprint(level), interface{}(1))
	if err != nil {
		glog.Errorf("SSDB set level error : %+v\n  data : %+v\n", err, data)
		return err
	}

	err = c.MultiHset(fmt.Sprintf("category_lv%v", level), data)
	if err != nil {
		glog.Errorf("SSDB category MultiHset error : %+v\n  data : %+v\n", err, data)
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
		glog.Errorf("SSDB category relation MultiHset error : %+v\n  data : %+v\n", err, relation)
		return err
	}

	err = c.MultiHset(fmt.Sprintf("category_enable_lv%v", level), enable)
	if err != nil {
		glog.Errorf("SSDB category enable MultiHset error : %+v\n  data : %+v\n", err, enable)
		return err
	}

	return nil
}

func (this *ConnectionPool) SetTailCate(data map[string]interface{}) error {
	if len(data) < 1 {
		glog.Errorln("SSDB set data nil")
		return errors.New("set data nil")
	}

	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	err := c.MultiHset("category_tail", data)
	if err != nil {
		glog.Errorf("SSDB category tail MultiHset error : %+v\n data : %+v\n", err, data)
		return err
	}

	enable := make(map[string]interface{})
	for key := range data {
		enable[key] = 1
	}

	err = c.MultiHset("category_enable_tail", enable)
	if err != nil {
		glog.Errorf("SSDB MultiHset error : %+v\n data : %+v\n", err, enable)
		return err
	}

	return nil
}

func (this *ConnectionPool) GetCategoryLinks(level int) (map[string]gossdb.Value, error) {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	enableKey, _, err := c.MultiHgetAllSlice(fmt.Sprintf("category_enable_lv%v", level))
	if err != nil {
		glog.Errorf("SSDB get enable category links level_%v error : %+v\n", level, err)
		return nil, err
	}

	result, err := c.MultiHget(fmt.Sprintf("category_lv%v", level), enableKey...)
	if err != nil {
		glog.Errorf("SSDB get category links level_%v error : %+v\n", level, err)
		return nil, err
	}

	return result, nil
}
