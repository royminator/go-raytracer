package mtl

import (
	"math"
	"testing"

	m "roytracer/math"

	"github.com/stretchr/testify/assert"
)

func TestPointLightHasPositionAndColor(t *testing.T) {
	l := PointLight{
		Pos:   m.Point4(0, 0, 0),
		Intensity: m.Color4(1, 1, 1, 0),
	}
	assert.Equal(t, m.Vec4{0, 0, 0, 1}, l.Pos)
	assert.Equal(t, m.Vec4{1, 1, 1, 0}, l.Intensity)
}

func TestLightingWhenEyeBetweenLightAndSurface(t *testing.T) {
	mat := DefaultMaterial()
	pos := m.Point4(0, 0, 0)
	eyev := m.Vector4(0, 0, -1)
	normalv := m.Vector4(0, 0, -1)
	l := PointLight{Pos: m.Point4(0, 0, -10), Intensity: m.Vec4{1, 1, 1, 0}}
	assert.Equal(t, m.Vec4{1.9, 1.9, 1.9, 0.0}, Lighting(mat, l, pos, eyev, normalv, false))
}

func TestLightingWhenEyeBetweenLightAndSurfaceWithOffset45Deg(t *testing.T) {
	mat := DefaultMaterial()
	pos := m.Point4(0, 0, 0)
	eyev := m.Vector4(0, math.Sqrt2/2.0, -math.Sqrt2/2.0)
	normalv := m.Vector4(0, 0, -1)
	l := PointLight{Pos: m.Point4(0, 0, -10), Intensity: m.Vec4{1, 1, 1, 0}}
	assert.Equal(t, m.Vec4{1, 1, 1, 0}, Lighting(mat, l, pos, eyev, normalv, false))
}

func TestLigthingWithEyeOppositeSurfaceLightOffset45Deg(t *testing.T) {
	mat := DefaultMaterial()
	pos := m.Point4(0, 0, 0)
	eyev := m.Vector4(0, 0, -1)
	normalv := m.Vector4(0, 0, -1)
	l := PointLight{Pos: m.Point4(0, 10, -10), Intensity: m.Vec4{1, 1, 1, 0}}
	actual := Lighting(mat, l, pos, eyev, normalv, false)
	expected := m.Vec4{0.7364, 0.7364, 0.7364, 0}
	assert.True(t, expected.ApproxEqual(actual))
}

func TestLightingWhenEyeInPathOfReflection(t *testing.T) {
	mat := DefaultMaterial()
	pos := m.Point4(0, 0, 0)
	eyev := m.Vector4(0, -math.Sqrt2/2.0, -math.Sqrt2/2.0)
	normalv := m.Vector4(0, 0, -1)
	l := PointLight{Pos: m.Point4(0, 10, -10), Intensity: m.Vec4{1, 1, 1, 0}}
	actual := Lighting(mat, l, pos, eyev, normalv, false)
	expected := m.Vec4{1.6364, 1.6364, 1.6364, 0}
	assert.True(t, expected.ApproxEqual(actual))
}

func TestLightingWhenLightBehindSurface(t *testing.T) {
	mat := DefaultMaterial()
	pos := m.Point4(0, 0, 0)
	eyev := m.Vector4(0, 0, -1)
	normalv := m.Vector4(0, 0, -1)
	l := PointLight{Pos: m.Point4(0, 0, 10), Intensity: m.Vec4{1, 1, 1, 0}}
	actual := Lighting(mat, l, pos, eyev, normalv, false)
	expected := m.Vec4{0.1, 0.1, 0.1, 0.0}
	assert.True(t, expected.ApproxEqual(actual))
}

func TestLightingWhenSurfaceInShadow(t *testing.T) {
	mat := DefaultMaterial()
	pos := m.Point4(0, 0, 0)
	eye := m.Vector4(0, 0, -1)
	normal := m.Vector4(0, 0, -1)
	light := PointLight{
		Pos: m.Point4(0, 0, -10),
		Intensity: m.Vec4{1, 1, 1},
	}
	assert.Equal(t, m.Vec4{0.1, 0.1, 0.1}, Lighting(mat, light, pos, eye, normal, true))
}
