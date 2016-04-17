package avr

const (
	// DisplayInfoQuery allows you to get the current display text.
	DisplayInfoQuery = "?FL"
)

const (
	// PowerOn turns on the device.
	PowerOn = "PO"
	// PowerOff turns off the device.
	PowerOff = "PF"

	// PowerQuery allows you to get the current power state.
	PowerQuery = "?P"
)

const (
	// VolumeUp turns up the volume.
	VolumeUp = "VU"
	// VolumeDown turns down the volume.
	VolumeDown = "VD"
	// VolumeSet sets the volume to a specific number.
	VolumeSet = "%03dVL"

	// VolumeQuery allows you to get the current volume setting.
	VolumeQuery = "?V"
)

const (
	// MuteOn mutes the device.
	MuteOn = "MO"
	// MuteOff unmutes the device.
	MuteOff = "MF"

	// MuteQuery allows you to get the current mute setting.
	MuteQuery = "?M"
)
