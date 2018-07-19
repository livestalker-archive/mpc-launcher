package main

import (
	"sync"
	"os/exec"
	"fmt"
	"net/url"
	"net/http"
	"strings"
	"strconv"
)

type App struct {
	Config  *Config
	Presets Presets
	WebUI   *WebUI
	Cmds    []*exec.Cmd
}

func (app *App) Init(configFilename string, presetsFilename string) {
	msgChan := make(chan string)
	app.Config = LoadConfig(configFilename)
	app.Presets = LoadPresets(presetsFilename)
	ui := &WebUI{}
	ui.Init(app.Config, msgChan)
	ui.App = app
	app.WebUI = ui
}

func (app *App) ExecuteMsg(msg string) {
	for i := 0; i < app.Config.MonCount; i++ {
		address := fmt.Sprintf("http://localhost:%d/command.html", app.Config.StartPort+i+1)
		data := url.Values{}
		data.Add("wm_command", msg)
		http.PostForm(address, data)
		//TODO process error
	}
}

func (app *App) LoadPreset(id int) {
	if id+1 > len(app.Presets) {
		return
	}
	preset := app.Presets[id]
	app.Cmds = make([]*exec.Cmd, app.Config.MonCount)
	for i := 0; i < app.Config.MonCount; i++ {
		app.Cmds[i] = exec.Command(app.Config.MpcPath, preset.Files[i].GetFullArgs()...)
		app.Cmds[i].Start()
	}
}

func main() {
	configFilename := "./conf/config.yml"
	presetsFilename := "./conf/presets.yml"
	app := App{}
	app.Init(configFilename, presetsFilename)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go supervisor(wg, &app)
	wg.Add(1)
	go app.WebUI.Start(wg)
	wg.Wait()
}

func supervisor(wg *sync.WaitGroup, app *App) {
	defer wg.Done()
	for msg := range app.WebUI.MsgChan {
		if strings.HasPrefix(msg, "preset") {
			elements := strings.Split(msg, ":")
			//TODO process error
			id, _ := strconv.Atoi(elements[1])
			app.LoadPreset(id)
		} else {
			app.ExecuteMsg(msg)
		}
	}
	for _, cmd := range app.Cmds {
		cmd.Wait()
	}
}
