package ray

import (
	"math"

	"github.com/google/uuid"
	m "roytracer/math"
	"roytracer/mtl"
)

type (
	Ray struct {
		Origin m.Vec4
		Dir    m.Vec4
	}

	Sphere struct {
		Id       uuid.UUID
		Tf       m.Mat4
		Material mtl.Material
	}

	Intersection struct {
		T  float64
		Id uuid.UUID
	}
)

// ////////////// RAY ////////////////
func (r Ray) Pos(t float64) m.Vec4 {
	return r.Origin.Add(r.Dir.Mul(t))
}

func (r Ray) Transform(tf m.Mat4) Ray {
	origin := tf.MulVec(r.Origin)
	dir := tf.MulVec(r.Dir)
	return Ray{Origin: origin, Dir: dir}
}

// ////////////// SPHERE ////////////////
func (s Sphere) Intersect(r Ray) []Intersection {
	r = r.Transform(s.Tf.Inv())
	sphereToRay := r.Origin.Sub(m.Point4(0, 0, 0))
	a := r.Dir.Dot(r.Dir)
	b := 2 * r.Dir.Dot(sphereToRay)
	c := sphereToRay.Dot(sphereToRay) - 1
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return []Intersection{}
	}
	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)
	if t1 < t2 {
		return []Intersection{
			{Id: s.Id, T: t1},
			{Id: s.Id, T: t2},
		}
	}
	return []Intersection{
		{Id: s.Id, T: t2},
		{Id: s.Id, T: t1},
	}
}

func NewSphere() Sphere {
	return Sphere{
		Tf: m.Mat4Ident(),
		Id: uuid.New(),
		Material: mtl.DefaultMaterial(),
	}
}

func (s Sphere) NormalAt(p m.Vec4) m.Vec4 {
	invTf := s.Tf.Inv()
	pObj := invTf.MulVec(p)
	nObj := pObj.Sub(m.Point4(0, 0, 0))
	nWorld := invTf.Tpose().MulVec(nObj)
	nWorld[3] = 0
	return nWorld.Normalize()
}

// ////////////// INTERSECTIONS ////////////////
func Intersections(isects ...Intersection) []Intersection {
	return isects
}

func Hit(isects []Intersection) (Intersection, bool) {
	res := Intersection{T: math.MaxFloat64}
	isHit := false
	for _, isect := range isects {
		if isect.T <= res.T && isect.T >= 0 {
			res = Intersection{Id: isect.Id, T: isect.T}
			isHit = true
		}
	}
	return res, isHit
}
