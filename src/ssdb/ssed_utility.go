package ssdbtool

import (
	"errors"
	"fmt"

	"github.com/golang/glog"
	"github.com/seefan/gossdb"
)

func (this *ConnectionPool) SetCate(level int, data map[string]interface{}, parent string) error {
	if len(data) < 1 {
		glog.Errorln("SSDB set data nil")
		return errors.New("set data nil")
	}
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	err := c.Hset("all_level", fmt.Sprint(level), interface{}(1))
	if err != nil {
		glog.Errorf("SSDB set level error => %+v\n  data => %+v", err, data)
		return err
	}

	err = c.MultiHset(fmt.Sprintf("category_lv%v", level), data)
	if err != nil {
		glog.Errorf("SSDB category MultiHset error => %+v\n  data : %+v", err, data)
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
		glog.Errorf("SSDB category relation MultiHset error => %+v\n  data => %+v", err, relation)
		return err
	}

	err = c.MultiHset(fmt.Sprintf("category_enable_lv%v", level), enable)
	if err != nil {
		glog.Errorf("SSDB category enable MultiHset error => %+v\n  data => %+v", err, enable)
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
		glog.Errorf("SSDB category tail MultiHset error => %+v data => %+v", err, data)
		return err
	}

	enable := make(map[string]interface{})
	for key := range data {
		enable[key] = 1
	}

	err = c.MultiHset("category_enable_tail", enable)
	if err != nil {
		glog.Errorf("SSDB MultiHset error => %+v\n data => %+v", err, enable)
		return err
	}

	return nil
}

func (this *ConnectionPool) GetCategoryLinks(level int) (map[string]gossdb.Value, error) {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	enableKey, _, err := c.MultiHgetAllSlice(fmt.Sprintf("category_enable_lv%v", level))
	if err != nil {
		glog.Errorf("SSDB get enable category links level_%v error => %+v", level, err)
		return nil, err
	}

	result, err := c.MultiHget(fmt.Sprintf("category_lv%v", level), enableKey...)
	if err != nil {
		glog.Errorf("SSDB get category links level_%v error => %+v", level, err)
		return nil, err
	}

	return result, nil
}

func (this *ConnectionPool) GetTailLinks() ([]string, error) {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	enableKey, _, err := c.MultiHgetAllSlice("category_enable_tail")
	if err != nil {
		glog.Errorf("SSDB get enable tail links error => %+v", err)
		return nil, err
	}

	result, _, err := c.MultiHgetSlice("category_tail", enableKey...)
	if err != nil {
		glog.Errorf("SSDB get tail category links error => %+v", err)
		return nil, err
	}

	return result, nil
}

func (this *ConnectionPool) SetLinkQueue() error {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	err := c.Qclear("tail_links")
	if err != nil {
		glog.Errorf("SSDB clear tail queue error => %+v", err)
		return err
	}

	tail, err := this.GetTailLinks()
	if err != nil {
		return err
	}

	convTail := make([]interface{}, len(tail))
	for key, _ := range tail {
		convTail[key] = tail[key]
	}

	size, err := c.Qpush_array("tail_links", convTail)
	if err != nil {
		glog.Errorf("SSDB push tail queue error => %+v", err)
		return err
	}
	glog.Infof("SSDB push tail queue size : %v", size)
	return nil
}

func (this *ConnectionPool) GetQueueLink() (string, string, int, error) {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	tailKey, err := c.Qpop_back("tail_links")
	if err != nil {
		glog.Errorf("SSDB pop queue key error => %+v", err)
		return "", "", 1, err
	}
	tk := tailKey.String()
	if tk == "" {
		glog.Errorln("SSDB tail queue empty")
		return "", "", 1, nil
	}

	data, err := c.Hget("category_tail", tk)
	if err != nil {
		glog.Errorf("SSDB hget tail error => %+v", err)
		return "", tk, 1, err
	}

	if data.String() == "" {
		glog.Errorf("SSDB hget tail nil, key => %v", tk)
		return "", tk, 1, nil
	}

	page, err := c.Hget("page_log", tk)
	if err != nil {
		glog.Errorf("SSDB get page node => %v error => %+v", tk, err)
		return "", tk, 1, err
	}

	if page.Int() == 0 || page.Int() >= 400 {
		this.ClearPageLog(tk)
		return data.String(), tk, 1, nil
	}

	return data.String(), tk, page.Int(), nil
}

func (this *ConnectionPool) GetQueueSize() (int64, error) {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	size, err := c.Qsize("tail_links")
	if err != nil {
		glog.Errorf("SSDB get qsize error => %+v", err)
		return 0, err
	}
	return size, nil
}

func (this *ConnectionPool) GetLevelSize(level int) (bool, error) {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	size, err := c.Hsize(fmt.Sprintf("category_enable_lv%v", level))
	if err != nil {
		glog.Errorf("SSDB get hsize level_%v error => %+v", level, err)
		return false, err
	}

	if size > 0 {
		return true, nil
	}

	return false, nil
}

func (this *ConnectionPool) SavePageLog(node string, page int) {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	err := c.Hset("page_log", node, page)
	if err != nil {
		glog.Errorf("SSDB save page error => %+v", err)
	}
}

func (this *ConnectionPool) ClearPageLog(node string) {
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	err := c.Hdel("page_log", node)
	if err != nil {
		glog.Errorf("SSDB clear page error => %+v", err)
	}
}

func (this *ConnectionPool) SetProductLink(data map[string]interface{}) error {
	if len(data) < 1 {
		glog.Errorln("SSDB set product link nil")
	}
	c := SSDBPool.GetSSDBClient()
	defer c.Close()

	err := c.MultiHset("product_links", data)
	if err != nil {
		glog.Errorf("SSDB set product link error => %+v", err)
		return err
	}
	return nil
}
