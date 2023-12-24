package ray

import (
	"fmt"
	"math"

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
	}

	Intersection struct {
		IsIntersect bool
		P1          m.Vec4
		P2          m.Vec4
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
	fmt.Println("a:", a)
	fmt.Println("b:", b)
	fmt.Println("c:", c)
	discriminant := b*b - 4*a*c
	fmt.Println("d:", discriminant)
	if discriminant < 0 {
		return Intersection{IsIntersect: false}
	}
	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)
	fmt.Println("t1:", t1)
	fmt.Println("t2:", t2)
	return Intersection{true, r.Origin.Sub(r.Pos(t1)), r.Pos(t2)}
}
