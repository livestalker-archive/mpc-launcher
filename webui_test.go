package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"fmt"
)

func TestWebUI_MainPage(t *testing.T) {
	ui := WebUI{}
	url := "http://localhost"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	ui.MainPage(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("Wrong StatusCode: got %d, expected %d", w.Code, http.StatusOK)
	}
}

func TestWebUI_LoadPreset(t *testing.T) {
	ui := WebUI{}
	ui.MsgChan = make(chan string)
	url := "http://localhost/load?preset=0"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	ui.LoadPreset(w, req)
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Wrong StatusCode: got %d, expected %d", w.Code, http.StatusTemporaryRedirect)
	}
	cmd := <-ui.MsgChan
	if cmd != fmt.Sprintf(PresetCmd, "0") {
		t.Errorf("Get wrong command: %s", cmd)
	}
}

func TestWebUI_PlayButton(t *testing.T) {
	ui := WebUI{}
	ui.MsgChan = make(chan string)
	url := "http://localhost/play"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	ui.PlayButton(w, req)
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Wrong StatusCode: got %d, expected %d", w.Code, http.StatusTemporaryRedirect)
	}
	cmd := <-ui.MsgChan
	if cmd != PlayCmd {
		t.Errorf("Get wrong command: %s", cmd)
	}
}

func TestWebUI_PauseButton(t *testing.T) {
	ui := WebUI{}
	ui.MsgChan = make(chan string)
	url := "http://localhost/pause"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	ui.PauseButton(w, req)
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Wrong StatusCode: got %d, expected %d", w.Code, http.StatusTemporaryRedirect)
	}
	cmd := <-ui.MsgChan
	if cmd != PauseCmd {
		t.Errorf("Get wrong command: %s", cmd)
	}
}

func TestWebUI_StopButton(t *testing.T) {
	ui := WebUI{}
	ui.MsgChan = make(chan string)
	url := "http://localhost/stop"
	req := httptest.NewRequest("GET", url, nil)
	w := httptest.NewRecorder()
	ui.StopButton(w, req)
	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Wrong StatusCode: got %d, expected %d", w.Code, http.StatusTemporaryRedirect)
	}
	cmd := <-ui.MsgChan
	if cmd != StopCmd {
		t.Errorf("Get wrong command: %s", cmd)
	}
}
