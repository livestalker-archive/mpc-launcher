package main

import (
	"testing"
	"bytes"
	"net"
	"fmt"
	"log"
	"sync"
)

func TestNewPacket(t *testing.T) {
	p := NewRCPacket()
	if string(p.ID) != "Stick_3A" {
		t.Error("Wrong ID:", string(p.ID))
	}
	if p.OpCode != 109 {
		t.Error("Wrong operation code. Value must be 109 but it:", p.OpCode)
	}
}

func TestRCPacket_SetScene(t *testing.T) {
	p := NewRCPacket()
	p.SetScene(0, 1)
	if p.SceneNr != 1 {
		t.Error("wrong scene number:", p.SceneNr)
	}
}

func TestRCPacket_SetCommand(t *testing.T) {
	p := NewRCPacket()
	p.SetCommand(SceneOff)
	if p.Command != 0 {
		t.Error("Wrong constant values")
	}
}

func TestRCPacket_GetBytes(t *testing.T) {
	b := []byte{0x53, 0x74, 0x69, 0x63, 0x6b, 0x5f, 0x33, 0x41, 0x6d, 0x00, 0x05, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	p := NewRCPacket()
	p.SetScene(0, 5)
	p.SetCommand(SceneOn)
	if bytes.Compare(b, p.GetBytes()) != 0 {
		t.Error("Wrong bytes slice created from RCPacket")
	}
}

func TestRCPacket_SendBytes(t *testing.T) {
	p := NewRCPacket()
	p.SetScene(0, 5)
	p.SetCommand(SceneOn)
	wg := sync.WaitGroup{}
	ch := make(chan struct{})
	udpServer := func() {
		defer wg.Done()
		address := fmt.Sprintf("%s:%s", "localhost", StickUDPPort)
		pc, err := net.ListenPacket("udp4", address)
		if err != nil {
			log.Fatal("Can not create local UDP server:", err)
		}
		close(ch)
		defer pc.Close()
		buf := make([]byte, 1024)
		n, _, err := pc.ReadFrom(buf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Server read %d bytes\n", n)
	}
	wg.Add(1)
	go udpServer()
	<-ch
	n, err := p.SendBytes("localhost")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("Client write %d bytes\n", n)
	wg.Wait()
}
