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

	PPMWriter struct {
		Header    string
		PixelData string
		Ppm       string
		buf       string
		m, n uint32
	}
)

const (
	PPMFormat          = "P3"
	PPMMaxColor        = 255
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

func clamp(color float64, maxVal float64) float64 {
	if color > 1 {
		return maxVal
	}
	if color < 0 {
		return 0
	}
	return color * maxVal
}

func floatToPpm(color float64, maxVal int) string {
	return fmt.Sprintf("%d",
		int(math.Round(clamp(color, PPMMaxColor))),
	)
}

func (w *PPMWriter) Write(c *Canvas) {
	w.WriteHeader(c)
	w.WritePixelData(c)
	w.Ppm = w.Header
	w.Ppm += "\n"
	w.Ppm += w.PixelData
}

func (w *PPMWriter) WriteHeader(c *Canvas) {
	w.Header = fmt.Sprintf("%s\n%d %d\n%d",
		PPMFormat, c.Width, c.Height, PPMMaxColor)
}

func (w *PPMWriter) WritePixelData(c *Canvas) {
	for w.m = 0; w.m < c.Height; w.m++ {
		for w.n = 0; w.n < c.Width; w.n++ {
			w.addColorToLine(c.Pixels[w.m][w.n], c)
		}
		w.addLine()
	}
}

func (w *PPMWriter) newlineIfLineTooLong(color string) bool {
	if len(w.buf)+len(color)+1 > PPMMaxCharsPerLine {
		w.addLine()
		return true
	}
	return false
}

func (w *PPMWriter) addLine() {
	w.PixelData += w.buf + "\n"
	w.buf = ""
}

func (w *PPMWriter) addColorToLine(color m.Vec4, c *Canvas) {
	for i := 0; i < 3; i++ {
		colorStr := floatToPpm(color[i], PPMMaxColor)
		w.newlineIfLineTooLong(colorStr)
		w.buf += colorStr
		if w.shouldAddSpace(i, c) {
			w.buf += " "
		}
	}
}

func (w *PPMWriter) shouldAddSpace(colorIndex int, c *Canvas) bool {
	isLastColorOnLine := (colorIndex == 2) && (w.n == c.Width-1)
	nextColorFitsOnLine := w.getNextColor(colorIndex, c)
	if !isLastColorOnLine && nextColorFitsOnLine {
		return true
	}
	return false
}

func (w *PPMWriter) getNextColor(colorIndex int, c *Canvas) bool {
	var m, n, i uint32
	isNextPixel := colorIndex == 2
	if isNextPixel {
		i = 0
		n = w.n+1
	}
	if n >= c.Width {
		n = 0
		m = w.m+1
	}

	if m >= c.Height {
		return false
	}
	nextColorStr := floatToPpm(c.Pixels[m][n][i], PPMMaxColor)

	return len(w.buf)+len(nextColorStr)+1 <= PPMMaxCharsPerLine
}
