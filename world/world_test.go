package world

import (
	"testing"

	"roytracer/light"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/shape"

	"github.com/stretchr/testify/assert"
)

func TestCreateWorld(t *testing.T) {
	w := World{}
	assert.Empty(t, w.Objects)
}

func TestDefaultWorld(t *testing.T) {
	w := DefaultWorld()
	assert.Len(t, w.Objects, 2)
}

func TestIntersectWorld(t *testing.T) {
	w := DefaultWorld()
	ray := shape.Ray{
		Origin: m.Point4(0, 0, -5),
		Dir:    m.Vector4(0, 0, 1),
	}
	isects := w.Intersect(ray)
	assert := assert.New(t)
	assert.Len(isects, 4)
	assert.Equal(4.0, isects[0].T)
	assert.Equal(4.5, isects[1].T)
	assert.Equal(5.5, isects[2].T)
	assert.Equal(6.0, isects[3].T)
}

func TestShadingIntersection(t *testing.T) {
	w := DefaultWorld()
	r := shape.Ray{
		Origin: m.Point4(0, 0, -5),
		Dir: m.Vector4(0, 0, 1),
	}
	s := w.Objects[0]
	isect := shape.Intersection{
		S: s,
		T: 4,
	}
	comps := isect.Prepare(r)
	c := w.ShadeHit(comps)
	assert.True(t, c.ApproxEqual(m.Vec4{0.38066, 0.47583, 0.2855, 0}))
}

func TestShadingIntersectionFromInside(t *testing.T) {
	w := DefaultWorld()
	w.Light = light.PointLight{
		Pos: m.Point4(0, 0.25, 0),
		Intensity: m.Vec4{1, 1, 1, 0},
	}
	r := shape.Ray{
		Origin: m.Point4(0, 0, 0),
		Dir: m.Vector4(0, 0, 1),
	}
	s := w.Objects[1]
	isect := shape.Intersection{
		S: s,
		T: 0.5,
	}
	comps := isect.Prepare(r)
	c := w.ShadeHit(comps)
	assert.True(t, c.ApproxEqual(m.Vec4{0.90498, 0.90498, 0.90498, 0}))
}

func TestColorWhenRayMisses(t *testing.T) {
	w := DefaultWorld()
	r := shape.Ray{
		Origin: m.Point4(0, 0, -5),
		Dir: m.Vector4(0, 1, 0),
	}
	c := w.ColorAt(r)
	assert.Equal(t, m.Vec4With(0), c)
}

func TestColorWhenRayHits(t *testing.T) {
	w := DefaultWorld()
	r := shape.Ray{
		Origin: m.Point4(0, 0, -5),
		Dir: m.Vector4(0, 0, 1),
	}
	c := w.ColorAt(r)
	assert.True(t, c.ApproxEqual(m.Vec4{0.38066, 0.47583, 0.2855, 0.0}))
}

func TestColorWithIntersectionBehindRay(t *testing.T) {
	w := DefaultWorld()
	mat := mtl.Material{Ambient: 1.0}
	outer := w.Objects[0]
	inner := w.Objects[1]
	outer.SetMat(mat)
	inner.SetMat(mat)
	r := shape.Ray{
		Origin: m.Point4(0, 0, 0.75),
		Dir: m.Vector4(0, 0, -1),
	}
	c := w.ColorAt(r)
	assert.Equal(t, c, inner.GetMat().Color)
}

func TestNoShadowWhenNothingCollinearWithPointAndLight(t *testing.T) {
	assert.False(t, DefaultWorld().IsShadowed(m.Point4(0, 10, 0)))
}

func TestShadowWhenSomethingIsBetweenPointAndLight(t *testing.T) {
	assert.True(t, DefaultWorld().IsShadowed(m.Point4(10, -10, 10)))
}

func TestNoShadowWhenObjectBehindLight(t *testing.T) {
	assert.False(t, DefaultWorld().IsShadowed(m.Point4(-20, 20, -20)))
}

func TestNoShadowWhenAnObjectIsBehindThePoint(t *testing.T) {
	assert.False(t, DefaultWorld().IsShadowed(m.Point4(-2, 2, -2)))
}

func TestShadeHitWhenIntersectionInShadow(t *testing.T) {
	s1 := shape.NewSphere()
	s2 := shape.NewSphere()
	s2.SetTf(m.Trans(0, 0, 10))
	w := World{
		Light: light.PointLight{
			Pos: m.Point4(0, 0, -10),
			Intensity: m.Vec4{1, 1, 1},
		},
		Objects: []shape.Shape{&s1, &s2},
	}
	ray := shape.Ray{Origin: m.Point4(0, 0, 5), Dir: m.Vector4(0, 0, 1)}
	i := shape.Intersection{T: 4, S: &s2}
	comps := i.Prepare(ray)
	c := w.ShadeHit(comps)
	assert.Equal(t, m.Vec4{0.1, 0.1, 0.1}, c)
}
