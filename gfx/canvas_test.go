package gfx

import (
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

func TestCanvasPPMHeader(t *testing.T) {
	expected := lines("P3\n5 3\n255\n", 0, 2)
	actual := lines(NewCanvas(5, 3, ColorBlack).ToPPM(), 0, 2)
	if expected != actual {
		t.Errorf("Expected %s, got %s", expected, actual)
	}
}

func TestCanvasPPMPixelData(t *testing.T) {
	expected := `255 0 0 0 0 0 0 0 0 0 0 0 0 0 0
0 0 0 0 0 0 0 128 0 0 0 0 0 0 0
0 0 0 0 0 0 0 0 0 0 0 0 0 0 255`

	c := NewCanvas(5, 3, ColorBlack)
	c.WritePixel(0, 0, m.Color4(1.5, 0, 0, 0))
	c.WritePixel(2, 1, m.Color4(0, 0.5, 0, 0))
	c.WritePixel(4, 2, m.Color4(-0.5, 0, 1, 0))

	actual := lines(c.ToPPM(), 3, 6)

	if expected != actual {
		t.Errorf("Expected %s, got %s", []byte(expected), []byte(actual))
	}
}

// ///////////////////// HELPERS ///////////////////////
func assertEqual(expected, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}

func lines(text string, from, to int) string {
	ls := strings.Split(text, "\n")
	var res string
	for i := from; i < to; i++ {
		res += ls[i]
		if i != to - 1 {
			res += "\n"
		}
	}
	return res
}
