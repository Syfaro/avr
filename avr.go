package avr

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math"
	"net"
	"strconv"
)

const (
	DisplayInfoQuery = "?FL"
)

const (
	PowerOn  = "PO"
	PowerOff = "PF"

	PowerQuery = "?P"
)

const (
	VolumeUp   = "VU"
	VolumeDown = "VD"
	VolumeSet  = "%dVL"

	VolumeQuery = "?V"
)

const (
	MuteOn  = "MO"
	MuteOff = "MF"

	MuteQuery = "?M"
)

type DisplayInfo struct {
	GUIDIcon      bool
	VolumeDisplay bool
	Message       string
}

type PowerInfo struct {
	On bool
}

type VolumeInfo struct {
	Volume float32
}

func (volume VolumeInfo) Convert() int {
	return int(math.Floor(float64(volume.Volume*185/100) + .5))
}

type MuteInfo struct {
	Muted bool
}

type AVR struct {
	Conn net.Conn
}

func (avr AVR) Start() <-chan interface{} {
	out := make(chan interface{})

	go func() {
		for {
			resp, err := bufio.NewReader(avr.Conn).ReadString('\n')
			if err != nil {
				close(out)
				return
			}

			// remove \r\n at end of each message
			out <- avr.decode(resp[:len(resp)-2])
		}
	}()

	return out
}

func (avr AVR) Send(cmd string, args ...interface{}) {
	fmt.Fprintf(avr.Conn, cmd+"\r\n", args...)
}

func (avr AVR) decode(line string) interface{} {
	if len(line) >= 2 && line[:2] == "FL" {
		status := byte(line[3])

		return DisplayInfo{
			Message:       avr.decodeHexString(line[2:]),
			GUIDIcon:      status&(1<<0) == (1 << 0),
			VolumeDisplay: status&(1<<1) == (1 << 1),
		}
	} else if len(line) >= 3 && line[:3] == "PWR" {
		return PowerInfo{
			On: line[3:4] == "0",
		}
	} else if len(line) >= 3 && line[:3] == "VOL" {
		vol, err := strconv.ParseInt(line[3:], 10, 32)
		if err != nil {
			log.Fatal(err)
		}

		return VolumeInfo{
			Volume: (float32(vol) * 100) / 185,
		}
	} else if len(line) >= 3 && line[:3] == "MUT" {
		return MuteInfo{
			Muted: line[3:4] == "0",
		}
	}

	return nil
}

func (avr AVR) decodeHexString(str string) string {
	b := bytes.Buffer{}

	for i := 0; i < len(str)-1; i += 2 {
		decoded, _ := hex.DecodeString(str[i : i+2])
		b.Write(decoded)
	}

	return b.String()
}

func NewAVR(host string) (*AVR, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return &AVR{
		Conn: conn,
	}, nil
}
