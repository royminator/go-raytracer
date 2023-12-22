package main

import (
	"roytracer/gfx"
	m "roytracer/math"
)

func main() {
	canvas := gfx.NewCanvas(80, 60, m.Color4(1, 0, 0, 0))
	ppmWriter := gfx.PPMWriter{}
	ppmWriter.Write(canvas)
	ppmWriter.SaveFile("scene.ppm")
}
