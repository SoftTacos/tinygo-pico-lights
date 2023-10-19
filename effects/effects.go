package effects

import (
	"machine"
	"math"
	"time"
	"tinygo-pico-lights/device"
)

type Effect interface {
	Draw(t time.Duration, device device.Device)
}

type Direction int

const (
	Clockwise        Direction = -1
	CounterClockwise Direction = 1
)

func add(c1, c2 []byte) (c3 []byte) {
	c3 = make([]byte, 4)
	c3[0] = c1[0] + c2[0]
	c3[1] = c1[1] + c2[1]
	c3[2] = c1[2] + c2[2]
	c3[3] = c1[3] + c2[3]
	return
}

func writeBlinkErr(d device.Device, b []byte) {
	_, err := d.Write(b)
	if err != nil {
		BlinkLED()
	}
}

func BlinkLED() {
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

// func NewSwirl(numPixels int, period time.Duration, direction Direction, colorSets [][][]byte) Swirl {
// 	return Swirl{
// 		ColorSets: colorSets,
// 		NumPixels: numPixels,
// 		Period:    period,
// 		Direction: direction,
// 	}
// }

type Swirl struct {
	ColorSets   [][][]byte    // represents [set per period][different colors for a period][color]
	Period      time.Duration // time between the movement of pixels
	ColorPeriod time.Duration // time between changes to LED color
	Direction   Direction     // direction to rotate/shift the pixels [ -1 | 1 ]
}

func (s Swirl) Draw(t time.Duration, d device.Device) {
	colors := cycleTime[[][][]byte](t, s.ColorPeriod, s.ColorSets)
	pixels := make([]byte, 0, d.NumPixels*4)
	for _, color := range colors {
		pixels = append(pixels, lightFade(d.NumPixels/len(colors), color, s.Direction)...)
	}
	pixels = rotateTime(t, s.Period, pixels, s.Direction)
	writeBlinkErr(d, pixels)
}

type Spread struct {
	colors [][]byte
	period time.Duration // time between the movement of pixels
}

func (s Spread) Draw(t time.Duration, d device.Device) {
	// pixels := spread(t, s.period, s.numPixels, s.color)
	pixels := make([]byte, 0, d.NumPixels*4)
	for _, color := range s.colors {
		pixels = append(pixels, spread(t, s.period, d.NumPixels/len(s.colors), color)...)
	}
	writeBlinkErr(d, pixels)
}

// sets a subsection of the pixels to the colors defined in ColorSets
type Band struct {
	BandSize    int           //number of pixels to set
	ColorSets   [][][]byte    // represents [set per period][different colors for a period][color]
	Period      time.Duration // time between the movement of pixels
	ColorPeriod time.Duration // time between changes to LED color
	Direction   Direction     // direction to rotate/shift the pixels [ -1 | 1 ]
}

func (b Band) Draw(t time.Duration, d device.Device) {
	pixels := make([]byte, 0, d.NumPixels*4)
	colors := cycleTime[[][][]byte](t, b.ColorPeriod, b.ColorSets)
	for _, color := range colors {
		pixels = append(pixels, band(d.NumPixels/len(colors), b.BandSize, color)...)
	}
	pixels = rotateTime(t, b.Period, pixels, b.Direction)
	writeBlinkErr(d, pixels)
}

func band(numPixels, bandSize int, color []byte) (pixels []byte) {
	pixels = make([]byte, (numPixels-bandSize)*4)
	for i := 0; i < bandSize; i++ {
		pixels = append(pixels, color...)
	}
	return
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

func lightFade(size int, color []byte, direction Direction) (pixels []byte) {
	pixels = make([]byte, 0, size)
	var (
		s float32
		c byte
	)
	if direction == CounterClockwise {
		s = 1
	}
	for i := 0; i < size; i++ {
		color = color
		for _, c = range color {
			// if ccw -> *(1 - ratio), if cw -> *(0 + ratio)
			pixels = append(pixels, uint8(float32(c)*(s+(-1.0*float32(direction))*(float32(i)/float32(size))))) // byte(float64(d.min) + float64(d.max-d.min)*(float64(len(partition))/float64(i+1)))
		}
	}

	return
}

// rotates 1 pixel per time period expired
func rotateTime(t, period time.Duration, pixels []byte, direction Direction) (p []byte) {
	shift := int(t/period) * int(direction)
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

func cycleTime[T []B, B any](t, period time.Duration, s T) B {
	return s[int(t/period)%len(s)]
}

// type Customizable struct {
// 	numPixels int
// }

// func (c Customizable) Draw(t time.Duration, device ws2812.Device) {
// 	pixels := make([]byte, 0, c.numPixels*4)

// 	writeBlinkErr(device, pixels)
// }

// func shiftTime(t, period time.Duration, pixels []byte, color []byte) []byte {

// 	return pixels
// }

// func colorTimeShift(t, period time.Duration, shift, pixels []byte) []byte {
// 	for i := 0; i < len(pixels)/4; i += 4 {

// 	}
// 	return pixels
// }
