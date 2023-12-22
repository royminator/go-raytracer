package gfx

import (
	"fmt"
	"strings"
	"testing"

	m "roytracer/math"
)

func TestNewCanvasIsAllBlack(t *testing.T) {
	c := NewCanvas(10, 20, ColorBlack)

	assertEqual(uint32(10), c.Width, t)
	assertEqual(uint32(20), c.Height, t)

	for i := range c.Pixels {
		for j := range c.Pixels[i] {
			assertEqual(c.Pixels[i][j], ColorBlack, t)
		}
	}
}

func TestCanvasWritePixel(t *testing.T) {
	c := NewCanvas(10, 20, ColorBlack)
	var i, j uint32 = 2, 3
	c.WritePixel(j, i, ColorRed)
	expected := ColorRed
	actual := c.Pixels[i][j]
	assertEqual(expected, actual, t)
}

func TestPPMWriterHeader(t *testing.T) {
	expected := lines("P3\n5 3\n255")
	c := NewCanvas(5, 3, ColorBlack)
	ppmWriter := PPMWriter{}
	ppmWriter.Write(c)
	actual := lines(string(ppmWriter.Ppm))[0:3]
	assertEqual(len(expected), len(actual), t)
	for i := range expected {
		assertEqual(expected[i], actual[i], t)
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

	ppmWriter := PPMWriter{}
	ppmWriter.Write(c)
	fmt.Println(string(ppmWriter.Ppm))
	actual := lines(string(ppmWriter.Ppm))[3:5]

	assertEqual(len(expected), len(actual), t)
	for i := range expected {
		assertEqual(expected[i], actual[i], t)
	}
}

func TestCanvasPPMWriterPixelDataMax70CharPerLine(t *testing.T) {
	expected := lines(`255 204 153 255 204 153 255 204 153 255 204 153 255 204 153 255 204
153 255 204 153 255 204 153 255 204 153 255 204 153
255 204 153 255 204 153 255 204 153 255 204 153 255 204 153 255 204
153 255 204 153 255 204 153 255 204 153 255 204 153`)

	c := NewCanvas(10, 2, m.Color4(1, 0.8, 0.6, 0.0))
	ppmWriter := PPMWriter{}
	ppmWriter.Write(c)
	actual := lines(string(ppmWriter.Ppm))[3:7]

	assertEqual(len(expected), len(actual), t)
	for i := range expected {
		assertEqual(expected[i]+";", actual[i]+";", t)
	}
}

// ///////////////////// HELPERS ///////////////////////
func assertEqual(expected, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func lines(text string) []string {
	return strings.Split(text, "\n")
}
