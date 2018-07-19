package main

import (
	"net/http"
	"strconv"
	"html/template"
	"fmt"
	"sync"
)

type WebUI struct {
	Srv     *http.Server
	Presets Presets
	MsgChan chan string
}

const (
	PlayCmd   = "play"
	PauseCmd  = "pause"
	StopCmd   = "stop"
	LoadCmd   = "load"
	PresetCmd = "preset:%s"
)

// Create new WebUI instance
func NewWebUI(config *Config, presets Presets, ch chan string) *WebUI {
	ui := WebUI{}
	ui.Presets = presets
	ui.MsgChan = ch
	port := strconv.Itoa(config.WebUIPort)
	ui.newServer(port)
	return &ui
}

// Create new http server
func (ui *WebUI) newServer(port string) {
	srv := &http.Server{Addr: ":" + port}
	http.HandleFunc("/", ui.MainPage)
	http.HandleFunc("/load", ui.LoadPreset)
	http.HandleFunc("/play", ui.PlayButton)
	http.HandleFunc("/pause", ui.PauseButton)
	http.HandleFunc("/stop", ui.StopButton)
	ui.Srv = srv
}

func (ui *WebUI) StartServer(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := ui.Srv.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

// Main page
func (ui *WebUI) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/webui.html")
	tmpl.Execute(w, ui.Presets)
}

// Load preset request
func (ui *WebUI) LoadPreset(w http.ResponseWriter, r *http.Request) {
	values, _ := r.URL.Query()["preset"]
	//TODO check errors
	msg := fmt.Sprintf(PresetCmd, values[0])
	go func() {
		ui.MsgChan <- msg
	}()
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// Play button
func (ui *WebUI) PlayButton(w http.ResponseWriter, r *http.Request) {
	go func() {
		ui.MsgChan <- PlayCmd
	}()
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// Pause button
func (ui *WebUI) PauseButton(w http.ResponseWriter, r *http.Request) {
	go func() {
		ui.MsgChan <- PauseCmd
	}()
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// Stop button
func (ui *WebUI) StopButton(w http.ResponseWriter, r *http.Request) {
	go func() {
		ui.MsgChan <- StopCmd
	}()
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
