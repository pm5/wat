
.PHONY: all install-dep build test dev

all: install-dep build

install-dep:
	go get github.com/domluna/watcher

build:
	go build

test:

dev: build
	./wat 'make build' 'make test' -- *.go Makefile
