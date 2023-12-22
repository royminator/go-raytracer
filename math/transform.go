package math

import "math"

func Trans(t Vec3) Mat4 {
	id := Mat4Ident()
	id.SetCol(3, t.ToVec4())
	return id
}

func Scale(t Vec3) Mat4 {
	return Mat4Diag(t.ToVec4())
}

func RotX(r float64) Mat4 {
	return Mat4FromRows(
		Vec4{1, 0, 0, 0},
		Vec4{0, math.Cos(r), -math.Sin(r), 0},
		Vec4{0, math.Sin(r), math.Cos(r), 0},
		Vec4{0, 0, 0, 1},
	)
}

func RotY(r float64) Mat4 {
	return Mat4FromRows(
		Vec4{math.Cos(r), 0, math.Sin(r), 0},
		Vec4{0, 1, 0, 0},
		Vec4{-math.Sin(r), 0, math.Cos(r), 0},
		Vec4{0, 0, 0, 1},
	)
}

func RotZ(r float64) Mat4 {
	return Mat4FromRows(
		Vec4{math.Cos(r), -math.Sin(r), 0, 0},
		Vec4{math.Sin(r), math.Cos(r), 0, 0},
		Vec4{0, 0, 1, 0},
		Vec4{0, 0, 0, 1},
	)
}

func Shear(xy, xz, yx, yz, zx, zy float64) Mat4 {
	return Mat4FromRows(
		Vec4{1, xy, xz, 0},
		Vec4{yx, 1, yz, 0},
		Vec4{zx, zy, 1, 0},
		Vec4{0, 0, 0, 1},
	)
}
