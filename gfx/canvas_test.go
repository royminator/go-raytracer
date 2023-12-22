package gfx

import (
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"

	m "roytracer/math"
)

func TestNewCanvasIsAllBlack(t *testing.T) {
	c := NewCanvas(10, 20, ColorBlack)

	assert := assert.New(t)
	assert.Equal(uint32(10), c.Width)
	assert.Equal(uint32(20), c.Height)

	for i := range c.Pixels {
		for j := range c.Pixels[i] {
			assert.Equal(c.Pixels[i][j], ColorBlack)
		}
	}
}

func TestCanvasWritePixel(t *testing.T) {
	c := NewCanvas(10, 20, ColorBlack)
	var i, j uint32 = 2, 3
	c.WritePixel(j, i, ColorRed)
	expected := ColorRed
	actual := c.Pixels[i][j]
	assert.Equal(t, expected, actual)
}

func TestPPMWriterHeader(t *testing.T) {
	expected := lines("P3\n5 3\n255")
	c := NewCanvas(5, 3, ColorBlack)
	ppmWriter := PPMWriter{MaxLineLength: 70}
	ppmWriter.Write(c)
	actual := lines(string(ppmWriter.Ppm))[0:3]
	assert.Equal(t, len(expected), len(actual))
	for i := range expected {
		assert.Equal(t, expected[i], actual[i])
	}
}

func TestPPMWriterPixelData(t *testing.T) {
	expected := lines(`255 0 0 0 0 0 0 0 0 0 0 0 0 0 0
0 0 0 0 0 0 0 128 0 0 0 0 0 0 0
0 0 0 0 0 0 0 0 0 0 0 0 0 0 255`)

	c := NewCanvas(5, 3, ColorBlack)
	c.WritePixel(0, 0, m.Color4(1.5, 0, 0, 0))
	c.WritePixel(2, 1, m.Color4(0, 0.5, 0, 0))
	c.WritePixel(4, 2, m.Color4(-0.5, 0, 1, 0))

	ppmWriter := PPMWriter{MaxLineLength: 70}
	ppmWriter.Write(c)
	actual := lines(string(ppmWriter.Ppm))[3:6]

	assert.Equal(t, len(expected), len(actual))
	for i := range expected {
		assert.Equal(t, expected[i], actual[i])
	}
}

func TestCanvasPPMWriterPixelDataMax70CharPerLine(t *testing.T) {
	expected := lines(`255 204 153 255 204 153 255 204 153 255 204 153 255 204 153 255 204
153 255 204 153 255 204 153 255 204 153 255 204 153
255 204 153 255 204 153 255 204 153 255 204 153 255 204 153 255 204
153 255 204 153 255 204 153 255 204 153 255 204 153`)

	c := NewCanvas(10, 2, m.Color4(1, 0.8, 0.6, 0.0))
	ppmWriter := PPMWriter{MaxLineLength: 70}
	ppmWriter.Write(c)
	actual := lines(string(ppmWriter.Ppm))[3:7]

	assert.Equal(t, len(expected), len(actual))
	for i := range expected {
		assert.Equal(t, expected[i], actual[i])
	}
}

func TestPPWWriterShouldInsertNewlineAtLinesEnd(t *testing.T) {
	c := NewCanvas(3, 2, ColorBlack)
	w := PPMWriter{MaxLineLength: 70}
	w.Write(c)
	assert := assert.New(t)
	assert.Equal(byte('\n'), w.Pixels[17])
	assert.Equal(byte('\n'), w.Pixels[35])
}

func TestPPWWriterShouldInsertNewlineWhenTooLong(t *testing.T) {
	c := NewCanvas(1, 2, ColorBlack)
	w := PPMWriter{MaxLineLength: 5}
	w.Write(c)
	assert := assert.New(t)
	assert.Equal(byte('\n'), w.Pixels[5])
}

// ///////////////////// HELPERS ///////////////////////
func lines(text string) []string {
	return strings.Split(text, "\n")
}
