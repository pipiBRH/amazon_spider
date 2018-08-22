GOPATH:=$(CURDIR)
export GOPATH

all: build clean

clean:
	rm ./log/*
build: 
	go build