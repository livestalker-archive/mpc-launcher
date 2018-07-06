package main

import (
	"path/filepath"
	"sync"
	"os/exec"
	"net/http"
	"fmt"
	"context"
	"net/url"
)

type Worker struct {
	Path string
	Args []string
	Port string
}

type Controller struct {
	Srv *http.Server
	Cmd chan string
}

func (ctrl *Controller) mainPage(writer http.ResponseWriter, request *http.Request) {
	content := `<html>
<head>
<title>Luncher</title>
</head>
<a href='/play'>PLAY</a> | <a href='/pause'>PAUSE</a><br/>
</html>`
	fmt.Fprintln(writer, content)
}

func (ctrl *Controller) shutdownButton(writer http.ResponseWriter, request *http.Request) {
	close(ctrl.Cmd)
	if err := ctrl.Srv.Shutdown(context.Background()); err != nil {
		panic(err)
	}
	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}

func (ctrl *Controller) playButton(writer http.ResponseWriter, request *http.Request) {
	ctrl.Cmd <- "play"
	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}

func (ctrl *Controller) pauseButton(writer http.ResponseWriter, request *http.Request) {
	ctrl.Cmd <- "pause"
	http.Redirect(writer, request, "/", http.StatusTemporaryRedirect)
}

func initCommands() []Worker {
	workers := make([]Worker, 2)
	programPath := filepath.FromSlash("C:\\Program Files\\MPC-HC\\mpc-hc64.exe")
	workers[0] = Worker{
		Path: programPath,
		Args: []string{"G:\\tmp\\1.mp4", "/webport", "8881", "/open", "/monitor", "1"},
		Port: "8881",
	}
	workers[1] = Worker{
		Path: programPath,
		Args: []string{"G:\\tmp\\1.mp4", "/webport", "8882", "/open", "/monitor", "2"},
		Port: "8882",
	}
	return workers
}

func main() {
	// wm_command
	// 887 - play
	// 888 - pause
	// address := "http://localhost:8888/command.html"
	wg := &sync.WaitGroup{}
	workers := initCommands()
	cmdChan := make(chan string)
	wg.Add(1)
	go supervisor(wg, workers, cmdChan)
	wg.Add(1)
	go httpServer(wg, cmdChan)
	wg.Wait()
	fmt.Println("Application closed")
}

func supervisor(wg *sync.WaitGroup, workers []Worker, cmdChan chan string) {
	defer wg.Done()
	cmds := make([]*exec.Cmd, len(workers))
	for ix, w := range workers {
		cmd := exec.Command(w.Path, w.Args...)
		cmd.Start()
		cmds[ix] = cmd
	}
	for msg := range cmdChan {
		switch msg {
		case "play":
			postPlay()
		case "pause":
			postPause()
		}
	}
	for _, cmd := range cmds {
		cmd.Wait()
	}
}
func postPause() {
	data := url.Values{}
	data.Add("wm_command", "888")
	http.PostForm("http://localhost:8881/command.html", data)
	http.PostForm("http://localhost:8882/command.html", data)
	fmt.Println("Send pause command")
}
func postPlay() {
	data := url.Values{}
	data.Add("wm_command", "887")
	http.PostForm("http://localhost:8881/command.html", data)
	http.PostForm("http://localhost:8882/command.html", data)
	fmt.Println("Send play command")
}

func httpServer(wg *sync.WaitGroup, cmdChan chan string) {
	defer wg.Done()
	srv := &http.Server{Addr: ":7777"}
	controller := Controller{Srv: srv, Cmd: cmdChan}
	http.HandleFunc("/", controller.mainPage)
	http.HandleFunc("/shutdown", controller.shutdownButton)
	http.HandleFunc("/play", controller.playButton)
	http.HandleFunc("/pause", controller.pauseButton)
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
