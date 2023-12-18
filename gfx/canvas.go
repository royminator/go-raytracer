package gfx

import (
	"fmt"
	"math"

	m "roytracer/math"
)

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

const (
	PPMFormat   = "P3"
	PPMMaxColor = 255
	PPMMaxCharsPerLine = 70
)

var (
	ColorWhite m.Vec4 = m.Vec4{1, 1, 1, 0}
	ColorBlack m.Vec4 = m.Vec4{0, 0, 0, 0}
	ColorRed   m.Vec4 = m.Vec4{1, 0, 0, 0}
)

func NewCanvas(width, height uint32, color m.Vec4) *Canvas {
	pixels := make([][]m.Vec4, height)
	for i := range pixels {
		row := make([]m.Vec4, width)
		for n := range row {
			row[n] = color
		}
		pixels[i] = row
	}
	return &Canvas{
		Width:  width,
		Height: height,
		Pixels: pixels,
	}
}

func (c *Canvas) WritePixel(n, m uint32, color m.Vec4) {
	c.Pixels[m][n] = color
}

func (c *Canvas) CreatePPMHeader() string {
	return fmt.Sprintf("%s\n%d %d\n%d\n", PPMFormat, c.Width, c.Height, PPMMaxColor)
}

func (c *Canvas) CreatePPMPixelData() string {
	var pixelData string
	for m := uint32(0); m < c.Height; m++ {
		var line string
		for n := uint32(0); n < c.Width; n++ {
			colorString := colorToPPM(c.Pixels[m][n], PPMMaxColor)
			line += colorString
			if n != c.Width-1 {
				line += " "
			}
		}
		pixelData += line
		pixelData += "\n"
	}
	return pixelData
}

func (c *Canvas) ToPPM() string {
	ppm := c.CreatePPMHeader()
	ppm += c.CreatePPMPixelData()
	return ppm
}

func clamp(color float64, maxVal float64) float64 {
	if color > 1 {
		return maxVal
	}
	if color < 0 {
		return 0
	}
	return color*maxVal
}
func colorToPPM(color m.Vec4, maxVal int) string {
	return fmt.Sprintf("%d %d %d",
		int(math.Round(clamp(color[0], PPMMaxColor))),
		int(math.Round(clamp(color[1], PPMMaxColor))),
		int(math.Round(clamp(color[2], PPMMaxColor))),
	)
}
