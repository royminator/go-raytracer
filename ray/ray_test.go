package ray

import (
	"testing"

	m "roytracer/math"

	"github.com/stretchr/testify/assert"
)

func TestCreateRay(t *testing.T) {
	origin := m.Point4(1, 2, 3)
	direction := m.Vector4(1, 0, 0)
	actual := Ray{origin, direction}
	assert.Equal(t, actual.Origin, origin)
	assert.Equal(t, actual.Dir, direction)
}

func TestComputePointOnRay(t *testing.T) {
	ray := Ray{m.Point4(2, 3, 4), m.Vector4(1, 0, 0)}
	assert := assert.New(t)
	assert.Equal(m.Point4(2, 3, 4), ray.Pos(0))
	assert.Equal(m.Point4(3, 3, 4), ray.Pos(1))
	assert.Equal(m.Point4(1, 3, 4), ray.Pos(-1))
	assert.Equal(m.Point4(4.5, 3, 4), ray.Pos(2.5))
}

func TestSphereIntersection(t *testing.T) {
	ray := Ray{m.Point4(0, 0, -5), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect.T))
	assert.Equal(4.0, isect.T[0])
	assert.Equal(6.0, isect.T[1])
}

func TestSphereIntersectionAtTangent(t *testing.T) {
	ray := Ray{m.Point4(0, 1, -5), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect.T))
	assert.Equal(5.0, isect.T[0])
	assert.Equal(5.0, isect.T[0])
}

func TestSphereIntersectionRayMisses(t *testing.T) {
	ray := Ray{m.Point4(0, 2, -5), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert.Equal(t, 0, len(isect.T))
}

func TestSphereIntersectionRayBehindSphere(t *testing.T) {
	ray := Ray{m.Point4(0, 0, 5), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect.T))
	assert.Equal(-6.0, isect.T[0])
	assert.Equal(-4.0, isect.T[1])
}

func TestSphereIntersectRayInsideSphere(t *testing.T) {
	ray := Ray{m.Point4(0, 0, 0), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect.T))
	assert.Equal(-1.0, isect.T[0])
	assert.Equal(1.0, isect.T[1])
}
