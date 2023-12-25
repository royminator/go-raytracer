package ray

import (
	"fmt"
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
	fmt.Println(isect)
	assert := assert.New(t)
	assert.True(isect.IsIntersect)
	assert.Equal(4.0, isect.D1)
	assert.Equal(6.0, isect.D2)
}

func TestSphereIntersectionAtTangent(t *testing.T) {
	ray := Ray{m.Point4(0, 1, -5), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert := assert.New(t)
	assert.True(isect.IsIntersect)
	assert.Equal(5.0, isect.D1)
	assert.Equal(5.0, isect.D2)
}

func TestSphereIntersectionRayMisses(t *testing.T) {
	ray := Ray{m.Point4(0, 2, -5), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert.False(t, isect.IsIntersect)
}

func TestSphereIntersectionRayBehindSphere(t *testing.T) {
	ray := Ray{m.Point4(0, 0, 5), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert.Equal(t, Intersection{true, -6.0, -4.0}, isect)
}

func TestSphereIntersectRayInsideSphere(t *testing.T) {
	ray := Ray{m.Point4(0, 0, 0), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert.Equal(t, Intersection{true, -1.0, 1.0}, isect)
}
