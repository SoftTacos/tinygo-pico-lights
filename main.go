package main

import (
	"machine"
	"time"
	"tinygo-pico-lights/device"
	"tinygo-pico-lights/effects"
)

func main() {
	machine.GP0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	numPixels := 24
	devices := []device.Device{
		// device.New(machine.GP0, ws2812.SK6812, numPixels),
		// ws2812.New(machine.GP0), //, ws2812.SK6812, numPixels),
		device.NewDevice(machine.GP0, numPixels),
	}

	wait := time.Millisecond * 50

	// clear any colors that might still be set if the device just restarted
	for _, ring := range devices {
		clear(ring, numPixels)
	}
	effects := setupEffects()
	e := 0
	effectDuration := time.Second * 10
	start := time.Now()
	var since time.Duration
	for {
		since = time.Since(start)
		e = (int(since / effectDuration)) % len(effects)
		for _, ring := range devices {
			effects[e].Draw(since, ring)
		}
		time.Sleep(wait)
	}
}

func setupEffects() (eff []effects.Effect) {

	// s := effects.Spread{
	// 	period: time.Millisecond * 250,
	// 	colors: [][]byte{
	// 		[]byte{0, 0, 255, 255},
	// 		[]byte{0, 0, 255, 0},
	// 	},
	// 	numPixels: 24,
	// }
	eff = []effects.Effect{
		effects.Band{
			ColorSets: [][][]byte{
				{
					{255, 0, 255, 0},
					{255, 0, 255, 0},
				},
				{
					{0, 0, 255, 255},
					{0, 0, 255, 255},
				},
			},
			BandSize:    5,
			Period:      time.Millisecond * 100,
			Direction:   effects.Clockwise,
			ColorPeriod: time.Second * 3,
		},
		effects.Swirl{
			Period:      time.Millisecond * 65,
			ColorPeriod: time.Second * 3,
			ColorSets: [][][]byte{
				{
					[]byte{0, 100, 255, 0},
					[]byte{0, 100, 255, 0},
					[]byte{0, 100, 255, 0},
				},
				{
					[]byte{0, 100, 255, 0},
					[]byte{0, 100, 255, 0},
				},
				{
					// []byte{0, 0, 255, 255},
					[]byte{0, 100, 255, 0},
				},
				{
					[]byte{0, 0, 255, 255},
					[]byte{0, 0, 255, 0},
				},
				{
					[]byte{0, 0, 255, 255},
					[]byte{255, 0, 255, 0},
				},
				{
					[]byte{0, 0, 255, 0},
					[]byte{255, 0, 255, 0},
				},
				{
					[]byte{0, 0, 255, 255},
					[]byte{255, 0, 255, 0},
				},
			},
			Direction: effects.Clockwise,
		},
	}
	return
}

func clear(device device.Device, numPixels int) {
	device.Write(make([]byte, numPixels*4))
	time.Sleep(time.Millisecond * 50)
}
