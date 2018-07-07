package main

import (
	"net/http"
	"strconv"
	"fmt"
	"sync"
	"context"
	"html/template"
)

type WebUI struct {
	App     *App
	Srv     *http.Server
	MsgChan chan string
}

func (ui *WebUI) Init(config *Config, msgChan chan string) {
	srv := &http.Server{Addr: ":" + strconv.Itoa(config.WebUIPort)}
	ui.Srv = srv
	ui.MsgChan = msgChan
	http.HandleFunc("/", ui.WebUI)
	http.HandleFunc("/shutdown", ui.ShutdownButton)
	http.HandleFunc("/load", ui.loadPreset)
	http.HandleFunc("/play", ui.PlayButton)
	http.HandleFunc("/pause", ui.PauseButton)
}

func (ui *WebUI) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	if err := ui.Srv.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func (ui *WebUI) WebUI(writer http.ResponseWriter, request *http.Request) {
	tmpl, _ := template.ParseFiles("./templates/webui.html")
	tmpl.Execute(writer, ui.App)
}

func (ui *WebUI) PlayButton(writer http.ResponseWriter, request *http.Request) {
	ui.MsgChan <- PlayCmd
	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}

func (ui *WebUI) PauseButton(writer http.ResponseWriter, request *http.Request) {
	ui.MsgChan <- PauseCmd
	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}

func (ui *WebUI) ShutdownButton(writer http.ResponseWriter, request *http.Request) {
	close(ui.MsgChan)
	if err := ui.Srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}

func (ui *WebUI) loadPreset(writer http.ResponseWriter, request *http.Request){
	values, _ := request.URL.Query()["preset"]
	//TODO check errors
	msg := fmt.Sprintf("preset:%s", values[0])
	ui.MsgChan <- msg
	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}
