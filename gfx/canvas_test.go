package gfx

import "testing"

func TestNewCanvasIsAllWhite(t *testing.T) {
	c := NewCanvas(10, 20, ColorWhite)

	assertEqual(uint32(10), c.Width, t)
	assertEqual(uint32(20), c.Height, t)

	for i := range c.Pixels {
		for j := range c.Pixels[i] {
			assertEqual(c.Pixels[i][j], ColorWhite, t)
		}
	}
}

func TestCanvasWritePixel(t *testing.T) {
	c := NewCanvas(10, 20, ColorWhite)
	var i, j uint32 = 2, 3
	c.WritePixel(i, j, ColorRed)
	expected := ColorRed
	actual := c.Pixels[i][j]
	assertEqual(expected, actual, t)
}

/////////////////////// HELPERS /////////////////////// 
func assertEqual(expected, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
