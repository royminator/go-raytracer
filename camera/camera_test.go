package camera

import (
	"math"
	"testing"

	m "roytracer/math"
	"roytracer/world"

	"github.com/stretchr/testify/assert"
)

func TestCreateCamera(t *testing.T) {
	hsize := 160
	vsize := 120
	fov := math.Pi / 2.0
	cam := NewCamera(hsize, vsize, fov)
	assert := assert.New(t)
	assert.Equal(hsize, cam.Hsize)
	assert.Equal(vsize, cam.Vsize)
	assert.Equal(fov, cam.Fov)
	assert.Equal(m.Mat4Ident(), cam.Tf)
}

func TestPixelSizeForHorizontalCanvas(t *testing.T) {
	cam := NewCamera(200, 125, math.Pi/2.0)
	assert.Equal(t, 0.01, cam.PixelSize)
}

func TestPixelSizeForVerticalCanvas(t *testing.T) {
	cam := NewCamera(125, 200, math.Pi/2.0)
	assert.Equal(t, 0.01, cam.PixelSize)
}

func TestConstructingARayThroughTheCenterOfTheCamera(t *testing.T) {
	cam := NewCamera(201, 101, math.Pi/2.0)
	r := cam.RayForPixel(100, 50)
	assert.Equal(t, m.Point4(0, 0, 0), r.Origin)
	assert.Equal(t, m.Vector4(0, 0, -1), r.Dir)
}

func TestConstructingARayThroughCornerOfCanvas(t *testing.T) {
	cam := NewCamera(201, 101, math.Pi/2.0)
	r := cam.RayForPixel(0, 0)
	assert.Equal(t, m.Point4(0, 0, 0), r.Origin)
	assert.True(t, m.Vector4(0.66519, 0.33259, -0.66851).ApproxEqual(r.Dir))
}

func TestConstructingARayWhenCameraIsTransformed(t *testing.T) {
	cam := NewCamera(201, 101, math.Pi/2.0)
	cam.SetTf(m.RotY(math.Pi / 4.0).Mul(m.Trans(0, -2, 5)))
	r := cam.RayForPixel(100, 50)
	assert.Equal(t, m.Point4(0, 2, -5), r.Origin)
	assert.True(t, m.Vector4(math.Sqrt2/2.0, 0, -math.Sqrt2/2.0).ApproxEqual(r.Dir))
}

func TestRenderDefaultWorld(t *testing.T) {
	w := world.DefaultWorld()
	cam := NewCamera(11, 11, math.Pi/2.0)
	from := m.Point4(0, 0, -5)
	to := m.Point4(0, 0, 0)
	up := m.Vector4(0, 1, 0)
	cam.SetTf(m.View(from, to ,up))
	image := cam.Render(w)
	assert.True(t, m.Vec4{0.38066, 0.47583, 0.2855}.ApproxEqual(image.Pixels[5][5]))
}
