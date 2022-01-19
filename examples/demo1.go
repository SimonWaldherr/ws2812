package main

import (
	"flag"
	pixarray "github.com/Jon-Bright/ledctl/pixarray"
	cc "github.com/SimonWaldherr/ColorConverterGo"
	"log"
	"os"
	"time"
)

var lpd8806Dev = flag.String("dev", "/dev/spidev0.0", "The SPI device on which LPD8806 LEDs are connected")
var lpd8806SpiSpeed = flag.Uint("spispeed", 1000000, "The speed to send data via SPI to LPD8806s, in Hz")
var ws281xFreq = flag.Uint("ws281xfreq", 800000, "The frequency to send data to WS2801x devices, in Hz")
var ws281xDma = flag.Int("ws281xdma", 10, "The DMA channel to use for sending data to WS281x devices")
var ws281xPin0 = flag.Int("ws281xpin0", 18, "The pin on which channel 0 should be output for WS281x devices")
var ws281xPin1 = flag.Int("ws281xpin1", 13, "The pin on which channel 1 should be output for WS281x devices")
var ledChip = flag.String("ledchip", "ws281x", "The type of LED strip to drive: one of ws281x, lpd8806")
var port = flag.Int("port", 24601, "The port that the server should listen to")
var pixels = flag.Int("pixels", 5*32, "The number of pixels to be controlled")
var pixelOrder = flag.String("order", "GRB", "The color ordering of the pixels")

func main() {
	flag.Parse()
	order := pixarray.StringOrders[*pixelOrder]
	var leds pixarray.LEDStrip
	var err error
	switch *ledChip {
	case "lpd8806":
		dev, err := os.OpenFile(*lpd8806Dev, os.O_RDWR, os.ModePerm)
		if err != nil {
			log.Fatalf("Failed opening SPI: %v", err)
		}
		leds, err = pixarray.NewLPD8806(dev, *pixels, 3, uint32(*lpd8806SpiSpeed), order)
		if err != nil {
			log.Fatalf("Failed creating LPD8806: %v", err)
		}
	case "ws281x":
		leds, err = pixarray.NewWS281x(*pixels, 3, order, *ws281xFreq, *ws281xDma, []int{*ws281xPin0, *ws281xPin1})
		if err != nil {
			log.Fatalf("Failed creating WS281x: %v", err)
		}
	default:
		log.Fatalf("Unrecognized LED type: %v", *ledChip)
	}
	pa := pixarray.NewPixArray(*pixels, 3, leds)

	var p pixarray.Pixel

	p.R = 0
	p.G = 0
	p.B = 0
	pa.SetAll(p)
	pa.Write()
	time.Sleep(1500 * time.Millisecond)

	for {
		for j := 0; j < 360; j++ {
			for i := 0; i < pa.NumPixels(); i++ {
				h := int(360.0 / float32(pa.NumPixels()) * float32(i))
				p.R, p.G, p.B = cc.HSV2RGB(h+j, 100, 100)
				pa.SetOne((i)%pa.NumPixels(), p)
			}
			time.Sleep(2 * time.Millisecond)
			pa.Write()
		}
	}
}
