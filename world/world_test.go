package world

import (
	"math"
	"testing"

	"roytracer/color"
	"roytracer/light"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/pattern"
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
		Dir:    m.Vector4(0, 0, 1),
	}
	s := w.Objects[0]
	isect := shape.Intersection{
		S: s,
		T: 4,
	}
	comps := isect.Prepare(r, []shape.Intersection{isect})
	c := w.ShadeHit(comps, 4)
	assert.True(t, c.ApproxEqual(m.Vec4{0.38066, 0.47583, 0.2855, 0}))
}

func TestShadingIntersectionFromInside(t *testing.T) {
	w := DefaultWorld()
	w.Light = light.PointLight{
		Pos:       m.Point4(0, 0.25, 0),
		Intensity: m.Vec4{1, 1, 1, 0},
	}
	r := shape.Ray{
		Origin: m.Point4(0, 0, 0),
		Dir:    m.Vector4(0, 0, 1),
	}
	s := w.Objects[1]
	isect := shape.Intersection{
		S: s,
		T: 0.5,
	}
	comps := isect.Prepare(r, []shape.Intersection{isect})
	c := w.ShadeHit(comps, 4)
	assert.True(t, c.ApproxEqual(m.Vec4{0.90498, 0.90498, 0.90498, 0}))
}

func TestColorWhenRayMisses(t *testing.T) {
	w := DefaultWorld()
	r := shape.Ray{
		Origin: m.Point4(0, 0, -5),
		Dir:    m.Vector4(0, 1, 0),
	}
	c := w.ColorAt(r, 4)
	assert.Equal(t, m.Vec4With(0), c)
}

func TestColorWhenRayHits(t *testing.T) {
	w := DefaultWorld()
	r := shape.Ray{
		Origin: m.Point4(0, 0, -5),
		Dir:    m.Vector4(0, 0, 1),
	}
	c := w.ColorAt(r, 4)
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
		Dir:    m.Vector4(0, 0, -1),
	}
	c := w.ColorAt(r, 4)
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
			Pos:       m.Point4(0, 0, -10),
			Intensity: m.Vec4{1, 1, 1},
		},
		Objects: []shape.Shape{&s1, &s2},
	}
	ray := shape.Ray{Origin: m.Point4(0, 0, 5), Dir: m.Vector4(0, 0, 1)}
	i := shape.Intersection{T: 4, S: &s2}
	comps := i.Prepare(ray, []shape.Intersection{i})
	c := w.ShadeHit(comps, 4)
	assert.Equal(t, m.Vec4{0.1, 0.1, 0.1}, c)
}

func TestReflectedColorForNonreflectiveMaterial(t *testing.T) {
	w := DefaultWorld()
	r := shape.Ray{
		Origin: m.Point4(0, 0, 0),
		Dir:    m.Vector4(0, 0, 1),
	}
	s := w.Objects[1]
	mat := s.GetMat()
	mat.Ambient = 1.0
	s.SetMat(mat)
	isect := shape.Intersection{T: 1.0, S: s}
	comps := isect.Prepare(r, []shape.Intersection{isect})
	c := w.ReflectedColor(comps, 4)
	assert.Equal(t, m.Vec4With(0), c)
}

func TestReflectedColorForReflectiveSurface(t *testing.T) {
	w := DefaultWorld()
	p := shape.NewPlane()
	p.O.Material.Reflective = 0.5
	p.SetTf(m.Trans(0, -1, 0))
	w.AddShape(&p)
	ray := shape.Ray{
		Origin: m.Point4(0, 0, -3),
		Dir:    m.Vector4(0, -math.Sqrt2/2.0, math.Sqrt2/2.0),
	}
	isect := shape.Intersection{S: &p, T: math.Sqrt2}
	comps := isect.Prepare(ray, []shape.Intersection{isect})
	actual := w.ReflectedColor(comps, 4)
	assert.True(t, actual.ApproxEqual(m.Vec4{0.19032, 0.2379, 0.14274}))
}

func TestShadeHitWithReflectiveMaterial(t *testing.T) {
	w := DefaultWorld()
	p := shape.NewPlane()
	p.O.Material.Reflective = 0.5
	p.SetTf(m.Trans(0, -1, 0))
	w.AddShape(&p)
	ray := shape.Ray{
		Origin: m.Point4(0, 0, -3),
		Dir:    m.Vector4(0, -math.Sqrt2/2.0, math.Sqrt2/2.0),
	}
	isect := shape.Intersection{S: &p, T: math.Sqrt2}
	comps := isect.Prepare(ray, []shape.Intersection{isect})
	actual := w.ShadeHit(comps, 4)
	assert.True(t, actual.ApproxEqual(m.Vec4{0.87677, 0.92436, 0.82918}))
}

func TestColorAtWithMutuallyReflectiveSurfaces(t *testing.T) {
	w := World{
		Light: light.PointLight{
			Pos:       m.Point4(0, 0, 0),
			Intensity: color.White,
		},
	}
	lower := shape.NewPlane()
	lower.O.Material.Reflective = 1.0
	lower.SetTf(m.Trans(0, -1, 0))
	w.AddShape(&lower)
	upper := shape.NewPlane()
	upper.O.Material.Reflective = 1.0
	upper.SetTf(m.Trans(0, 1, 0))
	w.AddShape(&upper)
	ray := shape.Ray{
		Origin: m.Point4(0, 0, 0),
		Dir:    m.Vector4(0, 1, 0),
	}
	w.ColorAt(ray, 4)
}

func TestReflectedColorAtTheMaximumRecursiveDepth(t *testing.T) {
	w := DefaultWorld()
	p := shape.NewPlane()
	p.O.Material.Reflective = 0.5
	p.SetTf(m.Trans(0, -1, 0))
	w.AddShape(&p)
	ray := shape.Ray{
		Origin: m.Point4(0, 0, -3),
		Dir:    m.Vector4(0, -math.Sqrt2/2.0, math.Sqrt2/2.0),
	}
	isect := shape.Intersection{S: &p, T: math.Sqrt2}
	comps := isect.Prepare(ray, []shape.Intersection{isect})
	actual := w.ReflectedColor(comps, 0)
	assert.Equal(t, actual, m.Vec4{})
}

func TestRefractedColorWithOpaqueSurface(t *testing.T) {
	w := DefaultWorld()
	s := w.Objects[0]
	r := shape.Ray{Origin: m.Point4(0, 0, -5), Dir: m.Vector4(0, 0, 1)}
	xs := []shape.Intersection{
		{T: 4, S: s},
		{T: 6, S: s},
	}
	comps := xs[0].Prepare(r, xs)
	c := w.RefractedColor(comps, 5)
	assert.Equal(t, m.Vec4{}, c)
}

func TestRefractedColorAtTheMaximumRecursiveDepth(t *testing.T) {
	w := DefaultWorld()
	s := w.Objects[0]
	mtl := s.GetMat()
	mtl.Transparency = 1.0
	mtl.RefractiveIndex = 1.5
	s.SetMat(mtl)
	r := shape.Ray{Origin: m.Point4(0, 0, -5), Dir: m.Vector4(0, 0, 1)}
	xs := []shape.Intersection{
		{S: s, T: 4},
		{S: s, T: 6},
	}
	comps := xs[0].Prepare(r, xs)
	c := w.RefractedColor(comps, 0)
	assert.Equal(t, color.Black, c)
}

func TestRefractedColorUnderTotalInternalReflection(t *testing.T) {
	w := DefaultWorld()
	s := w.Objects[0]
	mtl := s.GetMat()
	mtl.Transparency = 1.0
	mtl.RefractiveIndex = 1.5
	s.SetMat(mtl)
	r := shape.Ray{
		Origin: m.Point4(0, 0, math.Sqrt2/2.0),
		Dir:    m.Vector4(0, 1, 0),
	}
	xs := []shape.Intersection{
		{T: -math.Sqrt2 / 2.0, S: s},
		{T: math.Sqrt2 / 2.0, S: s},
	}
	comps := xs[1].Prepare(r, xs)
	c := w.RefractedColor(comps, 5)
	assert.Equal(t, color.Black, c)
}

func TestRefractedColorWithARefractedRay(t *testing.T) {
	w := DefaultWorld()
	a := w.Objects[0]
	p := pattern.NewTestPattern()
	mtl := a.GetMat()
	mtl.Ambient = 1.0
	mtl.Pattern = &p
	a.SetMat(mtl)
	b := w.Objects[1]
	mtlB := b.GetMat()
	mtlB.Transparency = 1.0
	mtlB.RefractiveIndex = 1.5
	b.SetMat(mtlB)

	r := shape.Ray{
		Origin: m.Point4(0, 0, 0.1),
		Dir:    m.Vector4(0, 1, 0),
	}
	xs := []shape.Intersection{
		{S: a, T: -0.9899},
		{S: b, T: -0.4899},
		{S: b, T: 0.4899},
		{S: a, T: 0.9899},
	}
	comps := xs[2].Prepare(r, xs)
	c := w.RefractedColor(comps, 5)
	assert.True(t, c.ApproxEqual(m.Vec4{0, 0.99888, 0.04725}))
}

func TestShadeHitWithATransparentMaterial(t *testing.T) {
	w := DefaultWorld()
	floor := shape.NewPlane()
	floor.SetTf(m.Trans(0, -1, 0))
	floor.O.Material.Transparency = 0.5
	floor.O.Material.RefractiveIndex = 1.5
	w.AddShape(&floor)

	ball := shape.NewSphere()
	ball.O.Material.Color = m.Vec4{1, 0, 0}
	ball.O.Material.Ambient = 0.5
	ball.SetTf(m.Trans(0, -3.5, -0.5))
	w.AddShape(&ball)

	r := shape.Ray{
		Origin: m.Point4(0, 0, -3),
		Dir:    m.Vector4(0, -math.Sqrt2/2.0, math.Sqrt2/2.0),
	}
	xs := []shape.Intersection{{T: math.Sqrt2, S: &floor}}
	comps := xs[0].Prepare(r, xs)
	color := w.ShadeHit(comps, 5)
	assert.True(t, m.Vec4{0.93642, 0.68642, 0.68642}.ApproxEqual(color))
}

func TestShadeHitWithReflectiveTransparentMaterial(t *testing.T) {
	w := DefaultWorld()
	r := shape.Ray{
		Origin: m.Point4(0, 0, -3),
		Dir: m.Vector4(0, -math.Sqrt2/2.0, math.Sqrt2/2.0),
	}
	floor := shape.NewPlane()
	floor.SetTf(m.Trans(0, -1, 0))
	floor.O.Material.Reflective = 0.5
	floor.O.Material.Transparency = 0.5
	floor.O.Material.RefractiveIndex = 1.5
	w.AddShape(&floor)

	ball := shape.NewSphere()
	ball.SetTf(m.Trans(0, -3.5, -0.5))
	ball.O.Material.Color = m.Vec4{1, 0, 0}
	ball.O.Material.Ambient = 0.5
	w.AddShape(&ball)

	xs := []shape.Intersection{{T: math.Sqrt2, S: &floor}}
	comps := xs[0].Prepare(r, xs)
	color := w.ShadeHit(comps, 5)
	assert.True(t, color.ApproxEqual(m.Vec4{0.93391, 0.69643, 0.69243}))
}
