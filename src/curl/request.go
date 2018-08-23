package curl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"

	"github.com/golang/glog"
)

func GetURLDataChrome(url string) (string, error) {
	cmd := exec.Command("node", "./chrome.js", url)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		glog.Errorf("CMD run Error : %+v", err)
		return "", err
	}
	return fmt.Sprintf("%s", out.String()), nil
}

func GetURLData(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		glog.Errorf("New Request Error : %+v", err)
		return "", err
	}
	req.Header.Add(
		"authority",
		"www.amazon.co.jp")

	req.Header.Add(
		"User-Agent",
		"User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:61.0) Gecko/20100101 Firefox/61.0")

	req.Header.Add(
		"cookie",
		`session-id=358-9454036-7733051; session-id-time=2082787201l; csm-hit=tb:9G1EAYNFHMH9M9EBV54Q+s-9G1EAYNFHMH9M9EBV54Q|1534998852733&adb:adblk_no; ubid-acbjp=358-4070913-7925309; x-wl-uid=1Unk1H6JskLAcb2dcvudl/aGMb4bakh7nVXyKd642RkZAlUOiUXrvvHW0XVDZ2YqnIeRxc+JX5F8=; session-token=1XO3ZdmrMzwcLtAB3fO1X7zdYu/8J/5yOo8tV5XE1IHzCukmsKqXmSGO0bstnYcODlcDbpN2z65F4leU4cQEeKGKvc42vCjckWKGS9WcNAhPQtq+UzrlT42SKkGvc+z5WQHw6ZsszF8gexWsiFtLm1yYI0SB6buHWNxTG6/owek7BsIN5wZQGg4YX1ZkSpPV`)

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
