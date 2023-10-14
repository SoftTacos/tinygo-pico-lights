package main

import (
	"machine"
	"math"
	"time"

	"tinygo.org/x/drivers/ws2812"
)

func main() {
	numPixels := 24
	wait := time.Millisecond * 50
	machine.GP0.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ring1 := ws2812.New(machine.GP0, ws2812.SK6812)

	// clear any colors that might still be set if the device just restarted
	clear(ring1, numPixels)

	start := time.Now()
	t := swirling{
		numPixels: numPixels,
		// partitionCount: 2,
		period: time.Millisecond * 65,
		colors: [][]byte{
			// []byte{100, 0, 255, 0},
			// []byte{0, 255, 100, 0},
			// []byte{50, 255, 0, 0},
			// []byte{0, 0, 0, 0},

			// []byte{100, 0, 255, 0},
			[]byte{0, 0, 255, 255},
			[]byte{0, 0, 255, 255},
		},
		colorSets: [][][]byte{
			{
				[]byte{0, 255, 0, 255},
				[]byte{255, 0, 0, 255},
				[]byte{0, 0, 255, 255},
			},
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
		direction: -1,
	}

	// ring1.WriteColors([]color.RGBA{
	// 	{R: 100},
	// 	{G: 100},
	// 	{B: 100},
	// 	{A: 255},
	// })
	s := Spread{
		period: time.Millisecond * 250,
		colors: [][]byte{
			[]byte{0, 0, 255, 255},
			[]byte{0, 0, 255, 0},
		},
		numPixels: 24,
	}
	// time.Sleep(time.Second * 10)
	for {
		t.Draw(time.Since(start), ring1)
		// s.Draw(time.Since(start), ring1)
		_ = t
		_ = s
		time.Sleep(wait)
	}
}

func add(c1, c2 []byte) (c3 []byte) {
	c3 = make([]byte, 4)
	c3[0] = c1[0] + c2[0]
	c3[1] = c1[1] + c2[1]
	c3[2] = c1[2] + c2[2]
	c3[3] = c1[3] + c2[3]
	return
}

type Drawable interface {
	Draw(t time.Duration, device ws2812.Device)
}

func clear(device ws2812.Device, numPixels int) {
	device.Write(make([]byte, numPixels*4))
	time.Sleep(time.Millisecond * 50)
}

func writeBlinkErr(d ws2812.Device, b []byte) {
	_, err := d.Write(b)
	if err != nil {
		blinkLED()
	}
}

func blinkLED() {
	led := machine.LED
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	wait := time.Millisecond * 500
	for {
		led.Low()
		time.Sleep(wait)

		led.High()
		time.Sleep(wait)
	}
}

// type Drawable interface {
// 	Draw(elapsed uint64, driver ws2812.Device)
// }

type swirling struct {
	colors    [][]byte
	colorSets [][][]byte
	numPixels int // number of LEDs in the whole set
	// partitionCount int           // number of sub sections of the total LEDs
	period    time.Duration // time between the movement of LEDs
	direction int           // [ -1 | 1 ]
}

func (d swirling) Draw(t time.Duration, driver ws2812.Device) {
	// colors := d.colors
	colors := cycleTime[[][][]byte](t, d.period*50, d.colorSets)
	pixels := make([]byte, 0, d.numPixels*4)
	for _, color := range colors {
		pixels = append(pixels, lightFade(d.numPixels/len(colors), color)...)
	}
	pixels = rotateTime(t, d.period, pixels, d.direction)
	writeBlinkErr(driver, pixels)
}

type Spread struct {
	colors    [][]byte
	numPixels int           // number of LEDs in the whole set
	period    time.Duration // time between the movement of LEDs
}

func (s Spread) Draw(t time.Duration, driver ws2812.Device) {
	// pixels := spread(t, s.period, s.numPixels, s.color)
	pixels := make([]byte, 0, len(s.colors)*4)
	for _, color := range s.colors {
		pixels = append(pixels, spread(t, s.period, s.numPixels/len(s.colors), color)...)
	}
	writeBlinkErr(driver, pixels)
}

func spread(t, period time.Duration, numPixels int, color []byte) (pixels []byte) {
	radius := int(math.Round((math.Sin(float64(t)/float64(period)) + 1.0) * float64(numPixels) / 4))
	center := numPixels / 2
	pixels = make([]byte, (center-radius)*4)
	for i := 0; i < radius*2; i++ {
		pixels = append(pixels, color...)
	}
	pixels = append(pixels, make([]byte, (center-radius)*4)...)
	return
}

func repeat(pixels []byte, n int) []byte {
	for i := 1; i < n; i++ {
		pixels = append(pixels, pixels...)
	}
	return pixels
}

func lightFade(size int, color []byte) (pixels []byte) {
	pixels = make([]byte, 0, size)
	for i := 0; i < size; i++ {
		color = color
		for c := range color {
			pixels = append(pixels, uint8(float64(color[c])*(float64(i)/float64(size)))) // byte(float64(d.min) + float64(d.max-d.min)*(float64(len(partition))/float64(i+1)))
		}
	}

	return
}

// rotates 1 pixel per time period expired
func rotateTime(t, period time.Duration, pixels []byte, direction int) (p []byte) {
	shift := int(t/period) * direction
	p = make([]byte, len(pixels))
	for i := range p {
		s := (i + shift*4) % len(pixels)
		if s < 0 {
			s += len(pixels)
		}
		p[i] = pixels[s]
	}
	return
}

// func colorTimeShift(t, period time.Duration, shift, pixels []byte) []byte {
// 	for i := 0; i < len(pixels)/4; i += 4 {

// 	}
// 	return pixels
// }

func cycleTime[T []B, B any](t, period time.Duration, s T) B {
	return s[int(t/period)%len(s)]
}
