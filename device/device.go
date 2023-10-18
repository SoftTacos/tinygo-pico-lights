package device

import (
	"machine"

	"tinygo.org/x/drivers/ws2812"
)

func NewDevice(pin machine.Pin, numPixels int) Device {
	pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	return Device{
		NumPixels: numPixels,
		Device:    ws2812.New(pin),
	}
}

// wrapper for ws2812.Device
// embeds ws2812.Device so we have easy access to all methods,
// constructor also sets up the new pins for easy configuration.
// Holds the NumPixels so effects don't need to
type Device struct {
	NumPixels int
	ws2812.Device
}
