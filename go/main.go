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
	var len int = 1 << mag
	var lenMask = len - 1
	// Total number of sites
	var size int = len * len
	var sizeMask = size - 1

	// Precalculate probabilities for delta in range [-4, ..., +4]
	var p [9]float64
	for i := range p {
		var beta = math.Log(1+math.Sqrt(2.0)) / 2 // 0.44068679350977147
		p[i] = math.Exp(-2 * beta * float64(i-4))
	}

	var ctx = createCanvas(len)

	var state = make([]int8, size)
	for i := range state {
		state[i] = 1
	}

	var rng = rand.New(rand.NewSource(0))
	var start = time.Now()
	var lastDraw = time.Now()
	var sweeps = 0
	for {
		for off := 0; off < 3; off++ {
			for center := off; center < size; center += 3 {
				var col = center & lenMask
				var row_off = center - col
				var down, up = (center + len) & sizeMask, (center + size - len) & sizeMask
				var right, left = row_off + (col+1)&lenMask, row_off + (col+len-1)&lenMask
				var delta = state[center] * (state[left] + state[right] + state[up] + state[down])
				if delta <= 0 || rng.Float64() < p[delta+4] {
					state[center] *= -1
				}
			}
		}
		sweeps++
		if time.Since(lastDraw) > 100*time.Millisecond {
			lastDraw = time.Now()
			fmt.Printf("sweep rate: %f/s\n", float64(sweeps)/time.Since(start).Seconds())
			draw(len, state, ctx)
			// Give the browser a chance to present
			time.Sleep(1 * time.Millisecond)
		}
	}
}

func createCanvas(len int) *dom.CanvasRenderingContext2D {
	var document = dom.GetWindow().Document()
	var body = document.(dom.HTMLDocument).Body()
	var canvas = document.CreateElement("canvas").(*dom.HTMLCanvasElement)
	canvas.Style().SetProperty("image-rendering", "pixelated", "")
	canvas.SetHeight(len)
	canvas.SetWidth(len)
	body.AppendChild(canvas)
	return canvas.GetContext2d()
}

func draw(len int, state []int8, ctx *dom.CanvasRenderingContext2D) {
	var pixels = make([]byte, len * len * 4)
	for i := 0; i < len*len; i++ {
		var pixel = pixels[i*4 : i*4+4]
		pixel[3] = 255 // alpha
		if state[i] < 0 {
			pixel[0] = 255 // red
		} else {
			pixel[2] = 255 // blue
		}
	}
	var arr = js.Global().Get("Uint8ClampedArray").New(len * len * 4)
	js.CopyBytesToJS(arr, pixels)
	ctx.PutImageData(&dom.ImageData{js.Global().Get("ImageData").New(arr, len)}, 0, 0)
}
