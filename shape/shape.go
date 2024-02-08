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

	Cube struct {
		O Object
	}

	Cylinder struct {
		O        Object
		Min, Max float64
		Closed   bool
	}

	Cone struct {
		O        Object
		Min, Max float64
		Closed   bool
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
		Intersect(Ray) ([]Intersection, int)
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
	res := r.Origin
	res.Add(r.Dir.Mul(t))
	return res
}

func (r *Ray) Transform(tf m.Mat4) {
	r.Origin.MulMat(tf)
	r.Dir.MulMat(tf)
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

func (s *Sphere) Intersect(r Ray) ([]Intersection, int) {
	r.Transform(s.O.InvTf)
	return s.localIntersect(r)
}

func (s *Sphere) localIntersect(r Ray) ([]Intersection, int) {
	s.O.SavedRay = r
	sphereToRay := r.Origin
	sphereToRay.Sub(m.Point4(0, 0, 0))
	a := r.Dir.Dot(r.Dir)
	b := 2 * r.Dir.Dot(sphereToRay)
	c := sphereToRay.Dot(sphereToRay) - 1
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return []Intersection{}, 0
	}
	sqrtDiscriminant := math.Sqrt(discriminant)
	twoA := 2 * a
	t1 := (-b - sqrtDiscriminant) / twoA
	t2 := (-b + sqrtDiscriminant) / twoA
	if t1 < t2 {
		return []Intersection{
			{S: s, T: t1},
			{S: s, T: t2},
		}, 2
	}
	return []Intersection{
		{S: s, T: t2},
		{S: s, T: t1},
	}, 2
}

func (s *Sphere) NormalAt(p m.Vec4) m.Vec4 {
	p.MulMat(s.O.InvTf)
	nWorld := s.localNormalAt(p)
	nWorld.MulMat(s.O.InvTf.Tpose())
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

func (p *Plane) localIntersect(ray Ray) ([]Intersection, int) {
	if math.Abs(ray.Dir[1]) < m.EPSILON {
		return []Intersection{}, 0
	}

	t := -ray.Origin[1] / ray.Dir[1]
	return []Intersection{{T: t, S: p}}, 1
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

func (p *Plane) Intersect(ray Ray) ([]Intersection, int) {
	ray.Transform(p.O.InvTf)
	return p.localIntersect(ray)
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

// ////////////// CUBE ////////////////
func NewCube() Cube {
	return Cube{
		O: defaultObject(),
	}
}

func (c *Cube) GetTf() m.Mat4 {
	return c.O.Tf
}

func (c *Cube) SetTf(tf m.Mat4) {
	c.O.Tf = tf
	c.O.InvTf = tf.Inv()
}

func (c *Cube) GetMat() mtl.Material {
	return c.O.Material
}

func (c *Cube) SetMat(mtl mtl.Material) {
	c.O.Material = mtl
}

func (c *Cube) Intersect(ray Ray) ([]Intersection, int) {
	ray.Transform(c.O.InvTf)
	return c.localIntersect(ray)
}

func (c *Cube) localIntersect(ray Ray) ([]Intersection, int) {
	xtmin, xtmax := checkAxis(ray.Origin[0], ray.Dir[0])
	ytmin, ytmax := checkAxis(ray.Origin[1], ray.Dir[1])
	ztmin, ztmax := checkAxis(ray.Origin[2], ray.Dir[2])

	tmin := math.Max(math.Max(xtmin, ytmin), ztmin)
	tmax := math.Min(math.Min(xtmax, ytmax), ztmax)

	if tmin > tmax {
		return []Intersection{}, 0
	}

	return []Intersection{{T: tmin, S: c}, {T: tmax, S: c}}, 2
}

func checkAxis(origin float64, dir float64) (float64, float64) {
	tmin_numerator := -1 - origin
	tmax_numerator := 1 - origin

	inf := math.Inf(1)
	tmin := tmin_numerator * inf
	tmax := tmax_numerator * inf

	if math.Abs(dir) >= m.EPSILON {
		tmin = tmin_numerator / dir
		tmax = tmax_numerator / dir
	}

	if tmin > tmax {
		return tmax, tmin
	}
	return tmin, tmax
}

func (c *Cube) GetSavedRay() Ray {
	return c.O.SavedRay
}

func (c *Cube) NormalAt(point m.Vec4) m.Vec4 {
	localPoint := c.O.InvTf.MulVec(point)
	localNormal := c.localNormalAt(localPoint)
	nWorld := c.O.InvTf.Tpose().MulVec(localNormal)
	nWorld[3] = 0
	return nWorld.Normalize()
}

func (c *Cube) localNormalAt(point m.Vec4) m.Vec4 {
	maxc := math.Max(math.Max(math.Abs(point[0]), math.Abs(point[1])), math.Abs(point[2]))
	if maxc == math.Abs(point[0]) {
		return m.Vector4(point[0], 0, 0)
	}
	if maxc == math.Abs(point[1]) {
		return m.Vector4(0, point[1], 0)
	}
	return m.Vector4(0, 0, point[2])
}

func (c *Cube) SamplePatternAt(point m.Vec4) m.Vec4 {
	objPoint := c.O.InvTf.MulVec(point)
	pattern := c.O.Material.Pattern
	patternPoint := pattern.GetInvTf().MulVec(objPoint)
	return pattern.SampleAt(patternPoint)
}

func (c *Cube) SetPattern(p pattern.Pattern) {
	c.O.Material.Pattern = p
}

// ////////////// CYLIDER ////////////////
func NewCylinder() Cylinder {
	return Cylinder{
		O:      defaultObject(),
		Min:    math.Inf(-1),
		Max:    math.Inf(1),
		Closed: false,
	}
}

func (c *Cylinder) GetTf() m.Mat4 {
	return c.O.Tf
}

func (c *Cylinder) SetTf(tf m.Mat4) {
	c.O.Tf = tf
	c.O.InvTf = tf.Inv()
}

func (c *Cylinder) GetMat() mtl.Material {
	return c.O.Material
}

func (c *Cylinder) SetMat(mtl mtl.Material) {
	c.O.Material = mtl
}

func (c *Cylinder) Intersect(ray Ray) ([]Intersection, int) {
	ray.Transform(c.O.InvTf)
	return c.localIntersect(ray)
}

func (cyl *Cylinder) localIntersect(ray Ray) ([]Intersection, int) {
	a := ray.Dir[0]*ray.Dir[0] + ray.Dir[2]*ray.Dir[2]

	var xs []Intersection

	if !m.EqApprox(a, 0.0) {
		cyl.intersectWalls(ray, &xs, a)
	}

	cyl.intersectCaps(ray, &xs)
	return xs, len(xs)
}

func (cyl *Cylinder) intersectWalls(ray Ray, xs *[]Intersection, a float64) {
	b := 2*ray.Origin[0]*ray.Dir[0] + 2*ray.Origin[2]*ray.Dir[2]
	c := ray.Origin[0]*ray.Origin[0] + ray.Origin[2]*ray.Origin[2] - 1.0

	disc := b*b - 4*a*c

	if disc < 0.0 {
		return
	}

	t0 := (-b - math.Sqrt(disc)) / (2 * a)
	t1 := (-b + math.Sqrt(disc)) / (2 * a)

	if t0 > t1 {
		t0, t1 = t1, t0
	}

	y0 := ray.Origin[1] + t0*ray.Dir[1]
	if cyl.Min < y0 && y0 < cyl.Max {
		*xs = append(*xs, Intersection{T: t0, S: cyl})
	}

	y1 := ray.Origin[1] + t1*ray.Dir[1]
	if cyl.Min < y1 && y1 < cyl.Max {
		*xs = append(*xs, Intersection{T: t1, S: cyl})
	}
}

func (c *Cylinder) intersectCaps(ray Ray, isects *[]Intersection) {
	if !c.Closed || m.EqApprox(ray.Dir[1], 0.0) {
		return
	}

	t := (c.Min - ray.Origin[1]) / ray.Dir[1]
	if c.checkCaps(ray, t) {
		*isects = append(*isects, Intersection{T: t, S: c})
	}

	t = (c.Max - ray.Origin[1]) / ray.Dir[1]
	if c.checkCaps(ray, t) {
		*isects = append(*isects, Intersection{T: t, S: c})
	}
}

func (c *Cylinder) checkCaps(ray Ray, t float64) bool {
	x := ray.Origin[0] + t*ray.Dir[0]
	z := ray.Origin[2] + t*ray.Dir[2]
	return (x*x + z*z) <= 1.0
}

func (c *Cylinder) GetSavedRay() Ray {
	return c.O.SavedRay
}

func (c *Cylinder) NormalAt(point m.Vec4) m.Vec4 {
	localPoint := c.O.InvTf.MulVec(point)
	localNormal := c.localNormalAt(localPoint)
	nWorld := c.O.InvTf.Tpose().MulVec(localNormal)
	nWorld[3] = 0
	return nWorld.Normalize()
}

func (c *Cylinder) localNormalAt(point m.Vec4) m.Vec4 {
	dist := point[0]*point[0] + point[2]*point[2]
	if dist < 1 && point[1] >= (c.Max-m.EPSILON) {
		return m.Vector4(0, 1, 0)
	}
	if dist < 1 && point[1] <= (c.Min+m.EPSILON) {
		return m.Vector4(0, -1, 0)
	}
	return m.Vector4(point[0], 0, point[2])
}

func (c *Cylinder) SamplePatternAt(point m.Vec4) m.Vec4 {
	objPoint := c.O.InvTf.MulVec(point)
	pattern := c.O.Material.Pattern
	patternPoint := pattern.GetInvTf().MulVec(objPoint)
	return pattern.SampleAt(patternPoint)
}

func (c *Cylinder) SetPattern(p pattern.Pattern) {
	c.O.Material.Pattern = p
}

// ////////////// CONE ////////////////
func NewCone() Cone {
	return Cone{
		O:      defaultObject(),
		Min:    math.Inf(-1),
		Max:    math.Inf(1),
		Closed: false,
	}
}

func (c *Cone) GetTf() m.Mat4 {
	return c.O.Tf
}

func (c *Cone) SetTf(tf m.Mat4) {
	c.O.Tf = tf
	c.O.InvTf = tf.Inv()
}

func (c *Cone) GetMat() mtl.Material {
	return c.O.Material
}

func (c *Cone) SetMat(mtl mtl.Material) {
	c.O.Material = mtl
}

func (c *Cone) Intersect(ray Ray) ([]Intersection, int) {
	ray.Transform(c.O.InvTf)
	return c.localIntersect(ray)
}

func (cone *Cone) localIntersect(ray Ray) ([]Intersection, int) {
	a := ray.Dir[0]*ray.Dir[0] - ray.Dir[1]*ray.Dir[1] + ray.Dir[2]*ray.Dir[2]

	b := 2*ray.Origin[0]*ray.Dir[0] - 2*ray.Origin[1]*ray.Dir[1] + 2*ray.Origin[2]*ray.Dir[2]

	c := ray.Origin[0]*ray.Origin[0] - ray.Origin[1]*ray.Origin[1] + ray.Origin[2]*ray.Origin[2]

	disc := b*b - 4*a*c

	var xs []Intersection
	if disc < 0.0 {
		return xs, 0
	}

	if m.EqApprox(a, 0.0) && m.EqApprox(b, 0.0) {
		return []Intersection{}, 0
	}

	if m.EqApprox(a, 0.0) && !m.EqApprox(b, 0.0) {
		xs = append(xs, Intersection{T: -c / (2 * b), S: cone})
	}

	t0 := (-b - math.Sqrt(disc)) / (2 * a)
	t1 := (-b + math.Sqrt(disc)) / (2 * a)

	if t0 > t1 {
		t0, t1 = t1, t0
	}

	y0 := ray.Origin[1] + t0*ray.Dir[1]
	if cone.Min < y0 && y0 < cone.Max {
		xs = append(xs, Intersection{T: t0, S: cone})
	}

	y1 := ray.Origin[1] + t1*ray.Dir[1]
	if cone.Min < y1 && y1 < cone.Max {
		xs = append(xs, Intersection{T: t1, S: cone})
	}

	cone.intersectCaps(ray, &xs)
	return xs, len(xs)
}

func (c *Cone) intersectCaps(ray Ray, isects *[]Intersection) {
	if !c.Closed || m.EqApprox(ray.Dir[1], 0.0) {
		return
	}

	t := (c.Min - ray.Origin[1]) / ray.Dir[1]
	if c.checkCaps(ray, t, c.Min) {
		*isects = append(*isects, Intersection{T: t, S: c})
	}

	t = (c.Max - ray.Origin[1]) / ray.Dir[1]
	if c.checkCaps(ray, t, c.Max) {
		*isects = append(*isects, Intersection{T: t, S: c})
	}
}

func (c *Cone) checkCaps(ray Ray, t, y float64) bool {
	x := ray.Origin[0] + t*ray.Dir[0]
	z := ray.Origin[2] + t*ray.Dir[2]
	return (x*x + z*z) <= math.Abs(y)
}

func (c *Cone) GetSavedRay() Ray {
	return c.O.SavedRay
}

func (c *Cone) NormalAt(point m.Vec4) m.Vec4 {
	localPoint := c.O.InvTf.MulVec(point)
	localNormal := c.localNormalAt(localPoint)
	nWorld := c.O.InvTf.Tpose().MulVec(localNormal)
	nWorld[3] = 0
	return nWorld.Normalize()
}

func (c *Cone) localNormalAt(point m.Vec4) m.Vec4 {
	dist := point[0]*point[0] + point[2]*point[2]
	if dist < 1 && point[1] >= (c.Max-m.EPSILON) {
		return m.Vector4(0, 1, 0)
	}
	if dist < 1 && point[1] <= (c.Min+m.EPSILON) {
		return m.Vector4(0, -1, 0)
	}

	y := math.Sqrt(point[0]*point[0] + point[2]*point[2])
	if point[1] > 0.0 {
		y = -y
	}

	return m.Vector4(point[0], y, point[2])
}

func (c *Cone) SamplePatternAt(point m.Vec4) m.Vec4 {
	objPoint := c.O.InvTf.MulVec(point)
	pattern := c.O.Material.Pattern
	patternPoint := pattern.GetInvTf().MulVec(objPoint)
	return pattern.SampleAt(patternPoint)
}

func (c *Cone) SetPattern(p pattern.Pattern) {
	c.O.Material.Pattern = p
}

// ////////////// INTERSECTIONS ////////////////
func Intersections(isects ...Intersection) []Intersection {
	return isects
}

func Hit(isects []Intersection) (Intersection, bool) {
	res := Intersection{T: math.MaxFloat64}
	isHit := false
	for i := 0; i < len(isects); i++ {
		if isects[i].T <= res.T && isects[i].T >= 0 {
			res = Intersection{S: isects[i].S, T: isects[i].T}
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
	op := pos
	op.Add(normal.Mul(m.EPSILON))
	up := pos
	up.Sub(normal.Mul(m.EPSILON))

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

func computeN(isect Intersection, isects []Intersection) (n1 float64, n2 float64) {
	containers := []Shape{}
	for i := 0; i < len(isects); i++ {
		if isects[i] == isect {
			if len(containers) == 0 {
				n1 = 1.0
			} else {
				n1 = containers[len(containers)-1].GetMat().RefractiveIndex
			}
		}
		if pie.Contains(containers, isects[i].S) {
			containers = pie.FilterNot(containers, func(shape Shape) bool { return shape == isects[i].S })
		} else {
			containers = append(containers, isects[i].S)
		}

		if isects[i] == isect {
			if len(containers) == 0 {
				n2 = 1.0
			} else {
				n2 = containers[len(containers)-1].GetMat().RefractiveIndex
			}
			break
		}
	}
	return
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
