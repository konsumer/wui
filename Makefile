.PHONY: run deps clean build release

run: server/main.go
	cd create-wui && go run ../server/main.go

deps:
	go get ./server

clean:
	rm -f runtimes/* releases/*

build: runtimes/wui-linux-amd64 runtimes/wui-linux-arm runtimes/wui-linux-arm64 runtimes/wui.exe runtimes/wui-osx

# TODO: target other than my own need some trickery with CGO_ENABLED=1 
release: releases/wui-linux-amd64.gz

runtimes/wui-linux-amd64: server/main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o $@ ./server/main.go
	upx --best $@

releases/wui-linux-amd64.gz: runtimes/wui-linux-amd64
	gzip -f -c $? > $@


runtimes/wui-linux-arm: server/main.go
	GOOS=linux GOARCH=arm go build -ldflags="-s -w" -o $@ $?
	upx --best $@

releases/wui-linux-arm.gz: runtimes/wui-linux-arm
	gzip -f -c $? > $@


runtimes/wui-linux-arm64: server/main.go
	GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $@ $?
	upx --best runtimes/wui-linux-arm64

releases/wui-linux-arm64.gz: runtimes/wui-linux-arm64
	gzip -f -c $? > $@


runtimes/wui.exe: server/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -H windowsgui" -o $@ $?
	upx --best $@

releases/wui-win-amd64.zip: runtimes/wui.exe
	cd runtimes && zip -9 $? $(shell basename "$@") && mv $@ ../releases


runtimes/wui-osx: server/main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $@ $?
	upx --best $@

releases/wui-osx-amd64.zip: runtimes/wui-osx
	cd runtimes && zip -9 $? $(shell basename "$@") && mv $@ ../releases