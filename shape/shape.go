package shape

import (
	"math"

	m "roytracer/math"
	"roytracer/mtl"
)

type (
	Ray struct {
		Origin m.Vec4
		Dir    m.Vec4
	}

	Sphere struct {
		Tf       m.Mat4
		InvTf    m.Mat4
		Material mtl.Material
	}

	Transform struct {
		m.Mat4
	}

	Intersection struct {
		O *Sphere
		T float64
	}

	IntersectionComps struct {
		O      *Sphere
		T      float64
		Point  m.Vec4
		Eye    m.Vec4
		Normal m.Vec4
		Inside bool
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
func (s *Sphere) SetTf(tf m.Mat4) {
	s.Tf = tf
	s.InvTf = tf.Inv()
}

// ////////////// SPHERE ////////////////
func NewSphere() Sphere {
	return Sphere{
		Tf:       m.Mat4Ident(),
		InvTf:    m.Mat4Ident(),
		Material: mtl.DefaultMaterial(),
	}
}

func (s *Sphere) Intersect(r Ray) []Intersection {
	r = r.Transform(s.InvTf)
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
			{O: s, T: t1},
			{O: s, T: t2},
		}
	}
	return []Intersection{
		{O: s, T: t2},
		{O: s, T: t1},
	}
}

func (s *Sphere) NormalAt(p m.Vec4) m.Vec4 {
	pObj := s.InvTf.MulVec(p)
	nObj := pObj.Sub(m.Point4(0, 0, 0))
	nWorld := s.InvTf.Tpose().MulVec(nObj)
	nWorld[3] = 0
	return nWorld.Normalize()
}

func (s *Sphere) SetMaterial(mat mtl.Material) {
	s.Material = mat
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
			res = Intersection{O: isect.O, T: isect.T}
			isHit = true
		}
	}
	return res, isHit
}

func (i Intersection) Prepare(ray Ray) IntersectionComps {
	pos := ray.Pos(i.T)
	normal := i.O.NormalAt(pos)
	eye := ray.Dir.Negate()
	inside := false
	if normal.Dot(eye) < 0.0 {
		inside = true
		normal = normal.Negate()
	}

	return IntersectionComps{
		T:      i.T,
		O:      i.O,
		Point:  pos,
		Eye:    eye,
		Normal: normal,
		Inside: inside,
	}
}
