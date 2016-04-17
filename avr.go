package avr

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strconv"
)

// AVR contains the connection to the device.
type AVR struct {
	Conn net.Conn
}

// Start opens a goroutine to check for new data and decode it.
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

// Send sends a command with data to the device.
func (avr AVR) Send(cmd string, args ...interface{}) {
	fmt.Fprintf(avr.Conn, cmd+"\r\n", args...)
}

// decode attempts to determine what kind of data has been sent.
func (avr AVR) decode(line string) interface{} {
	if len(line) >= 2 && line[:2] == "FL" {
		return DisplayInfo{
			Message: avr.decodeHexString(line[2:]),
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

// decodeHexString decodes the hex string from the display info.
func (avr AVR) decodeHexString(str string) string {
	b := bytes.Buffer{}

	for i := 0; i < len(str)-1; i += 2 {
		decoded, _ := hex.DecodeString(str[i : i+2])
		b.Write(decoded)
	}

	return b.String()
}

// NewAVR connects to the device at IP:port.
func NewAVR(host string) (*AVR, error) {
	conn, err := net.Dial("tcp", host)
	if err != nil {
		return nil, err
	}

	return &AVR{
		Conn: conn,
	}, nil
}
