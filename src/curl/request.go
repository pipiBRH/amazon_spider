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

	req.Header.Add(
		"cookie",
		`session-id=355-8078946-8525209; ubid-acbjp=355-4342028-5983344; x-wl-uid=1pIGHhg2QVTJ2j2Tx0h3Lb5OKLccu0uj5BL2Afof5UekLAU0miMfcNDwG8VFK8bMNUfbaex+c3fg=; session-token="0LJYyXedoNap06mEeF8phHfwCgQnk/2uz6Euiigl3EXRldwVoadlZ90c6aNP3lOpGLu2iFvTjAWOvNSG4HkdYc+XhufA1jvK7pVLlV6Tf4zwodS84bdsugeV1hMOR+9VJrvhy6Ujg3BIUQAD3itKpARk6PyZcY98jQYfwKgJV7/16FrM0i/gaF7E+FJuXrsUEXfrn8b3ysgxYJBcP+kIFQ=="; x-amz-captcha-1=1534921756664722; x-amz-captcha-2=qRINUJ+Wr9d/OXk1tvh7Gg==; amznacsleftnav-d8c684c1-6ab5-3ef9-a6eb-63ba4c0d071d=1; session-id-time=2082787201l; csm-hit=tb:2215R08Q1S1DAX43B18Y+s-2215R08Q1S1DAX43B18Y|1534918284911&adb:adblk_yes`)

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
