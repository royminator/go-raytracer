package gfx

import m "raytracer/math"

type (
	Canvas struct {
		Pixels [][]m.Vec4
		Width  uint32
		Height uint32
	}

	Pixel struct {
		Color m.Vec4
	}
)

var (
	ColorWhite m.Vec4 = m.Vec4{1, 1, 1, 0}
	ColorRed   m.Vec4 = m.Vec4{1, 0, 0, 0}
)

func NewCanvas(width, height uint32, color m.Vec4) Canvas {
	pixels := make([][]m.Vec4, width)
	for i := range pixels {
		col := make([]m.Vec4, height)
		for j := range col {
			col[j] = color
		}
		pixels[i] = col
	}
	return Canvas{
		Width:  width,
		Height: height,
		Pixels: pixels,
	}
}

func (c *Canvas) WritePixel(i, j uint32, color m.Vec4) {
	c.Pixels[i][j] = color
}
