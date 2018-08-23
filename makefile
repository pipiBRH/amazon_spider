GOPATH:=$(CURDIR)
export GOPATH

all: clean build rmlog

clean:
	go clean
build: 
	go build

rmlog:
	rm ./log/*

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o linux_amazon_spider	