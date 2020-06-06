package main

import (
	"fmt"
	"net"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/zserge/webview"
	"github.com/imdario/mergo"
)

type Config struct {
  Title string `json:"title"`
  Icon string `json:"icon""`
  Url string `json:"url"`
  Debug bool `json:"debug"`
}

type PayloadFileWrite struct {
  Filename string `json:"filename"`
  Contents string `json:"contents"`
}

type PayloadFile struct {
  Filename string `json:"filename"`
}

type PayloadDir struct {
  Dirname string `json:"dirname"`
}

type PayloadEnv struct {
  Name string `json:"name"`
  Value string `json:"value"`
}

type PayloadExec struct {
  Command string `json:"command"`
}

type PayloadFileRead struct {
  Contents string `json:"contents"`
}

var config Config

func main() {
	r, err := os.Open("./app/settings.json")
	byteValue, _ := ioutil.ReadAll(r)
	
	// defaults
	config = Config {
		Title: "WUI App",
		Icon: "./app.png",
		Url: "/",
	}

	json.Unmarshal(byteValue, &config)

	fs := http.FileServer(http.Dir("./app"))
	http.Handle("/", fs)
	
	http.HandleFunc("/_api/settings", handlSettings)
	http.HandleFunc("/_api/write", handleWrite)
	http.HandleFunc("/_api/read", handleRead)
	http.HandleFunc("/_api/mkdir", handleMkdir)
	http.HandleFunc("/_api/ls", handleLs)
	http.HandleFunc("/_api/stat", handleStat)
	http.HandleFunc("/_api/rm", handleRm)
	http.HandleFunc("/_api/env", handleEnv)
	http.HandleFunc("/_api/exec", handleExec)
	http.HandleFunc("/_api/exit", handleExit)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}

	if config.Debug {
		fmt.Println("Listening on http://" + listener.Addr().String() + config.Url)
	}

	go http.Serve(listener, nil)

	w := webview.New(config.Debug)
	defer w.Destroy()
	w.SetTitle(config.Title)
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("http://" + listener.Addr().String() + config.Url)
	w.Run()
}


// TODO: use webview.Eval & webview.Bind to connect this all, directly?
func handlSettings (w http.ResponseWriter, r *http.Request){
  if (r.Method == http.MethodPost) {
		var newConfig Config
		err := json.NewDecoder(r.Body).Decode(&newConfig)
		if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  	}
  	err = mergo.Merge(&config, newConfig, mergo.WithOverride)
  	if err != nil {
      http.Error(w, err.Error(), http.StatusBadRequest)
      return
  	}
  }
  json.NewEncoder(w).Encode(config)
}

func handleWrite(w http.ResponseWriter, r *http.Request){
	p := PayloadFileWrite {
		Filename: "",
		Contents: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
	}
  err = ioutil.WriteFile(p.Filename, []byte(p.Contents), 0644)
  if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func handleRead(w http.ResponseWriter, r *http.Request){
	p := PayloadFile {
		Filename: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
	}
	contents, err := ioutil.ReadFile(p.Filename)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
	}

	out := PayloadFileRead {
		Contents: string(contents),
	}

	json.NewEncoder(w).Encode(out)
}

func handleMkdir(w http.ResponseWriter, r *http.Request){
	p := PayloadDir {
		Dirname: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
	}
}

func handleLs(w http.ResponseWriter, r *http.Request){
	p := PayloadDir {
		Dirname: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
	}
}

func handleStat(w http.ResponseWriter, r *http.Request){
	p := PayloadFile {
		Filename: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
	}
}

func handleRm(w http.ResponseWriter, r *http.Request){
	p := PayloadFile {
		Filename: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
	}
}

func handleEnv(w http.ResponseWriter, r *http.Request){
	p := PayloadEnv {
		Name: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
	}
}

func handleExec(w http.ResponseWriter, r *http.Request){
	p := PayloadExec {
		Command: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
    http.Error(w, err.Error(), http.StatusBadRequest)
    return
	}
}

func handleExit(w http.ResponseWriter, r *http.Request){
	os.Exit(0)
}