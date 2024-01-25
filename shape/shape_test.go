package shape

import (
	"math"
	"testing"

	"roytracer/color"
	m "roytracer/math"
	"roytracer/pattern"

	"roytracer/mtl"

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
	s := NewSphere()
	isect := s.Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect))
	assert.Equal(4.0, isect[0].T)
	assert.Equal(6.0, isect[1].T)
}

func TestSphereIntersectionAtTangent(t *testing.T) {
	ray := Ray{m.Point4(0, 1, -5), m.Vector4(0, 0, 1)}
	s := NewSphere()
	isect := s.Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect))
	assert.Equal(5.0, isect[0].T)
	assert.Equal(5.0, isect[1].T)
}

func TestSphereIntersectionRayMisses(t *testing.T) {
	ray := Ray{m.Point4(0, 2, -5), m.Vector4(0, 0, 1)}
	s := NewSphere()
	isect := s.Intersect(ray)
	assert.Equal(t, 0, len(isect))
}

func TestSphereIntersectionRayBehindSphere(t *testing.T) {
	ray := Ray{m.Point4(0, 0, 5), m.Vector4(0, 0, 1)}
	s := NewSphere()
	isect := s.Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect))
	assert.Equal(-6.0, isect[0].T)
	assert.Equal(-4.0, isect[1].T)
}

func TestSphereIntersectRayInsideSphere(t *testing.T) {
	ray := Ray{m.Point4(0, 0, 0), m.Vector4(0, 0, 1)}
	s := NewSphere()
	isect := s.Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect))
	assert.Equal(-1.0, isect[0].T)
	assert.Equal(1.0, isect[1].T)
}

func TestIntersectStoresObject(t *testing.T) {
	s := NewSphere()
	isect := Intersection{S: &s, T: 3.5}
	assert.Equal(t, &s, isect.S)
	assert.Equal(t, 3.5, isect.T)
}

func TestInterectionAggregateTValues(t *testing.T) {
	s := NewSphere()
	i1 := Intersection{S: &s, T: 3.5}
	i2 := Intersection{S: &s, T: 5.5}
	is := Intersections(i1, i2)
	assert.Equal(t, 2, len(is))
	assert.Equal(t, 3.5, is[0].T)
	assert.Equal(t, 5.5, is[1].T)
}

func TestIntersectRaySetsObject(t *testing.T) {
	s := NewSphere()
	ray := Ray{m.Point4(0, 0, -5), m.Vector4(0, 0, 1)}
	isects := s.Intersect(ray)
	assert.Equal(t, 2, len(isects))
	assert.Equal(t, &s, isects[0].S)
	assert.Equal(t, &s, isects[1].S)
}

func TestHitAllIntersectionGT0(t *testing.T) {
	type testData struct {
		res    Intersection
		isects []Intersection
		any    bool
	}

	s := NewSphere()
	td := []testData{
		{
			isects: []Intersection{
				{T: 1.0}, {T: 2},
			},
			any: true,
			res: Intersection{S: &s, T: 1},
		},
		{
			isects: []Intersection{
				{T: -1}, {T: 1},
			},
			any: true,
			res: Intersection{S: &s, T: 1},
		},
		{
			isects: []Intersection{
				{T: -2}, {T: -1},
			},
			any: false,
			res: Intersection{},
		},
		{
			isects: []Intersection{
				{T: 5}, {T: 7}, {T: -3}, {T: 2},
			},
			any: true,
			res: Intersection{S: &s, T: 2.0},
		},
	}

	assert := assert.New(t)
	for _, d := range td {
		hit, isHit := Hit(d.isects)
		assert.Equal(d.any, isHit)
		if isHit {
			assert.Equal(d.res.T, hit.T)
		}
	}
}

func TestTranslateRay(t *testing.T) {
	ray := Ray{m.Point4(1, 2, 3), m.Vector4(0, 1, 0)}
	tf := m.Trans(3, 4, 5)
	expected := Ray{m.Point4(4, 6, 8), m.Vector4(0, 1, 0)}
	ray.Transform(tf)
	assert.Equal(t, expected, ray)
}

func TestScaleRay(t *testing.T) {
	ray := Ray{m.Point4(1, 2, 3), m.Vector4(0, 1, 0)}
	tf := m.Scale(2, 3, 4)
	expected := Ray{m.Point4(2, 6, 12), m.Vector4(0, 3, 0)}
	ray.Transform(tf)
	assert.Equal(t, expected, ray)
}

func TestSpheresDefaultTransformIsIdentity(t *testing.T) {
	assert.Equal(t, m.Mat4Ident(), NewSphere().O.Tf)
}

func TestSphereSetTransform(t *testing.T) {
	s := NewSphere()
	expected := m.Trans(2, 3, 4)
	s.SetTf(expected)
	assert.Equal(t, expected, s.O.Tf)
}

func TestIntersectScaledShape(t *testing.T) {
	ray := Ray{m.Point4(0, 0, -5), m.Vector4(0, 0, 1)}
	s := NewTestShape()
	s.SetTf(m.Scale(2, 2, 2))
	s.Intersect(ray)
	assert.Equal(t, m.Point4(0, 0, -2.5), s.GetSavedRay().Origin)
	assert.Equal(t, m.Vector4(0, 0, 0.5), s.GetSavedRay().Dir)
}

func TestIntersectTranslatedShape(t *testing.T) {
	ray := Ray{m.Point4(0, 0, -5), m.Vector4(0, 0, 1)}
	s := NewTestShape()
	s.SetTf(m.Trans(5, 0, 0))
	s.Intersect(ray)
	assert.Equal(t, m.Point4(-5, 0, -5), s.GetSavedRay().Origin)
	assert.Equal(t, m.Vector4(0, 0, 1), s.GetSavedRay().Dir)
}

func TestSphereNormal(t *testing.T) {
	type testData struct {
		p   m.Vec4
		res m.Vec4
	}
	k := math.Sqrt(3.0) / 3.0
	td := []testData{
		{m.Point4(1, 0, 0), m.Vector4(1, 0, 0)},
		{m.Point4(0, 1, 0), m.Vector4(0, 1, 0)},
		{m.Point4(k, k, k), m.Vector4(k, k, k)},
	}

	s := NewSphere()

	assert := assert.New(t)
	for _, d := range td {
		assert.Equal(d.res, s.NormalAt(d.p))
	}
}

func TestSphereNormalIsNormalized(t *testing.T) {
	s := NewSphere()
	k := math.Sqrt(3.0) / 3.0
	n := s.NormalAt(m.Point4(k, k, k))
	assert.Equal(t, n, n.Normalize())
}

func TestSphereNormalOnTranslatedSphere(t *testing.T) {
	s := NewSphere()
	s.SetTf(m.Trans(0, 1, 0))
	expected := m.Vector4(0, 0.70711, -0.70711)
	assert.True(t, expected.ApproxEqual(s.NormalAt(m.Point4(0, 1.70711, -0.70711))))
}

func TestSphereNormalOnTransformedSphere(t *testing.T) {
	s := NewSphere()
	s.SetTf(m.Scale(1, 0.5, 1).Mul(m.RotZ(math.Pi / 5.0)))
	expected := m.Vector4(0, 0.97014, -0.24254)
	assert.True(t, expected.ApproxEqual(s.NormalAt(m.Point4(0, math.Sqrt2/2.0, -math.Sqrt2/2.0))))
}

func TestSphereHasDefaultMaterial(t *testing.T) {
	s := NewSphere()
	assert.Equal(t, mtl.DefaultMaterial(), s.O.Material)
}

func TestSphereMayBeAssignedMaterial(t *testing.T) {
	s := NewSphere()
	material := mtl.Material{Ambient: 1}
	s.O.Material = material
	assert.Equal(t, material, s.O.Material)
}

func TestPreComputeStateOfIntersection(t *testing.T) {
	ray := Ray{
		Origin: m.Point4(0, 0, -5),
		Dir:    m.Vector4(0, 0, 1),
	}
	s := NewSphere()
	isect := Intersection{T: 4, S: &s}
	comps := isect.Prepare(ray, []Intersection{isect})
	assert := assert.New(t)
	assert.Equal(isect.T, comps.T)
	assert.Equal(isect.S, comps.S)
	assert.Equal(m.Point4(0, 0, -1), comps.Point)
	assert.Equal(m.Vector4(0, 0, -1), comps.Eye)
	assert.Equal(m.Vector4(0, 0, -1), comps.Normal)
}

func TestHitWhenIntersectionOnOutside(t *testing.T) {
	r := Ray{
		Origin: m.Point4(0, 0, -5),
		Dir:    m.Vector4(0, 0, 1),
	}
	s := NewSphere()
	isect := Intersection{S: &s, T: 4}
	comps := isect.Prepare(r, []Intersection{isect})
	assert.False(t, comps.Inside)
}

func TestHitWhenIntersectionOnInside(t *testing.T) {
	r := Ray{
		Origin: m.Point4(0, 0, 0),
		Dir:    m.Vector4(0, 0, 1),
	}
	s := NewSphere()
	isect := Intersection{S: &s, T: 1}
	comps := isect.Prepare(r, []Intersection{isect})
	assert := assert.New(t)
	assert.True(comps.Inside)
	assert.Equal(m.Point4(0, 0, 1), comps.Point)
	assert.Equal(m.Vector4(0, 0, -1), comps.Eye)
	assert.Equal(m.Vector4(0, 0, -1), comps.Normal)
	assert.True(comps.Inside)
}

func TestHitShouldOffsetThePoint(t *testing.T) {
	r := Ray{Origin: m.Point4(0, 0, -5), Dir: m.Vector4(0, 0, 1)}
	s := NewSphere()
	s.SetTf(m.Trans(0, 0, 1))
	i := Intersection{T: 5, S: &s}
	comps := i.Prepare(r, []Intersection{i})
	assert.Less(t, comps.OverPoint[2], m.EPSILON)
	assert.Greater(t, comps.Point[2], comps.OverPoint[2])
}

func TestShapeDefaultTransformIsIdentity(t *testing.T) {
	s := NewTestShape()
	assert.Equal(t, m.Mat4Ident(), s.GetTf())
}

func TestShapeSetTransform(t *testing.T) {
	s := NewTestShape()
	tf := m.Trans(2, 3, 4)
	s.SetTf(tf)
	assert.Equal(t, tf, s.GetTf())
}

func TestShapeGetMaterial(t *testing.T) {
	s := NewTestShape()
	assert.Equal(t, mtl.DefaultMaterial(), s.GetMat())
}

func TestShapeSetMaterial(t *testing.T) {
	s := NewTestShape()
	mat := mtl.DefaultMaterial()
	mat.Ambient = 1
	s.SetMat(mat)
	assert.Equal(t, mat, s.GetMat())
}

func TestNormalOnTranslatedShape(t *testing.T) {
	s := NewTestShape()
	s.SetTf(m.Trans(0, 1, 0))
	n := s.NormalAt(m.Point4(0, 1.70711, -0.70711))
	expected := m.Vector4(0, 0.70711, -0.70711)
	assert.True(t, expected.ApproxEqual(n))
}

func TestNormalOnTransformedShape(t *testing.T) {
	s := NewTestShape()
	s.SetTf(m.Scale(1, 0.5, 1).Mul(m.RotZ(math.Pi / 5.0)))
	n := s.NormalAt(m.Point4(0, math.Sqrt2/2.0, -math.Sqrt2/2.0))
	expected := m.Vector4(0, 0.97014, -0.24254)
	assert.True(t, expected.ApproxEqual(n))
}

func TestNormalOfPlaneIsConstantEverywhere(t *testing.T) {
	p := NewPlane()
	n1 := p.localNormalAt(m.Point4(0, 0, 0))
	n2 := p.localNormalAt(m.Point4(10, 0, -10))
	n3 := p.localNormalAt(m.Point4(-5, 0, 150))
	assert := assert.New(t)
	assert.Equal(m.Vector4(0, 1, 0), n1)
	assert.Equal(m.Vector4(0, 1, 0), n2)
	assert.Equal(m.Vector4(0, 1, 0), n3)
}

func TestIntersectPlaneAndParallelRay(t *testing.T) {
	p := NewPlane()
	r := Ray{Origin: m.Point4(0, 10, 0), Dir: m.Vector4(0, 0, 1)}
	isects := p.localIntersect(r)
	assert.Empty(t, isects)
}

func TestIntersectPlaneAndCoplanarRay(t *testing.T) {
	p := NewPlane()
	r := Ray{Origin: m.Point4(0, 0, 0), Dir: m.Vector4(0, 0, 1)}
	isects := p.localIntersect(r)
	assert.Empty(t, isects)
}

func TestIntersectPlaneWhenRayFromAbove(t *testing.T) {
	p := NewPlane()
	r := Ray{Origin: m.Point4(0, 1, 0), Dir: m.Vector4(0, -1, 0)}
	isects := p.localIntersect(r)
	assert := assert.New(t)
	assert.Len(isects, 1)
	assert.Equal(1.0, isects[0].T)
	assert.Equal(&p, isects[0].S)
}

func TestIntersectPlaneWhenRayFromBelow(t *testing.T) {
	p := NewPlane()
	r := Ray{Origin: m.Point4(0, -1, 0), Dir: m.Vector4(0, 1, 0)}
	isects := p.localIntersect(r)
	assert := assert.New(t)
	assert.Len(isects, 1)
	assert.Equal(1.0, isects[0].T)
	assert.Equal(&p, isects[0].S)
}

func TestSampleStripesWhenShapeTransformed(t *testing.T) {
	s := NewSphere()
	s.SetTf(m.Scale(2, 2, 2))
	p := pattern.NewStripePattern(color.White, color.Black)
	s.SetPattern(&p)
	actual := s.SamplePatternAt(m.Point4(1.5, 0, 0))
	assert.Equal(t, color.White, actual)
}

func TestStripesWhenPatternTransformed(t *testing.T) {
	s := NewSphere()
	p := pattern.NewStripePattern(color.White, color.Black)
	p.SetTf(m.Scale(2, 2, 2))
	s.SetPattern(&p)
	actual := s.SamplePatternAt(m.Point4(1.5, 0, 0))
	assert.Equal(t, color.White, actual)
}

func TestStripesWhenPatternAndShapeTransformed(t *testing.T) {
	s := NewSphere()
	s.SetTf(m.Scale(2, 2, 2))
	p := pattern.NewStripePattern(color.White, color.Black)
	p.SetTf(m.Scale(2, 2, 2))
	s.SetPattern(&p)
	actual := s.SamplePatternAt(m.Point4(2.5, 0, 0))
	assert.Equal(t, color.White, actual)
}

func TestPatternWithObjectTransformed(t *testing.T) {
	s := NewSphere()
	s.SetTf(m.Scale(2, 2, 2))
	p := pattern.NewTestPattern()
	s.SetPattern(&p)
	actual := s.SamplePatternAt(m.Point4(2, 3, 4))
	assert.Equal(t, m.Vec4{1, 1.5, 2}, actual)
}

func TestPatternWithPatternTransformed(t *testing.T) {
	s := NewSphere()
	p := pattern.NewTestPattern()
	p.SetTf(m.Scale(2, 2, 2))
	s.SetPattern(&p)
	actual := s.SamplePatternAt(m.Point4(2, 3, 4))
	assert.Equal(t, m.Vec4{1, 1.5, 2}, actual)
}

func TestPatternWithShapeAndPatternTransformed(t *testing.T) {
	s := NewSphere()
	s.SetTf(m.Scale(2, 2, 2))
	p := pattern.NewTestPattern()
	p.SetTf(m.Trans(0.5, 1, 1.5))
	s.SetPattern(&p)
	actual := s.SamplePatternAt(m.Point4(2.5, 3, 3.5))
	assert.Equal(t, m.Vec4{0.75, 0.5, 0.25}, actual)
}

func TestPrepareReflectionVector(t *testing.T) {
	p := NewPlane()
	r := Ray{
		Origin: m.Point4(0, 1, -1),
		Dir:    m.Vector4(0, -math.Sqrt2/2.0, math.Sqrt2/2.0),
	}
	isect := Intersection{T: math.Sqrt2, S: &p}
	comps := isect.Prepare(r, []Intersection{isect})
	assert.Equal(t, m.Vector4(0, math.Sqrt2/2.0, math.Sqrt2/2.0), comps.Reflect)
}

func TestCreateGlassSphere(t *testing.T) {
	s := NewGlassSphere()
	assert := assert.New(t)
	assert.Equal(m.Mat4Ident(), s.O.Tf)
	assert.Equal(1.0, s.O.Material.Transparency)
	assert.Equal(1.5, s.O.Material.RefractiveIndex)
}

func TestFindingN1AndN2AtVariousIntersections(t *testing.T) {
	a := NewGlassSphere()
	a.SetTf(m.Scale(2, 2, 2))
	a.O.Material.RefractiveIndex = 1.5

	b := NewGlassSphere()
	b.SetTf(m.Trans(0, 0, -0.25))
	b.O.Material.RefractiveIndex = 2.0

	c := NewGlassSphere()
	c.SetTf(m.Trans(0, 0, 0.25))
	c.O.Material.RefractiveIndex = 2.5

	r := Ray{
		Origin: m.Point4(0, 0, -4),
		Dir:    m.Vector4(0, 0, 1),
	}
	isects := []Intersection{
		{S: &a, T: 2},
		{S: &b, T: 2.75},
		{S: &c, T: 3.25},
		{S: &b, T: 4.75},
		{S: &c, T: 5.25},
		{S: &a, T: 6},
	}

	type testData struct {
		i      int
		n1, n2 float64
	}
	td := []testData{
		{i: 0, n1: 1.0, n2: 1.5},
		{i: 1, n1: 1.5, n2: 2.0},
		{i: 2, n1: 2.0, n2: 2.5},
		{i: 3, n1: 2.5, n2: 2.5},
		{i: 4, n1: 2.5, n2: 1.5},
		{i: 5, n1: 1.5, n2: 1.0},
	}

	assert := assert.New(t)
	for _, d := range td {
		comps := isects[d.i].Prepare(r, isects)
		assert.Equal(comps.N1, d.n1)
		assert.Equal(comps.N2, d.n2)
	}
}

func TestUnderPointIsOffsetBelowTheSurface(t *testing.T) {
	r := Ray{Origin: m.Point4(0, 0, -5), Dir: m.Vector4(0, 0, 1)}
	s := NewGlassSphere()
	s.SetTf(m.Trans(0, 0, 1))
	isect := Intersection{T: 5, S: &s}
	xs := []Intersection{isect}
	comps := isect.Prepare(r, xs)
	assert.Greater(t, comps.UnderPoint[2], m.EPSILON/2.0)
	assert.Less(t, comps.Point[2], comps.UnderPoint[2])
}

func TestSchlickUnderTotalInternalReflection(t *testing.T) {
	shape := NewGlassSphere()
	r := Ray{
		Origin: m.Point4(0, 0, math.Sqrt2/2.0),
		Dir:    m.Vector4(0, 1, 0),
	}
	xs := []Intersection{
		{T: -math.Sqrt2 / 2.0, S: &shape},
		{T: math.Sqrt2 / 2.0, S: &shape},
	}

	comps := xs[1].Prepare(r, xs)
	assert.Equal(t, 1.0, comps.Schlick())
}

func TestSchlickWithPerpendicularViewingAngle(t *testing.T) {
	shape := NewGlassSphere()
	r := Ray{
		Origin: m.Point4(0, 0, 0),
		Dir:    m.Vector4(0, 1, 0),
	}
	xs := []Intersection{
		{T: -1, S: &shape},
		{T: 1, S: &shape},
	}

	comps := xs[1].Prepare(r, xs)
	assert.True(t, m.EqApprox(0.04, comps.Schlick()))
}

func TestSchlickWithSmallAngleAndN2GTN1(t *testing.T) {
	shape := NewGlassSphere()
	r := Ray{
		Origin: m.Point4(0, 0.99, -2),
		Dir:    m.Vector4(0, 0, 1),
	}
	xs := []Intersection{
		{T: 1.8589, S: &shape},
	}

	comps := xs[0].Prepare(r, xs)
	assert.True(t, m.EqApprox(0.48873, comps.Schlick()))
}

func TestRayIntersectsCube(t *testing.T) {
	c := NewCube()

	type testData struct {
		ray    Ray
		t1, t2 float64
	}
	td := []testData{
		{ray: Ray{Origin: m.Point4(5, 0.5, 0), Dir: m.Vector4(-1, 0, 0)}, t1: 4, t2: 6},
		{ray: Ray{Origin: m.Point4(-5, 0.5, 0), Dir: m.Vector4(1, 0, 0)}, t1: 4, t2: 6},
		{ray: Ray{Origin: m.Point4(0.5, 5, 0), Dir: m.Vector4(0, -1, 0)}, t1: 4, t2: 6},
		{ray: Ray{Origin: m.Point4(0.5, -5, 0), Dir: m.Vector4(0, 1, 0)}, t1: 4, t2: 6},
		{ray: Ray{Origin: m.Point4(0.5, 0, 5), Dir: m.Vector4(0, 0, -1)}, t1: 4, t2: 6},
		{ray: Ray{Origin: m.Point4(0.5, 0, -5), Dir: m.Vector4(0, 0, 1)}, t1: 4, t2: 6},
		{ray: Ray{Origin: m.Point4(0, 0.5, 0), Dir: m.Vector4(0, 0, 1)}, t1: -1, t2: 1},
	}
	assert := assert.New(t)
	for _, d := range td {
		xs := c.localIntersect(d.ray)
		assert.Equal(2, len(xs))
		assert.Equal(d.t1, xs[0].T)
		assert.Equal(d.t2, xs[1].T)
	}
}

func TestRayMissesCube(t *testing.T) {
	c := NewCube()

	td := []Ray{
		{Origin: m.Point4(-2, 0, 0), Dir: m.Vector4(0.2673, 0.5345, 0.8018)},
		{Origin: m.Point4(0, -2, 0), Dir: m.Vector4(0.8018, 0.2673, 0.5345)},
		{Origin: m.Point4(0, 0, -2), Dir: m.Vector4(0.5345, 0.8018, 0.2673)},
		{Origin: m.Point4(2, 0, 2), Dir: m.Vector4(0, 0, -1)},
		{Origin: m.Point4(0, 2, 2), Dir: m.Vector4(0, -1, 0)},
		{Origin: m.Point4(2, 2, 0), Dir: m.Vector4(-1, 0, 0)},
	}

	assert := assert.New(t)

	for _, d := range td {
		xs := c.localIntersect(d)
		assert.Len(xs, 0)
	}
}

func TestCubeNormalAt(t *testing.T) {
	c := NewCube()
	type testData struct {
		point  m.Vec4
		normal m.Vec4
	}
	td := []testData{
		{point: m.Point4(1, 0.5, -0.8), normal: m.Vector4(1, 0, 0)},
		{point: m.Point4(-1, -0.2, 0.9), normal: m.Vector4(-1, 0, 0)},
		{point: m.Point4(-0.4, 1, -0.1), normal: m.Vector4(0, 1, 0)},
		{point: m.Point4(0.3, -1, -0.7), normal: m.Vector4(0, -1, 0)},
		{point: m.Point4(-0.6, 0.3, 1), normal: m.Vector4(0, 0, 1)},
		{point: m.Point4(0.4, 0.4, -1), normal: m.Vector4(0, 0, -1)},
		{point: m.Point4(1, 1, 1), normal: m.Vector4(1, 0, 0)},
		{point: m.Point4(-1, -1, -1), normal: m.Vector4(-1, 0, 0)},
	}

	assert := assert.New(t)
	for _, d := range td {
		normal := c.localNormalAt(d.point)
		assert.Equal(d.normal, normal)
	}
}

func TestRayMissesCylinder(t *testing.T) {
	c := NewCylinder()

	testData := []Ray{
		{Origin: m.Point4(1, 0, 0), Dir: m.Vector4(0, 1, 0)},
		{Origin: m.Point4(0, 0, 0), Dir: m.Vector4(0, 1, 0)},
		{Origin: m.Point4(0, 0, -5), Dir: m.Vector4(1, 1, 1)},
	}

	assert := assert.New(t)
	for _, d := range testData {
		r := d
		r.Dir = r.Dir.Normalize()
		xs := c.localIntersect(r)
		assert.Len(xs, 0)
	}
}

func TestRayHitsCylinder(t *testing.T) {
	c := NewCylinder()

	type testData struct {
		origin, dir m.Vec4
		t0, t1      float64
	}
	td := []testData{
		{origin: m.Point4(1, 0, -5), dir: m.Vector4(0, 0, 1), t0: 5, t1: 5},
		{origin: m.Point4(0, 0, -5), dir: m.Vector4(0, 0, 1), t0: 4, t1: 6},
		{origin: m.Point4(0.5, 0, -5), dir: m.Vector4(0.1, 1, 1), t0: 6.80798, t1: 7.08872},
	}

	assert := assert.New(t)

	for _, d := range td {
		r := Ray{
			Origin: d.origin,
			Dir:    d.dir.Normalize(),
		}
		xs := c.localIntersect(r)
		assert.Len(xs, 2)
		assert.True(m.EqApprox(d.t0, xs[0].T))
		assert.True(m.EqApprox(d.t1, xs[1].T))
	}
}

func TestNormalOnCylinder(t *testing.T) {
	c := NewCylinder()

	type testData struct{ point, normal m.Vec4 }
	td := []testData{
		{point: m.Point4(1, 0, 0), normal: m.Vector4(1, 0, 0)},
		{point: m.Point4(0, 5, -1), normal: m.Vector4(0, 0, -1)},
		{point: m.Point4(0, -2, 1), normal: m.Vector4(0, 0, 1)},
		{point: m.Point4(-1, 1, 0), normal: m.Vector4(-1, 0, 0)},
	}

	assert := assert.New(t)

	for _, d := range td {
		assert.Equal(d.normal, c.localNormalAt(d.point))
	}
}

func TestDefaultMinAndMaxForCylinder(t *testing.T) {
	c := NewCylinder()
	assert.Equal(t, math.Inf(-1), c.Min)
	assert.Equal(t, math.Inf(1), c.Max)
}

func TestIntersectingAConstrainedCylinder(t *testing.T) {
	c := NewCylinder()
	c.Min = 1.0
	c.Max = 2.0

	type testdata struct {
		point, dir m.Vec4
		count      int
	}
	td := []testdata{
		{point: m.Point4(0, 1.5, 0), dir: m.Vector4(0.1, 1, 0), count: 0},
		{point: m.Point4(0, 3, -5), dir: m.Vector4(0, 0, 1), count: 0},
		{point: m.Point4(0, 0, -5), dir: m.Vector4(0, 0, 1), count: 0},
		{point: m.Point4(0, 2, -5), dir: m.Vector4(0, 0, 1), count: 0},
		{point: m.Point4(0, 1, -5), dir: m.Vector4(0, 0, 1), count: 0},
		{point: m.Point4(0, 1.5, -2), dir: m.Vector4(0, 0, 1), count: 2},
	}

	assert := assert.New(t)

	for _, d := range td {
		r := Ray{
			Origin: d.point,
			Dir:    d.dir.Normalize(),
		}
		xs := c.localIntersect(r)
		assert.Len(xs, d.count)
	}
}

func TestDefaultClosedValueForCylinder(t *testing.T) {
	c := NewCylinder()
	assert.False(t, c.Closed)
}

func TestIntersectingCapsOfCylinders(t *testing.T) {
	c := NewCylinder()
	c.Min = 1
	c.Max = 2
	c.Closed = true

	type testdata struct {
		point, dir m.Vec4
		count      int
	}
	td := []testdata{
		{point: m.Point4(0, 3, 0), dir: m.Vector4(0, -1, 0), count: 2},
		{point: m.Point4(0, 3, -2), dir: m.Vector4(0, -1, 2), count: 2},
		{point: m.Point4(0, 4, -2), dir: m.Vector4(0, -1, 1), count: 2},
		{point: m.Point4(0, 0, -2), dir: m.Vector4(0, 1, 2), count: 2},
		{point: m.Point4(0, -1, -2), dir: m.Vector4(0, 1, 1), count: 2},
	}

	assert := assert.New(t)

	for _, d := range td {
		r := Ray{
			Origin: d.point,
			Dir:    d.dir.Normalize(),
		}
		xs := c.localIntersect(r)
		assert.Len(xs, d.count)
	}
}

func TestNormalOnCylinderEndCaps(t *testing.T) {
	c := NewCylinder()
	c.Min = 1.0
	c.Max = 2.0
	c.Closed = true

	type testdata struct{ point, normal m.Vec4 }
	td := []testdata{
		{point: m.Point4(0, 1, 0), normal: m.Vector4(0, -1, 0)},
		{point: m.Point4(0.5, 1, 0), normal: m.Vector4(0, -1, 0)},
		{point: m.Point4(0, 1, 0.5), normal: m.Vector4(0, -1, 0)},
		{point: m.Point4(0, 2, 0), normal: m.Vector4(0, 1, 0)},
		{point: m.Point4(0.5, 2, 0), normal: m.Vector4(0, 1, 0)},
		{point: m.Point4(0, 2, 0.5), normal: m.Vector4(0, 1, 0)},
	}

	assert := assert.New(t)

	for _, d := range td {
		assert.Equal(d.normal, c.localNormalAt(d.point))
	}
}

func TestIntersectConeWithRay(t *testing.T) {
	c := NewCone()

	type testdata struct {
		origin, dir m.Vec4
		t0, t1      float64
	}

	td := []testdata{
		{origin: m.Point4(0, 0, -5), dir: m.Vector4(0, 0, 1), t0: 5, t1: 5},
		{origin: m.Point4(0, 0, -5), dir: m.Vector4(1, 1, 1), t0: 8.66025, t1: 8.66025},
		{origin: m.Point4(1, 1, -5), dir: m.Vector4(-0.5, -1, 1), t0: 4.55006, t1: 49.44994},
	}

	assert := assert.New(t)

	for _, d := range td {
		r := Ray{
			Origin: d.origin,
			Dir:    d.dir.Normalize(),
		}
		xs := c.localIntersect(r)
		assert.Len(xs, 2)
		assert.True(m.EqApprox(d.t0, xs[0].T))
		assert.True(m.EqApprox(d.t1, xs[1].T))
	}
}

func TestIntersectConeWithARayParallelToOneOfItsHalves(t *testing.T) {
	c := NewCone()
	r := Ray{
		Origin: m.Point4(0, 0, -1),
		Dir:    m.Vector4(0, 1, 1).Normalize(),
	}
	xs := c.localIntersect(r)
	assert.Len(t, xs, 1)
	assert.True(t, m.EqApprox(xs[0].T, 0.35355))
}

func TestIntersectConeEndCaps(t *testing.T) {
	c := NewCone()
	c.Min = -0.5
	c.Max = 0.5
	c.Closed = true

	type testdata struct {
		origin, dir m.Vec4
		count       int
	}
	td := []testdata{
		{origin: m.Point4(0, 0, -5), dir: m.Vector4(0, 1, 0), count: 0},
		{origin: m.Point4(0, 0, -0.25), dir: m.Vector4(0, 1, 1), count: 2},
		{origin: m.Point4(0, 0, -0.25), dir: m.Vector4(0, 1, 0), count: 4},
	}

	assert := assert.New(t)

	for _, d := range td {
		r := Ray{
			Origin: d.origin,
			Dir:    d.dir.Normalize(),
		}
		xs := c.localIntersect(r)
		assert.Len(xs, d.count)
	}
}

func TestConeNormal(t *testing.T) {
	c := NewCone()

	type testdata struct{ point, normal m.Vec4 }
	td := []testdata{
		{point: m.Point4(0, 0, 0), normal: m.Vector4(0, 0, 0)},
		{point: m.Point4(1, 1, 1), normal: m.Vector4(1, -math.Sqrt2, 1)},
		{point: m.Point4(-1, -1, 0), normal: m.Vector4(-1, 1, 0)},
	}

	assert := assert.New(t)

	for _, d := range td {
		assert.Equal(d.normal, c.localNormalAt(d.point))
	}
}
