GOPATH:=$(CURDIR)
export GOPATH

all: build clean

clean:
	go clean
	rm ./log/*
build: 
	go build