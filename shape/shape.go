package shape

import (
	"math"

	"github.com/elliotchance/pie/v2"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/pattern"
)

type (
	Ray struct {
		Origin m.Vec4
		Dir    m.Vec4
	}

	Object struct {
		Material mtl.Material
		Tf       m.Mat4
		InvTf    m.Mat4
		SavedRay Ray
	}

	Sphere struct {
		O Object
	}

	Plane struct {
		O Object
	}

	Intersection struct {
		S Shape
		T float64
	}

	Shape interface {
		GetTf() m.Mat4
		SetTf(m.Mat4)
		GetMat() mtl.Material
		SetMat(mtl.Material)
		Intersect(Ray) []Intersection
		GetSavedRay() Ray
		NormalAt(m.Vec4) m.Vec4
		SamplePatternAt(m.Vec4) m.Vec4
		SetPattern(pattern.Pattern)
	}

	IntersectionComps struct {
		S          Shape
		T          float64
		Point      m.Vec4
		Eye        m.Vec4
		Normal     m.Vec4
		Inside     bool
		OverPoint  m.Vec4
		UnderPoint m.Vec4
		Reflect    m.Vec4
		N1, N2     float64
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
func NewSphere() Sphere {
	o := Object{
		Tf:       m.Mat4Ident(),
		InvTf:    m.Mat4Ident(),
		Material: mtl.DefaultMaterial(),
	}
	return Sphere{O: o}
}

func NewGlassSphere() Sphere {
	o := Object{
		Tf:    m.Mat4Ident(),
		InvTf: m.Mat4Ident(),
		Material: mtl.Material{
			Transparency:    1.0,
			RefractiveIndex: 1.5,
		},
	}
	return Sphere{O: o}
}

func (s *Sphere) Intersect(r Ray) []Intersection {
	localRay := r.Transform(s.O.InvTf)
	return s.localIntersect(localRay)
}

func (s *Sphere) localIntersect(r Ray) []Intersection {
	s.O.SavedRay = r
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
			{S: s, T: t1},
			{S: s, T: t2},
		}
	}
	return []Intersection{
		{S: s, T: t2},
		{S: s, T: t1},
	}
}

func (s *Sphere) NormalAt(p m.Vec4) m.Vec4 {
	localPoint := s.O.InvTf.MulVec(p)
	localNormal := s.localNormalAt(localPoint)
	nWorld := s.O.InvTf.Tpose().MulVec(localNormal)
	nWorld[3] = 0
	return nWorld.Normalize()
}

func (s *Sphere) localNormalAt(p m.Vec4) m.Vec4 {
	return m.Vector4(p[0], p[1], p[2])
}

func (s *Sphere) GetTf() m.Mat4 {
	return s.O.Tf
}

func (s *Sphere) SetTf(tf m.Mat4) {
	s.O.Tf = tf
	s.O.InvTf = tf.Inv()
}

func (s *Sphere) GetMat() mtl.Material {
	return s.O.Material
}

func (s *Sphere) SetMat(mat mtl.Material) {
	s.O.Material = mat
}

func (s *Sphere) GetSavedRay() Ray {
	return s.O.SavedRay
}

func (s *Sphere) SamplePatternAt(point m.Vec4) m.Vec4 {
	objPoint := s.O.InvTf.MulVec(point)
	pattern := s.O.Material.Pattern
	patternPoint := pattern.GetInvTf().MulVec(objPoint)
	return pattern.SampleAt(patternPoint)
}

func (s *Sphere) SetPattern(p pattern.Pattern) {
	s.O.Material.Pattern = p
}

// ////////////// PLANE ////////////////
func NewPlane() Plane {
	return Plane{O: defaultObject()}
}

func (p *Plane) localNormalAt(_ m.Vec4) m.Vec4 {
	return m.Vector4(0, 1, 0)
}

func (p *Plane) localIntersect(ray Ray) []Intersection {
	if math.Abs(ray.Dir[1]) < m.EPSILON {
		return []Intersection{}
	}

	t := -ray.Origin[1] / ray.Dir[1]
	return []Intersection{{T: t, S: p}}
}

func (p *Plane) GetMat() mtl.Material {
	return p.O.Material
}

func (p *Plane) SetMat(mat mtl.Material) {
	p.O.Material = mat
}

func (p *Plane) GetTf() m.Mat4 {
	return p.O.Tf
}

func (p *Plane) SetTf(tf m.Mat4) {
	p.O.Tf = tf
	p.O.InvTf = tf.Inv()
}

func (p *Plane) GetSavedRay() Ray {
	return p.O.SavedRay
}

func (p *Plane) Intersect(ray Ray) []Intersection {
	localRay := ray.Transform(p.O.InvTf)
	return p.localIntersect(localRay)
}

func (p *Plane) NormalAt(point m.Vec4) m.Vec4 {
	localPoint := p.O.InvTf.MulVec(point)
	localNormal := p.localNormalAt(localPoint)
	nWorld := p.O.InvTf.Tpose().MulVec(localNormal)
	nWorld[3] = 0
	return nWorld.Normalize()
}

func (p *Plane) SamplePatternAt(point m.Vec4) m.Vec4 {
	objPoint := p.O.InvTf.MulVec(point)
	pattern := p.O.Material.Pattern
	patternPoint := pattern.GetInvTf().MulVec(objPoint)
	return pattern.SampleAt(patternPoint)
}

func (p *Plane) SetPattern(pat pattern.Pattern) {
	p.O.Material.Pattern = pat
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
			res = Intersection{S: isect.S, T: isect.T}
			isHit = true
		}
	}
	return res, isHit
}

func (i Intersection) Prepare(ray Ray, isects []Intersection) IntersectionComps {
	pos := ray.Pos(i.T)
	normal := i.S.NormalAt(pos)
	eye := ray.Dir.Negate()
	inside := false
	if normal.Dot(eye) < 0.0 {
		inside = true
		normal = normal.Negate()
	}
	op := pos.Add(normal.Mul(m.EPSILON))
	up := pos.Sub(normal.Mul(m.EPSILON))

	n1, n2 := computeN(i, isects)

	return IntersectionComps{
		T:          i.T,
		S:          i.S,
		Point:      pos,
		Eye:        eye,
		Normal:     normal,
		Inside:     inside,
		OverPoint:  op,
		UnderPoint: up,
		Reflect:    ray.Dir.Reflect(normal),
		N1:         n1,
		N2:         n2,
	}
}

func computeN(isect Intersection, isects []Intersection) (float64, float64) {
	var n1, n2 float64
	containers := []Shape{}
	for _, i := range isects {
		if i == isect {
			if len(containers) == 0 {
				n1 = 1.0
			} else {
				n1 = containers[len(containers)-1].GetMat().RefractiveIndex
			}
		}
		if pie.Contains(containers, i.S) {
			containers = pie.FilterNot(containers, func(shape Shape) bool { return shape == i.S })
		} else {
			containers = append(containers, i.S)
		}

		if i == isect {
			if len(containers) == 0 {
				n2 = 1.0
			} else {
				n2 = containers[len(containers)-1].GetMat().RefractiveIndex
			}
			break
		}
	}
	return n1, n2
}

// ////////////// SHAPE ////////////////
func NewTestShape() Shape {
	s := Sphere{O: defaultObject()}
	return &s
}

func defaultObject() Object {
	return Object{
		Tf:       m.Mat4Ident(),
		InvTf:    m.Mat4Ident(),
		Material: mtl.DefaultMaterial(),
	}
}

// ////////////// INTERSECTION COMPS ////////////////
func (comps IntersectionComps) Schlick() float64 {
	cos := comps.Eye.Dot(comps.Normal)

	if comps.N1 > comps.N2 {
		n := comps.N1 / comps.N2
		sin2t := n * n * (1.0 - cos*cos)
		if sin2t > 1.0 {
			return 1.0
		}
		cosT := math.Sqrt(1.0 - sin2t)
		cos = cosT
	}
	f := ((comps.N1 - comps.N2) / (comps.N1 + comps.N2))
	r0 := f * f
	return r0 + (1.0-r0)*math.Pow(1.0-cos, 5)
}
