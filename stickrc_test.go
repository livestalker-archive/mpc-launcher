package main

import (
	"testing"
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