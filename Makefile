####

.PHONY: all test dev build

all:

test:

install-dep:
	go get github.com/domluna/watcher

build: install-dep
	go build

dev:
	./wat 'make build' *.go
