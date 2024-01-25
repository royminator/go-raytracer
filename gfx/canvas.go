package gfx

import (
	"fmt"
	"math"
	"os"

	m "roytracer/math"
)

type (
	Canvas struct {
		Pixels [][]m.Vec4
		Width  uint32
		Height uint32
	}

	PPMWriter struct {
		Header    []byte
		Pixels []byte
		Ppm       []byte
		m, n, i uint32
		nBytes uint32
		lastNewline uint32
		cursor uint32
		MaxLineLength uint32
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
	w.Ppm = append(w.Ppm, w.Pixels...)
}

func (w *PPMWriter) WriteHeader(c *Canvas) {
	w.Header = []byte(fmt.Sprintf("%s\n%d %d\n%d\n",
		PPMFormat, c.Width, c.Height, PPMMaxColor))
}

func (w *PPMWriter) WritePixelData(c *Canvas) {
	w.CalcBytes(c)
	w.Pixels = make([]byte, w.nBytes)

	for w.m = 0; w.m < c.Height; w.m++ {
		for w.n = 0; w.n < c.Width; w.n++ {
			color := c.Pixels[w.m][w.n]
			for w.i = 0; w.i < 3; w.i++ {
				colorStr := floatToPpm(color[w.i], PPMMaxColor)
				nextColor, _ := w.getNextColor(c)
				nextColorStr := floatToPpm(nextColor, PPMMaxColor)
				sep := w.getSeparator(colorStr, nextColorStr, c.Width)
				w.insertIntoBuffer([]byte(colorStr + sep))
			}
		}
	}
}

func (w *PPMWriter) getNextColor(c *Canvas) (float64, bool) {
	if w.i < 2 {
		return c.Pixels[w.m][w.n][w.i+1], true
	}
	if w.n+1 < c.Width {
		return c.Pixels[w.m][w.n+1][0], true
	}
	if w.m+1 < c.Height {
		return c.Pixels[w.m+1][0][0], true
	}
	return 0, false
}

func (w *PPMWriter) insertIntoBuffer(bytes []byte) {
	start, end :=w.cursor, w.cursor+uint32(len(bytes))
	copy(w.Pixels[start:end], bytes)
	if bytes[len(bytes)-1] == byte('\n') {
		w.lastNewline = end
	}
	w.cursor = end
}

func (w *PPMWriter) getSeparator(colorStr, nextColor string, width uint32) string {
	if w.isNewline(colorStr, nextColor, width) {
		return "\n"
	}
	return " "
}

func (w *PPMWriter) isNewline(color, nextColor string, width uint32) bool {
	isTooLong := w.cursor+uint32(len(color))+uint32(len(nextColor))+2-w.lastNewline >= w.MaxLineLength
	isEndOfRow := w.n == width-1 && w.i == 2
	return isTooLong || isEndOfRow
}

func (w *PPMWriter) CalcBytes(c *Canvas) {
	nPixels := c.Height*c.Width
	nColors := nPixels*3
	nChars := nColors*4
	w.nBytes = nChars
}

func (w *PPMWriter) SaveFile(filePath string) {
	os.WriteFile(filePath, []byte(w.Ppm), 0644)
}
