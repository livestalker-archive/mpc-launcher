package main

import (
	"sync"
	"os/exec"
	"fmt"
	"net/url"
	"net/http"
	"time"
	"log"
	"strings"
	"strconv"
)

// Application structure
type App struct {
	// Application config
	Config  *Config
	// Video presets
	Presets Presets
	// User web interface
	WebUI   *WebUI
	// External commands for video player
	Cmds    []*exec.Cmd
	// Timer for light scenes
	Timer   *time.Timer
}

// Init application
func (app *App) Init(configFilename string, presetsFilename string) {
	var err error
	msgChan := make(chan string)
	app.Config, err= LoadConfig(configFilename)
	if err != nil {
		log.Fatal("Can not load config file: ", err)
	}
	app.Presets, err = LoadPresets(presetsFilename)
	if err != nil {
		log.Fatal("Can not load presets file: ", err)
	}
	app.WebUI = NewWebUI(app.Config, app.Presets, msgChan)
}

// Start application
func (app *App) StartApp() {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go app.StartSupervisor(wg)
	wg.Add(1)
	go app.WebUI.StartServer(wg)
	wg.Wait()
}

func (app *App) StartSupervisor(wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range app.WebUI.MsgChan {
		if strings.HasPrefix(msg, "preset") {
			parts := strings.Split(msg, ":")
			if len(parts) != 2 {
				// wrong command can not load preset
				continue
			}
			id, err := strconv.Atoi(parts[1])
			if err != nil {
				// can not convert string to preset number
				continue
			}
			app.LoadPreset(id)
		} else {
			app.SendCommand(msg)
		}
	}
	for _, cmd := range app.Cmds {
		cmd.Wait()
	}
}

func (app *App) SendCommand(msg string) {
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
	if len(app.Cmds) > 0 {
		for i := 0; i < app.Config.MonCount; i++ {
			if app.Cmds[i] != nil && app.Cmds[i].Process != nil {
				app.Cmds[i].Process.Kill()
			}
		}
	}
	//p := NewPacket(0, 6, SceneOn)
	//p.Do()
	app.Cmds = make([]*exec.Cmd, app.Config.MonCount)
	for i := 0; i < app.Config.MonCount; i++ {
		app.Cmds[i] = exec.Command(app.Config.MpcPath, preset.Files[i].GetFullArgs()...)
		app.Cmds[i].Start()
	}
	if app.Timer != nil {
		app.Timer.Stop()
	}
	app.Timer = time.NewTimer(time.Duration(preset.Light.Time) * time.Second)
	go func() {
		<-app.Timer.C
		//p := NewPacket(0, preset.Light.Number, SceneOn)
		//p.Do()
		fmt.Println("Timer 2 expired")
	}()
}

func main() {
	configFilename := "./conf/config.yml"
	presetsFilename := "./conf/presets.yml"
	app := App{}
	app.Init(configFilename, presetsFilename)
	app.StartApp()
}

