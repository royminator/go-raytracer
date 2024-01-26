package math

import (
	"math"
)

type (
	Vec2 [2]float64
	Vec3 [3]float64
	Vec4 [4]float64
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

func Color4(r, g, b, a float64) Vec4 {
	return Vec4{r, g, b, a}
}

func Vec4With(v float64) Vec4 {
	return Vec4{v, v, v, v}
}

func (lhs *Vec4) Add(rhs Vec4) {
	lhs[0] += rhs[0]
	lhs[1] += rhs[1]
	lhs[2] += rhs[2]
	lhs[3] += rhs[3]
}

func (lhs *Vec4) Sub(rhs Vec4) {
	lhs[0] -= rhs[0]
	lhs[1] -= rhs[1]
	lhs[2] -= rhs[2]
	lhs[3] -= rhs[3]
}

func (v Vec4) Mul(x float64) Vec4 {
	return Vec4{v[0] * x, v[1] * x, v[2] * x, v[3] * x}
}

func (v Vec4) Div(x float64) Vec4 {
	return Vec4{v[0] / x, v[1] / x, v[2] / x, v[3] / x}
}

func (v Vec4) Negate() Vec4 {
	return v.Mul(-1)
}

func (v Vec4) Magnitude() float64 {
	return math.Sqrt(v[0]*v[0] + v[1]*v[1] + v[2]*v[2] + v[3]*v[3])
}

func (v Vec4) Normalize() Vec4 {
	return v.Div(v.Magnitude())
}

func (lhs Vec4) ApproxEqual(rhs Vec4) bool {
	for i := range lhs {
		if !EqApprox(lhs[i], rhs[i]) {
			return false
		}
	}
	return true
}

func (lhs Vec4) Dot(rhs Vec4) float64 {
	return lhs[0]*rhs[0] + lhs[1]*rhs[1] + lhs[2]*rhs[2] + lhs[3]*rhs[3]
}

func (lhs Vec4) Cross(rhs Vec4) Vec4 {
	return Vector4(
		lhs[1]*rhs[2]-lhs[2]*rhs[1],
		lhs[2]*rhs[0]-lhs[0]*rhs[2],
		lhs[0]*rhs[1]-lhs[1]*rhs[0],
	)
}

func (lhs Vec4) MulElem(rhs Vec4) Vec4 {
	return Vec4{
		lhs[0] * rhs[0],
		lhs[1] * rhs[1],
		lhs[2] * rhs[2],
		lhs[3] * rhs[3],
	}
}

func (vec *Vec4) MulMat(m Mat4) {
	v := *vec
	vec[0] = m[0]*v[0]+m[1]*v[1]+m[2]*v[2]+m[3]*v[3]
	vec[1] = m[4]*v[0]+m[5]*v[1]+m[6]*v[2]+m[7]*v[3]
	vec[2] = m[8]*v[0]+m[9]*v[1]+m[10]*v[2]+m[11]*v[3]
	vec[3] = m[12]*v[0]+m[13]*v[1]+m[14]*v[2]+m[15]*v[3]
}

func (v Vec4) Reflect(n Vec4) Vec4 {
	v.Sub(n.Mul(2).Mul(v.Dot(n)))
	return v
}

func (v Vec4) ToVec3() Vec3 {
	return Vec3{v[0], v[1], v[2]}
}

func (v Vec3) ToVec4() Vec4 {
	return Vec4{v[0], v[1], v[2], 1}
}

func EqApprox(a, b float64) bool {
	return math.Abs(a-b) < EPSILON
}
