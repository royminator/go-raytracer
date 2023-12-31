package light

import (
	"math"
	"testing"

	"roytracer/color"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/pattern"
	"roytracer/shape"

	"github.com/stretchr/testify/assert"
)

var (
	mat = mtl.DefaultMaterial()
	pos = m.Point4(0, 0, 0)
)

func teardown() {
	mat = mtl.DefaultMaterial()
	pos = m.Point4(0, 0, 0)
}

func TestPointLightHasPositionAndColor(t *testing.T) {
	l := PointLight{
		Pos:   m.Point4(0, 0, 0),
		Intensity: m.Color4(1, 1, 1, 0),
	}
	assert.Equal(t, m.Vec4{0, 0, 0, 1}, l.Pos)
	assert.Equal(t, m.Vec4{1, 1, 1, 0}, l.Intensity)
}

func TestLightingWhenEyeBetweenLightAndSurface(t *testing.T) {
	eyev := m.Vector4(0, 0, -1)
	normalv := m.Vector4(0, 0, -1)
	l := PointLight{Pos: m.Point4(0, 0, -10), Intensity: m.Vec4{1, 1, 1, 0}}
	assert.Equal(t, m.Vec4{1.9, 1.9, 1.9, 0.0}, Lighting(mat, shape.NewTestShape(), l, pos, eyev, normalv, false))
}

func TestLightingWhenEyeBetweenLightAndSurfaceWithOffset45Deg(t *testing.T) {
	eyev := m.Vector4(0, math.Sqrt2/2.0, -math.Sqrt2/2.0)
	normalv := m.Vector4(0, 0, -1)
	l := PointLight{Pos: m.Point4(0, 0, -10), Intensity: m.Vec4{1, 1, 1, 0}}
	assert.Equal(t, m.Vec4{1, 1, 1, 0}, Lighting(mat, shape.NewTestShape(), l, pos, eyev, normalv, false))
}

func TestLigthingWithEyeOppositeSurfaceLightOffset45Deg(t *testing.T) {
	eyev := m.Vector4(0, 0, -1)
	normalv := m.Vector4(0, 0, -1)
	l := PointLight{Pos: m.Point4(0, 10, -10), Intensity: m.Vec4{1, 1, 1, 0}}
	actual := Lighting(mat, shape.NewTestShape(), l, pos, eyev, normalv, false)
	expected := m.Vec4{0.7364, 0.7364, 0.7364, 0}
	assert.True(t, expected.ApproxEqual(actual))
}

func TestLightingWhenEyeInPathOfReflection(t *testing.T) {
	eyev := m.Vector4(0, -math.Sqrt2/2.0, -math.Sqrt2/2.0)
	normalv := m.Vector4(0, 0, -1)
	l := PointLight{Pos: m.Point4(0, 10, -10), Intensity: m.Vec4{1, 1, 1, 0}}
	actual := Lighting(mat, shape.NewTestShape(), l, pos, eyev, normalv, false)
	expected := m.Vec4{1.6364, 1.6364, 1.6364, 0}
	assert.True(t, expected.ApproxEqual(actual))
}

func TestLightingWhenLightBehindSurface(t *testing.T) {
	eyev := m.Vector4(0, 0, -1)
	normalv := m.Vector4(0, 0, -1)
	l := PointLight{Pos: m.Point4(0, 0, 10), Intensity: m.Vec4{1, 1, 1, 0}}
	actual := Lighting(mat, shape.NewTestShape(), l, pos, eyev, normalv, false)
	expected := m.Vec4{0.1, 0.1, 0.1, 0.0}
	assert.True(t, expected.ApproxEqual(actual))
}

func TestLightingWhenSurfaceInShadow(t *testing.T) {
	eye := m.Vector4(0, 0, -1)
	normal := m.Vector4(0, 0, -1)
	light := PointLight{
		Pos: m.Point4(0, 0, -10),
		Intensity: m.Vec4{1, 1, 1},
	}
	assert.Equal(t, m.Vec4{0.1, 0.1, 0.1}, Lighting(mat, shape.NewTestShape(), light, pos, eye, normal, true))
}

func TestLightWithPatternApplied(t *testing.T) {
	defer teardown()
	p := pattern.NewStripePattern(color.White, color.Black)
	mat.Pattern = &p
	mat.Ambient = 1
	mat.Diffuse = 0
	mat.Specular = 0
	eye := m.Vector4(0, 0, -1)
	normal := m.Vector4(0, 0, -1)
	light := PointLight{Pos: m.Point4(0, 0, -10), Intensity: color.White}
	s := shape.NewTestShape()
	s.SetMat(mat)
	c1 := Lighting(mat, s, light, m.Point4(0.9, 0, 0), eye, normal, false)
	c2 := Lighting(mat, s, light, m.Point4(1.1, 0, 0), eye, normal, false)
	assert.Equal(t, color.White, c1)
	assert.Equal(t, color.Black, c2)
}
