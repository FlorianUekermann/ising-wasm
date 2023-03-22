package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
	"time"

	"honnef.co/go/js/dom/v2"
)

func main() {
	// Magnitude of side length
	var mag int = 8
	// Side length (e.g. 256 for mag=8)
	var length int = 1 << mag

	// Precalculate probabilities for delta in range [-4, ..., +4]
	var p [9]float64
	for i := range p {
		var beta = math.Log(1+math.Sqrt(2.0)) / 2 // 0.44068679350977147
		p[i] = math.Exp(-2 * beta * float64(i-4))
	}

	var ctx, pixels = createCanvas(length)

	var state = make([]int8, length*length)
	for i := range state {
		state[i] = int8(2*(i%2) - 1)
	}
	var rng = rand.New(rand.NewSource(0))

	var start = time.Now()
	var lastDraw = time.Now()
	for sweeps := 0; true; sweeps++ {
		for off := 0; off < 3; off++ {
			for i := off; i < length*length; i += 3 {
				var y = i >> mag
				var x = i - y
				var center, left, right, top, bottom = xy2i(x, y, mag), xy2i(x-1, y, mag), xy2i(x+1, y, mag), xy2i(x, y-1, mag), xy2i(x, y+1, mag)
				var delta = state[center] * (state[left] + state[right] + state[top] + state[bottom])
				if delta <= 0 || rng.Float64() < p[delta+4] {
					state[center] *= -1
				}
			}
		}
		if time.Since(lastDraw) > 100*time.Millisecond {
			lastDraw = time.Now()
			fmt.Printf("sweep rate: %f/s\n", float64(sweeps)/time.Since(start).Seconds())
			draw(length, state, ctx, pixels)
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func xy2i(x int, y int, mag int) int {
	var mask = 1<<mag - 1
	return x&mask + (y&mask)<<mag
}

func createCanvas(len int) (*dom.CanvasRenderingContext2D, []byte) {
	var document = dom.GetWindow().Document()
	var body = document.(dom.HTMLDocument).Body()
	var canvas = document.CreateElement("canvas").(*dom.HTMLCanvasElement)
	canvas.Style().SetProperty("image-rendering", "pixelated", "important")
	canvas.SetHeight(len)
	canvas.SetWidth(len)
	body.AppendChild(canvas)
	var ctx = canvas.GetContext2d()
	ctx.SetFillStyle("blue")
	var pixels = make([]byte, len*len*4)
	return ctx, pixels
}

func draw(len int, state []int8, ctx *dom.CanvasRenderingContext2D, pixels []byte) {
	for i := 0; i < len*len; i++ {
		var pixel = pixels[i*4 : i*4+4]
		pixel[3] = 255 // alpha
		if state[i] < 0 {
			pixel[0] = 255 // red
			pixel[1] = 0
			pixel[2] = 0
		} else {
			pixel[0] = 0
			pixel[1] = 0
			pixel[2] = 255 // blue
		}
	}
	var arr = js.Global().Get("Uint8ClampedArray").New(len * len * 4)
	js.CopyBytesToJS(arr, pixels)
	ctx.PutImageData(&dom.ImageData{js.Global().Get("ImageData").New(arr, len)}, 0, 0)
}
