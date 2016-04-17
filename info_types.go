package avr

import (
	"math"
)

// DisplayInfo is what information is currently on the display.
type DisplayInfo struct {
	Message string
}

// PowerInfo is the device's current power state.
type PowerInfo struct {
	On bool
}

// VolumeInfo is the device's current volume.
type VolumeInfo struct {
	Volume float32
}

// Convert turns the percentage value into a value the device understands.
func (volume VolumeInfo) Convert() int {
	return int(math.Floor(float64(volume.Volume*185/100) + .5))
}

// MuteInfo is the device's current mute state.
type MuteInfo struct {
	Muted bool
}
