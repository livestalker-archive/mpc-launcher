package main

import (
	"sync"
	"os/exec"
	"fmt"
	"net/url"
	"net/http"
	"strings"
	"strconv"
	"time"
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

func (app *App) LoadPreset(id int) {
	if id+1 > len(app.Presets) {
		return
	}
	preset := app.Presets[id]
	client := http.Client{}
	for i := 0; i < app.Config.MonCount; i++ {
		filename := preset.Files[i]
		address := fmt.Sprintf("http://localhost:%d/browser.html", app.Config.StartPort+i+1)
		req, _ := http.NewRequest(http.MethodGet, address, nil)
		q := req.URL.Query()
		q.Add("path", filename)
		req.URL.RawQuery = q.Encode()
		client.Do(req)
		//TODO process error
	}
	time.Sleep(1 * time.Second)
	app.ExecuteMsg(FullscreenCmd)
}

func (app *App) StartCmd() {
	for _, cmd := range app.Cmds {
		err := cmd.Start()
		if err != nil {
			//TODO process error
			fmt.Println(err)
		}
	}
}

func (app *App) StopCmd() {
	for _, cmd := range app.Cmds {
		if cmd.Process != nil {
			err := cmd.Process.Kill()
			if err != nil {
				//TODO process error
				fmt.Println(err)
			}
			cmd.Wait()
		}
	}
	app.Cmds = make([]*exec.Cmd, app.Config.MonCount)
	for i := 0; i < app.Config.MonCount; i++ {
		app.Cmds[i] = exec.Command(app.Config.MpcPath, app.Config.GetNArgs(i+1)...)
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
	app.StartCmd()
	time.Sleep(1 * time.Second)
	app.ExecuteMsg(AlwaysTopCmd)
	for msg := range app.WebUI.MsgChan {
		if strings.HasPrefix(msg, "preset") {
			elements := strings.Split(msg, ":")
			//TODO process error
			id, _ := strconv.Atoi(elements[1])
			app.LoadPreset(id)
		} else if msg == StartCmd {
			app.StartCmd()
		} else if msg == StopCmd {
			app.StopCmd()
		} else {
			app.ExecuteMsg(msg)
		}
	}
	for _, cmd := range app.Cmds {
		cmd.Wait()
	}
}
