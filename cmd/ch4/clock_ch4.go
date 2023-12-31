package main

import (
	"math"
	. "roytracer/gfx"
	. "roytracer/math"
)

func main() {
	canvas := NewCanvas(300, 300, ColorBlack)
	drawClock(canvas)
	w := PPMWriter{MaxLineLength: 70}
	w.Write(canvas)
	w.SaveFile("scene.ppm")
}

func drawClock(c *Canvas) {
	center := Vec3{float64(c.Width)/2, float64(c.Height)/2, 0}
	radius := 120.0
	angle := 2.0*math.Pi/12.0
	arm := Vec4{0, -radius, 0, 1}

	for i := 0; i < 12; i++ {
		m, n := project(Trans(center[0], center[1], center[2]).MulVec(arm), c)
		c.WritePixel(n, m, ColorWhite)
		arm = RotZ(angle).MulVec(arm)
	}
}

func project(point Vec4, c *Canvas) (uint32, uint32) {
	return uint32(point[1]), uint32(point[0])
}
