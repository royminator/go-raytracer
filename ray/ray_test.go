package ray

import (
	"math"
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
	assert.Equal(2, len(isect))
	assert.Equal(4.0, isect[0].T)
	assert.Equal(6.0, isect[1].T)
}

func TestSphereIntersectionAtTangent(t *testing.T) {
	ray := Ray{m.Point4(0, 1, -5), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect))
	assert.Equal(5.0, isect[0].T)
	assert.Equal(5.0, isect[1].T)
}

func TestSphereIntersectionRayMisses(t *testing.T) {
	ray := Ray{m.Point4(0, 2, -5), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert.Equal(t, 0, len(isect))
}

func TestSphereIntersectionRayBehindSphere(t *testing.T) {
	ray := Ray{m.Point4(0, 0, 5), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect))
	assert.Equal(-6.0, isect[0].T)
	assert.Equal(-4.0, isect[1].T)
}

func TestSphereIntersectRayInsideSphere(t *testing.T) {
	ray := Ray{m.Point4(0, 0, 0), m.Vector4(0, 0, 1)}
	isect := NewSphere().Intersect(ray)
	assert := assert.New(t)
	assert.Equal(2, len(isect))
	assert.Equal(-1.0, isect[0].T)
	assert.Equal(1.0, isect[1].T)
}

func TestIntersectStoresObjectId(t *testing.T) {
	s := NewSphere()
	isect := Intersection{Id: s.Id, T: 3.5}
	assert.Equal(t, s.Id, isect.Id)
	assert.Equal(t, 3.5, isect.T)
}

func TestInterectionAggregateTValues(t *testing.T) {
	s := NewSphere()
	i1 := Intersection{Id: s.Id, T: 3.5}
	i2 := Intersection{Id: s.Id, T: 5.5}
	is := Intersections(i1, i2)
	assert.Equal(t, 2, len(is))
	assert.Equal(t, 3.5, is[0].T)
	assert.Equal(t, 5.5, is[1].T)
}

func TestInsersectRaySetsObjectId(t *testing.T) {
	s := NewSphere()
	ray := Ray{m.Point4(0, 0, -5), m.Vector4(0, 0, 1)}
	isects := s.Intersect(ray)
	assert.Equal(t, 2, len(isects))
	assert.Equal(t, s.Id, isects[0].Id)
	assert.Equal(t, s.Id, isects[1].Id)
}

func TestHitAllIntersectionGT0(t *testing.T) {
	type testData struct {
		isects []Intersection
		any bool
		res Intersection
	}

	s := NewSphere()
	td := []testData{
		{ 
			isects: []Intersection{
				{ T: 1.0 }, { T: 2 },
			},
			any: true,
			res: Intersection{Id: s.Id, T: 1},
		},
		{ 
			isects: []Intersection{
				{ T: -1 }, { T: 1 },
			},
			any: true,
			res: Intersection{Id: s.Id, T: 1},
		},
		{ 
			isects: []Intersection{
				{ T: -2 }, { T: -1 },
			},
			any: false,
			res: Intersection{},
		},
		{ 
			isects: []Intersection{
				{ T: 5 }, { T: 7 }, { T: -3 }, { T: 2 },
			},
			any: true,
			res: Intersection{Id: s.Id, T: 2.0},
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
	tf := m.Trans(m.Vec3{3, 4, 5})
	expected := Ray{m.Point4(4, 6, 8), m.Vector4(0, 1, 0)}
	actual := ray.Transform(tf)
	assert.Equal(t, expected, actual)
}

func TestScaleRay(t *testing.T) {
	ray := Ray{m.Point4(1, 2, 3), m.Vector4(0, 1, 0)}
	tf := m.Scale(m.Vec3{2, 3, 4})
	expected := Ray{m.Point4(2, 6, 12), m.Vector4(0, 3, 0)}
	actual := ray.Transform(tf)
	assert.Equal(t, expected, actual)
}

func TestSpheresDefaultTransformIsIdentity(t *testing.T) {
	assert.Equal(t, m.Mat4Ident(), NewSphere().Tf)
}

func TestSphereSetTransform(t *testing.T) {
	s := NewSphere()
	expected := m.Trans(m.Vec3{2, 3, 4})
	s.Tf = expected
	assert.Equal(t, expected, s.Tf)
}

func TestIntersectScaledSphereWithRay(t *testing.T) {
	ray := Ray{m.Point4(0, 0, -5), m.Vector4(0, 0, 1)}
	s := NewSphere()
	s.Tf = m.Scale(m.Vec3{2, 2, 2})
	isects := s.Intersect(ray)
	assert.Equal(t, 2, len(isects))
	assert.Equal(t, 3.0, isects[0].T)
	assert.Equal(t, 7.0, isects[1].T)
}

func TestIntersectTranslatedSphereWithRay(t *testing.T) {
	ray := Ray{m.Point4(0, 0, -5), m.Vector4(0, 0, 1)}
	s := NewSphere()
	s.Tf = m.Trans(m.Vec3{5, 0, 0})
	isects := s.Intersect(ray)
	assert.Equal(t, 0, len(isects))
}

func TestSphereNormal(t *testing.T) {
	type testData struct {p m.Vec4; res m.Vec4}
	k := math.Sqrt(3.0)/3.0
	td := []testData{
		{ m.Point4(1, 0, 0), m.Vector4(1, 0, 0)},
		{ m.Point4(0, 1, 0), m.Vector4(0, 1, 0)},
		{ m.Point4(k, k, k), m.Vector4(k, k, k)},
	}

	s := NewSphere()

	assert := assert.New(t)
	for _, d := range td {
		assert.Equal(d.res, s.NormalAt(d.p))
	}
}

func TestSphereNormalIsNormalized(t *testing.T) {
	s := NewSphere()
	k := math.Sqrt(3.0)/3.0
	n := s.NormalAt(m.Point4(k, k, k))
	assert.Equal(t, n, n.Normalize())
}

func TestSphereNormalOnTranslatedSphere(t *testing.T) {
	s := NewSphere()
	s.Tf = m.Trans(m.Vec3{0, 1, 0})
	expected := m.Vector4(0, 0.70711, -0.70711)
	assert.True(t, expected.ApproxEqual(s.NormalAt(m.Point4(0, 1.70711, -0.70711))))
}

func TestSphereNormalOnTransformedSphere(t *testing.T) {
	s := NewSphere()
	s.Tf = m.Scale(m.Vec3{1, 0.5, 1}).Mul(m.RotZ(math.Pi/5.0))
	expected := m.Vector4(0, 0.97014, -0.24254)
	assert.True(t, expected.ApproxEqual(s.NormalAt(m.Point4(0, math.Sqrt2/2.0, -math.Sqrt2/2.0))))
}
