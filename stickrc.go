// It is simple module implement network remote protocol
// for STICK-DE3 device. For more details see:
// https://www.nicolaudie.com/stick-de3.htm
package main

import (
	"encoding/binary"
	"net"
	"fmt"
)

type Packet struct {
	ID         [8]byte
	OpCode     [2]byte
	SceneNr    [2]byte
	ZoneSyncID [1]byte
	Command    [1]byte
	DimmerVal  [2]byte
	SpeedVal   [2]byte
	Unused     [2]byte
	ColorVal   [4]byte
}

// Package format. Fore more details see:
// https://storage.googleapis.com/nicolaudie-eu-litterature/Release/stick3_remote_protocol.pdf
type RCPacket struct {
	ID         []byte
	OpCode     uint16
	SceneNr    uint16
	ZoneSyncID uint8
	Command    uint8
	DimmerVal  uint16
	SpeedVal   uint16
	Unused1    uint8
	Unused2    uint8
	ColorVal   uint32
}

// Quick triggering commands
const (
	SceneOff       = iota
	SceneOn
	ScenePauseOff
	ScenePauseOn
	SceneReset
	SceneDimmerSet
	SceneSpeedSet
	SceneColorSet
	BlackOutOff
	BlackOutOn
)

const (
	StickUDPPort = "2430"
)

// Create new RCPacket
func NewRCPacket() *RCPacket {
	p := RCPacket{
		ID:     []byte("Stick_3A"),
		OpCode: 109,
	}
	return &p
}

// Set scene number. SceneNr = Page Number * 50 + Scene Number
// The maximum number of scenes per page is 50. If more then
// 50 scenes have been added to a page, a second page will be
// allocated even if it does not appear this way on the device display.
func (p *RCPacket) SetScene(pn, sn uint16) {
	p.SceneNr = pn*50 + sn
}

// Set command for RCPackage
func (p *RCPacket) SetCommand(cmd uint8) {
	p.Command = cmd
}

// Get slice of bytes from RCPackage
func (p *RCPacket) GetBytes() []byte {
	b := make([]byte, 24)
	copy(b[0:8], p.ID)
	binary.LittleEndian.PutUint16(b[8:10], p.OpCode)
	binary.LittleEndian.PutUint16(b[10:12], p.SceneNr)
	b[12] = p.ZoneSyncID
	b[13] = p.Command
	return b
}

// Send a UDP packet to 2430 port to trigger the STICK.
func (p *RCPacket) SendBytes(host string) (int, error) {
	address := fmt.Sprintf("%s:%s", host, StickUDPPort)
	conn, err := net.Dial("udp4", address)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	n, err := conn.Write(p.GetBytes())
	if err != nil {
		return 0, err
	}
	return n, nil
}
