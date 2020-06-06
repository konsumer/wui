package main

import (
	"fmt"
	"net"
	"net/http"

	"github.com/zserge/webview"
)

func main() {
	fs := http.FileServer(http.Dir("./app"))
	http.Handle("/", fs)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	fmt.Println("listening on http://" + listener.Addr().String())

	http.Serve(listener, nil)

	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Title Here")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://" + listener.Addr().String())
	w.Run()
}
