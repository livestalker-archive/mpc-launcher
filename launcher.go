package main

import (
	"sync"
	"os/exec"
	"fmt"
)

type App struct {
	Config  *Config
	Presets *Presets
	Cmds []*exec.Cmd
}

func (app *App) Init(configFilename string, presetsFilename string) {
	app.Config = LoadConfig(configFilename)
	app.Presets = LoadPresets(presetsFilename)
	app.Cmds = make([]*exec.Cmd, app.Config.MonCount)
	for i := 0; i< app.Config.MonCount; i++ {
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
	for _, cmd := range app.Cmds {
		cmd.Wait()
	}
}
