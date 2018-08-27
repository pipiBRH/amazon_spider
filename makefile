all: clean build rmlog

clean:
	vgo clean

build: 
	vgo build

rmlog:
	rm ./log/*

linux: clean rmlog
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 vgo build -o linux_amazon_spider	

dep:
	go get github.com/seefan/gossdb
	go get github.com/BurntSushi/toml
	go get github.com/PuerkitoBio/goquery
	go get github.com/golang/glog

