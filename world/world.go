package world

import (
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/shape"
	"sort"
)

type World struct {
	Objects []*shape.Sphere
	Light   mtl.PointLight
}

func DefaultWorld() World {
	s1 := shape.NewSphere()
	s1.Material.Color = m.Color4(0.8, 1.0, 0.6, 0.0)
	s1.Material.Diffuse = 0.7
	s1.Material.Specular = 0.2

	s2 := shape.NewSphere()
	s2.SetTf(m.Scale(m.Vec3{0.5, 0.5, 0.5}))

	return World{
		Light: mtl.PointLight{
			Intensity: m.Color4(1, 1, 1, 0),
			Pos: m.Point4(-10, 10, -10),
		},
		Objects: []*shape.Sphere{&s1, &s2},
	}
}

func (w *World) Intersect(ray shape.Ray) []shape.Intersection {
	var isects []shape.Intersection
	for _, obj := range w.Objects {
		isect := obj.Intersect(ray)
		isects = append(isects, isect...)
	}
	sort.Slice(isects, func(i, j int) bool {
		return isects[i].T < isects[j].T
	})
	return isects
}

func (w *World) ShadeHit(comps shape.IntersectionComps) m.Vec4 {
	return mtl.Lighting(comps.O.Material, w.Light,
		comps.Point, comps.Eye, comps.Normal)
}


func (w *World) ColorAt(ray shape.Ray) m.Vec4 {
	isects := w.Intersect(ray)
	if hit, isHit := shape.Hit(isects); isHit {
		comps := hit.Prepare(ray)
		return w.ShadeHit(comps)
	}
	return m.Vec4With(0)
}
