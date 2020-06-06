run: server/main.go
	cd create-wui && go run ../server/main.go

create-wui/wui-linux-amd64: server/main.go
	GOOS=linux GOARCH=amd64 go build -o create-wui/wui-linux-amd64 ./server/main.go

create-wui/wui-linux-arm: server/main.go
	GOOS=linux GOARCH=arm go build -o create-wui/wui-linux-arm ./server/main.go

create-wui/wui-linux-arm64: server/main.go
	GOOS=linux GOARCH=arm64 go build -o create-wui/wui-linux-arm64 ./server/main.go

create-wui/wui-windows: server/main.go
	GOOS=windows GOARCH=amd64 go build -o create-wui/wui-windows ./server/main.go

create-wui/wui-osx: server/main.go
	GOOS=windows GOARCH=amd64 go build -o create-wui/wui-osx ./server/main.go

clean:
	rm -f create-wui/wui-*