package ray

import (
	"math"

	"github.com/google/uuid"
	m "roytracer/math"
)

type (
	Ray struct {
		Origin m.Vec4
		Dir    m.Vec4
	}

	Sphere struct {
		Radius float64
		Center m.Vec4
		Id     uuid.UUID
	}

	Intersection struct {
		T  []float64
		Id uuid.UUID
	}
)

func (r Ray) Pos(t float64) m.Vec4 {
	return r.Origin.Add(r.Dir.Mul(t))
}

func (s Sphere) Intersect(r Ray) Intersection {
	sphereToRay := r.Origin.Sub(m.Point4(0, 0, 0))
	a := r.Dir.Dot(r.Dir)
	b := 2 * r.Dir.Dot(sphereToRay)
	c := sphereToRay.Dot(sphereToRay) - 1
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return Intersection{Id: s.Id}
	}
	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)
	if t1 < t2 {
		return Intersection{Id: s.Id, T: []float64{t1, t2}}
	}
	return Intersection{Id: s.Id, T: []float64{t2, t1}}
}

func NewSphere() Sphere {
	return Sphere{0, m.Vector4(0, 0, 0), uuid.New()}
}
