package main

import (
	"math"
	"math/rand"
	"time"

	"honnef.co/go/js/dom/v2"
)

const mag uint = 8
const len uint = 1 << mag
const scale uint = 3

func main() {
	var ctx = createCanvas(len, scale)

	var p [9]float64
	for i := range p {
		var beta = math.Log(1+math.Sqrt(2.0)) / 2
		p[i] = math.Exp(-2 * beta * float64(i-4))
	}

	var state = make([]int8, len*len)
	for i := range state {
		state[i] = int8(2*(i%2) - 1)
	}

	var start = time.Now()
	var count = 0
	for {
		for i := uint(0); i < len*len; i++ {
			var r = uint(rand.Int())
			var x, y = r % len, (r / len) % len
			var center = x + y*len
			var right = (x+1)%len + y*len
			var left = (x-1)%len + y*len
			var bottom = x + ((y+1)%len)*len
			var top = x + ((y-1)%len)*len
			var sum = state[right] + state[left] + state[bottom] + state[top]
			var delta = state[center] * sum
			if delta <= 0 || rand.Float64() < p[delta+4] {
				state[center] *= -1
			}
		}
		count += 1
		if t := time.Since(start); t > 300*time.Millisecond {
			// fmt.Printf("updates: %d; rate: %f/s\n", count, float64(count)/t.Seconds())
			start = time.Now()
			count = 0
			draw(len, scale, state, ctx)
			time.Sleep(2 * time.Millisecond)
		}
	}

}

func createCanvas(len, scale uint) *dom.CanvasRenderingContext2D {
	var document = dom.GetWindow().Document()
	var body = document.(dom.HTMLDocument).Body()
	var canvas = document.CreateElement("canvas").(*dom.HTMLCanvasElement)
	canvas.SetHeight(int(len * scale))
	canvas.SetWidth(int(len * scale))
	body.AppendChild(canvas)
	var ctx = canvas.GetContext2d()
	ctx.SetFillStyle("blue")
	return ctx
}

func draw(len, scale uint, state []int8, ctx *dom.CanvasRenderingContext2D) {
	for x := uint(0); x < len; x++ {
		for y := uint(0); y < len; y++ {
			if state[x+y*len] < 0 {
				ctx.ClearRect(float64(x*scale), float64(y*scale), float64(scale), float64(scale))
			} else {
				ctx.FillRect(float64(x*scale), float64(y*scale), float64(scale), float64(scale))
			}
		}
	}
}
