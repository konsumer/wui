package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/imdario/mergo"
	"github.com/zserge/webview"
)

// Config - settings format
type Config struct {
	Title  string `json:"title"`
	Icon   string `json:"icon"`
	URL    string `json:"url"`
	Debug  bool   `json:"debug"`
	Tray   Tray   `json:"tray"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// Tray - definition of Tray-menu
type Tray struct {
	Title string     `json:"title"`
	Icon  string     `json:"icon"`
	Menu  []TrayItem `json:"menu"`
}

// TrayItem - single tray-menu item
type TrayItem struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Tooltip string `json:"tooltip"`
}

// PayloadFileWrite - write to a file
type PayloadFileWrite struct {
	Filename string `json:"filename"`
	Contents string `json:"contents"`
}

// PayloadFile - reference a file
type PayloadFile struct {
	Filename string `json:"filename"`
}

// PayloadDir - reference a directory
type PayloadDir struct {
	Dirname string `json:"dirname"`
}

// PayloadEnv - reference an environment-variable
type PayloadEnv struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// PayloadExec - reference a command
type PayloadExec struct {
	Command string `json:"command"`
}

// PayloadFileRead - reference file-contents
type PayloadFileRead struct {
	Contents string `json:"contents"`
}

// PayloadFileInfo - reference file-info
type PayloadFileInfo struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	Mode    string    `json:"mode"`
	ModTime time.Time `json:"modTime"`
	IsDir   bool      `json:"isDir"`
}

var config Config
var w webview.WebView

func main() {
	r, err := os.Open("./app/settings.json")
	byteValue, _ := ioutil.ReadAll(r)

	// defaults
	config = Config{
		Title:  "WUI App",
		Icon:   "./app.png",
		URL:    "/",
		Width:  800,
		Height: 600,
		Tray:   Tray{},
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
		fmt.Println("Listening on http://" + listener.Addr().String() + config.URL)
	}

	go http.Serve(listener, nil)

	if !strings.HasPrefix(config.URL, "http") {
		config.URL = "http://" + listener.Addr().String() + config.URL
	}

	w = webview.New(config.Debug)
	defer w.Destroy()
	w.SetTitle(config.Title)
	w.SetSize(config.Width, config.Height, webview.HintNone)
	w.Navigate(config.URL)

	if config.Tray.Icon != "" {
		systray.Run(onTrayReady, nil)
	}

	w.Run()
}

func onTrayReady() {
	systray.SetTitle(config.Tray.Title)
	for _, item := range config.Tray.Menu {
		mItem := systray.AddMenuItem(item.Title, item.Tooltip)
		go func(item TrayItem) {
			<-mItem.ClickedCh
			fmt.Println(item.Title)
			w.Navigate(item.URL)
		}(item)
	}
}

// TODO: use webview.Eval & webview.Bind to connect this all, directly?
func handlSettings(w http.ResponseWriter, r *http.Request) {
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
	json.NewEncoder(w).Encode(config)
}

func handleWrite(w http.ResponseWriter, r *http.Request) {
	p := PayloadFileWrite{
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
	fmt.Fprintf(w, "1")
}

func handleRead(w http.ResponseWriter, r *http.Request) {
	p := PayloadFile{
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

	out := PayloadFileRead{
		Contents: string(contents),
	}

	json.NewEncoder(w).Encode(out)
}

func handleMkdir(w http.ResponseWriter, r *http.Request) {
	p := PayloadDir{
		Dirname: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = os.MkdirAll(p.Dirname, 755)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Fprintf(w, "1")
}

func handleLs(w http.ResponseWriter, r *http.Request) {
	p := PayloadDir{
		Dirname: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	files, err := ioutil.ReadDir(p.Dirname)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	out := []string{}

	for _, file := range files {
		out = append(out, file.Name())
	}

	json.NewEncoder(w).Encode(out)
}

func handleStat(w http.ResponseWriter, r *http.Request) {
	p := PayloadFile{
		Filename: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	info, err := os.Stat(p.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	out := PayloadFileInfo{
		Name:  info.Name(),
		Size:  info.Size(),
		Mode:  info.Mode().Perm().String(),
		IsDir: info.IsDir(),
	}
	json.NewEncoder(w).Encode(out)
}

func handleRm(w http.ResponseWriter, r *http.Request) {
	p := PayloadFile{
		Filename: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func handleEnv(w http.ResponseWriter, r *http.Request) {
	p := PayloadEnv{
		Name: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func handleExec(w http.ResponseWriter, r *http.Request) {
	p := PayloadExec{
		Command: "",
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func handleExit(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}
