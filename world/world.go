package world

import (
	"math"
	"sort"

	"roytracer/color"
	"roytracer/light"
	m "roytracer/math"
	"roytracer/shape"
)

type World struct {
	Objects []shape.Shape
	Light   light.PointLight
}

func DefaultWorld() *World {
	s1 := shape.NewSphere()
	s1.O.Material.Color = m.Color4(0.8, 1.0, 0.6, 0.0)
	s1.O.Material.Diffuse = 0.7
	s1.O.Material.Specular = 0.2

	s2 := shape.NewSphere()
	s2.SetTf(m.Scale(0.5, 0.5, 0.5))

	return &World{
		Light: light.PointLight{
			Intensity: m.Color4(1, 1, 1, 0),
			Pos:       m.Point4(-10, 10, -10),
		},
		Objects: []shape.Shape{&s1, &s2},
	}
}

func (w *World) AddShape(s shape.Shape) {
	w.Objects = append(w.Objects, s)
}

func (w *World) Intersect(ray shape.Ray) []shape.Intersection {
	var isects []shape.Intersection
	for i := 0; i < len(w.Objects); i++ {
		isect := w.Objects[i].Intersect(ray)
		isects = append(isects, isect...)
	}
	sort.Slice(isects, func(i, j int) bool {
		return isects[i].T < isects[j].T
	})
	return isects
}

func (w *World) ShadeHit(comps shape.IntersectionComps, remaining int) m.Vec4 {
	shadowed := w.IsShadowed(comps.OverPoint)
	surface := light.Lighting(comps.S.GetMat(), comps.S, w.Light, comps.OverPoint,
		comps.Eye, comps.Normal, shadowed)
	reflected := w.ReflectedColor(comps, remaining)
	refracted := w.RefractedColor(comps, remaining)
	mtl := comps.S.GetMat()
	if mtl.Reflective > 0.0 && mtl.Transparency > 0.0 {
		reflectance := comps.Schlick()
		res := surface
		res.Add(reflected.Mul(reflectance))
		res.Add(refracted.Mul(1.0 - reflectance))
		return res
	}
	res := surface
	res.Add(reflected)
	res.Add(refracted)
	return res
}

func (w *World) ColorAt(ray shape.Ray, remaining int) m.Vec4 {
	isects := w.Intersect(ray)
	if hit, isHit := shape.Hit(isects); isHit {
		comps := hit.Prepare(ray, isects)
		return w.ShadeHit(comps, remaining)
	}
	return m.Vec4With(0)
}

func (w *World) IsShadowed(p m.Vec4) bool {
	v := w.Light.Pos
	v.Sub(p)
	dist := v.Magnitude()
	dir := v.Normalize()
	r := shape.Ray{Origin: p, Dir: dir}
	isects := w.Intersect(r)
	if hit, isHit := shape.Hit(isects); isHit && hit.T < dist {
		return true
	}
	return false
}

func (w *World) ReflectedColor(comps shape.IntersectionComps, remaining int) m.Vec4 {
	if remaining <= 0 {
		return color.Black
	}
	mat := comps.S.GetMat()
	if mat.Reflective == 0.0 {
		return m.Vec4With(0)
	}
	reflectedRay := shape.Ray{
		Origin: comps.OverPoint,
		Dir:    comps.Reflect,
	}
	color := w.ColorAt(reflectedRay, remaining-1)
	return color.Mul(mat.Reflective)
}

func (w *World) RefractedColor(comps shape.IntersectionComps, remaining int) m.Vec4 {
	if remaining <= 0 {
		return color.Black
	}
	if comps.S.GetMat().Transparency == 0.0 {
		return color.Black
	}
	return w.hasTotalInternalReflection(comps, remaining)
}

func (w *World) hasTotalInternalReflection(comps shape.IntersectionComps, remaining int) m.Vec4 {
	nRatio := comps.N1 / comps.N2
	cosI := comps.Eye.Dot(comps.Normal)
	sin2t := nRatio * nRatio * (1.0 - cosI*cosI)
	if sin2t > 1.0 {
		return color.Black
	}
	cosT := math.Sqrt(1.0 - sin2t)
	dir := comps.Normal.Mul(nRatio*cosI - cosT)
	dir.Sub(comps.Eye.Mul(nRatio))
	refractedRay := shape.Ray{
		Origin: comps.UnderPoint,
		Dir:    dir,
	}
	return w.ColorAt(refractedRay, remaining-1).Mul(comps.S.GetMat().Transparency)
}
