package curl

import (
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
)

func GetURLData(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Errorf("New Request Error : %+v", err)
		return "", err
	}

	req.Header.Add(
		"User-Agent",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_13_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		glog.Errorf("Get Request Error : %+v", err)
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		glog.Errorf("Read Body Error : %+v", err)
		return "", err
	}

	return string(body), nil
}
