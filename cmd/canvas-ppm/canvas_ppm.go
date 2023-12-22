package main

import (
	"roytracer/gfx"
	m "roytracer/math"
)

func main() {
	canvas := gfx.NewCanvas(9, 7, m.Color4(1, 0, 0, 0))
	canvas.WritePixel(0, 0, m.Color4(0, 1, 0, 0))
	canvas.WritePixel(0, 6, m.Color4(0, 1, 0, 0))
	canvas.WritePixel(8, 0, m.Color4(0, 1, 0, 0))
	canvas.WritePixel(8, 6, m.Color4(0, 1, 0, 0))
	canvas.WritePixel(4, 3, m.Color4(0, 0, 1, 0))
	ppmWriter := gfx.PPMWriter{MaxLineLength: 70}
	ppmWriter.Write(canvas)
	ppmWriter.SaveFile("scene.ppm")
}
