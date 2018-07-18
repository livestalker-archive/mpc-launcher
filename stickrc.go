package main

import (
	"bytes"
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

const (
	SceneOff = 0
	SceneOn  = 1
)

func NewPacket(pn, sn, cmd int) *Packet {
	p := Packet{}
	copy(p.ID[:], "Stick_3A")
	p.OpCode[1] = 0
	p.OpCode[0] = 109
	p.SceneNr[0] = byte(pn*50 + sn)
	p.Command[0] = byte(cmd)
	return &p
}

func (p *Packet) Join() []byte {
	res := make([][]byte, 0)
	res = append(res, p.ID[:], p.OpCode[:], p.SceneNr[:], p.ZoneSyncID[:], p.Command[:])
	res = append(res, p.DimmerVal[:], p.SpeedVal[:], p.Unused[:], p.ColorVal[:])
	return bytes.Join(res, nil)
}

func (p *Packet) Do() {
	conn, err := net.Dial("udp", "192.168.16.139:2430")
	defer conn.Close()
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	_, err = conn.Write(p.Join())
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
}
/*
func main() {
	p := NewPacket(0, 4, 9)
	fmt.Println(hex.Dump(p.Join()))
	conn, err := net.Dial("udp", "192.168.16.139:2430")
	defer conn.Close()
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	n, err := conn.Write(p.Join())
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	fmt.Println(n)
	//message, _ := bufio.NewReader(conn).ReadString('\n')
	//fmt.Print("Message from server: " + message)
}
*/
