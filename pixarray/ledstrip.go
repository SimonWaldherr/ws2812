package pixarray

import (
	rpi "github.com/SimonWaldherr/ws2812/rpi"
)

type LEDStrip interface {
	RPi() *rpi.RPi
	MaxPerChannel() int
	GetPixel(i int) Pixel
	SetPixel(i int, p Pixel)
	Write() error
}
