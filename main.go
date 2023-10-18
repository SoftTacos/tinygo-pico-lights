package main

import (
	"machine"
	"time"
	"tinygo-pico-lights/device"
	"tinygo-pico-lights/effects"
)

func main() {
	devices := setupDevices()
	effects := setupEffects()

	orchestrator := Orchestrator{
		refreshDuration: time.Millisecond * 50,
		effectDuration:  time.Second * 10,
	}

	orchestrator.Start(devices, effects)
}

type Orchestrator struct {
	effectDuration  time.Duration
	refreshDuration time.Duration
}

func (o Orchestrator) Start(devices []device.Device, effects []effects.Effect) {
	var (
		e     = 0
		start = time.Now()
		since time.Duration
	)
	for {
		since = time.Since(start)
		e = (int(since / o.effectDuration)) % len(effects)
		for _, ring := range devices {
			effects[e].Draw(since, ring)
		}
		time.Sleep(o.refreshDuration)
	}
}

func setupDevices() (devices []device.Device) {
	devices = []device.Device{
		device.NewDevice(machine.GP0, 24),
	}
	for _, ring := range devices {
		clear(ring)
	}
	return
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

func clear(device device.Device) {
	device.Write(make([]byte, device.NumPixels*4))
	time.Sleep(time.Millisecond * 50)
}
