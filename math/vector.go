package math

import (
	"math"
)

type (
	Vec4 struct { X, Y, Z, W float64 }
)

const (
	EPSILON = 0.0001
)

func Point4(x, y, z float64) Vec4 {
	return Vec4{x, y, z, 1.0}
}

func Vector4(x, y, z float64) Vec4 {
	return Vec4{x, y, z, 0.0}
}

func Vec4With(v float64) Vec4 {
	return Vec4{v, v, v, v}
}

func (lhs *Vec4) Add(rhs Vec4) Vec4 {
	return Vec4{lhs.X+rhs.X, lhs.Y+rhs.Y, lhs.Z+rhs.Z, lhs.W+rhs.W}
}

func (lhs *Vec4) Sub(rhs Vec4) Vec4 {
	return Vec4{lhs.X-rhs.X, lhs.Y-rhs.Y, lhs.Z-rhs.Z, lhs.W-rhs.W}
}

func (v *Vec4) Mul(x float64) Vec4 {
	return Vec4{v.X*x, v.Y*x, v.Z*x, v.W*x}
}

func (v *Vec4) Div(x float64) Vec4 {
	return Vec4{v.X/x, v.Y/x, v.Z/x, v.W/x}
}

func (v *Vec4) Negate() Vec4 {
	return v.Mul(-1)
}

func (v *Vec4) Magnitude() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z + v.W*v.W)
}

func eqApprox(a, b float64) bool {
	if math.Abs(a - b) < EPSILON {
		return true
	}
	return false
}

