package main

import (
	"sync"
	"os/exec"
	"fmt"
	"net/url"
	"net/http"
)

type App struct {
	Config  *Config
	Presets *Presets
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
	app.Cmds = make([]*exec.Cmd, app.Config.MonCount)
	for i := 0; i < app.Config.MonCount; i++ {
		app.Cmds[i] = exec.Command(app.Config.MpcPath, app.Config.GetNArgs(i+1)...)
	}
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
	for _, cmd := range app.Cmds {
		err := cmd.Start()
		if err != nil {
			//TODO process error
			fmt.Println(err)
		}
	}
	for msg := range app.WebUI.MsgChan {
		app.ExecuteMsg(msg)
	}
	for _, cmd := range app.Cmds {
		cmd.Wait()
	}
}
