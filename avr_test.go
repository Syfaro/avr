package avr

import (
	"log"
	"testing"
)

func TestDecodeFL(t *testing.T) {
	testString := "FL024465657020536561204465746563"

	avr := &AVR{}
	result := avr.decode(testString)

	info, ok := result.(DisplayInfo)
	if !ok {
		t.Fail()
	}

	if info.Message != "Deep Sea Detec" {
		t.Fail()
	}
}

func ExampleNewAVR() {
	avr, err := NewAVR("192.168.1.224:23")
	if err != nil {
		log.Fatal(err)
	}

	for line := range avr.Start() {
		switch v := line.(type) {
		case DisplayInfo:
			log.Printf("Got display update: %s\n", v.Message)

			avr.Send(PowerOff)
		}
	}
}
