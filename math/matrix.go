package math

type (
	Mat2 [4]float64
	Mat3 [9]float64
	Mat4 [16]float64
)

// //////////////////////////// MAT2 //////////////////////////////
func Mat2FromRows(m0, m1 Vec2) Mat2 {
	return Mat2{
		m0[0], m0[1],
		m1[0], m1[1],
	}
}

func (mat Mat2) At(m, n int) float64 {
	return mat[2*m+n]
}

// //////////////////////////// MAT3 //////////////////////////////
func Mat3FromRows(m0, m1, m2 Vec3) Mat3 {
	return Mat3{
		m0[0], m0[1], m0[2],
		m1[0], m1[1], m1[2],
		m2[0], m2[1], m2[2],
	}
}

func (mat Mat3) At(m, n int) float64 {
	return mat[3*m+n]
}

// //////////////////////////// MAT4 //////////////////////////////
func Mat4FromRows(m0, m1, m2, m3 Vec4) Mat4 {
	return Mat4{
		m0[0], m0[1], m0[2], m0[3],
		m1[0], m1[1], m1[2], m1[3],
		m2[0], m2[1], m2[2], m2[3],
		m3[0], m3[1], m3[2], m3[3],
	}
}

func Mat4Ident() Mat4 {
	return Mat4FromRows(
		Vec4{1, 0, 0, 0},
		Vec4{0, 1, 0, 0},
		Vec4{0, 0, 1, 0},
		Vec4{0, 0, 0, 1},
	)
}

func Mat4Diag(v Vec4) Mat4 {
	return Mat4FromRows(
		Vec4{v[0], 0, 0, 0},
		Vec4{0, v[1], 0, 0},
		Vec4{0, 0, v[2], 0},
		Vec4{0, 0, 0, v[3]},
	)
}

func (mat Mat4) At(m, n int) float64 {
	return mat[4*m+n]
}

func (mat *Mat4) Set(m, n int, v float64) {
	mat[4*m+n] = v
}

func (m1 Mat4) Add(m2 Mat4) Mat4 {
	for m := 0; m < 4; m++ {
		for n := 0; n < 4; n++ {
			m1.Set(m, n, m1.At(m, n)+m2.At(m, n))
		}
	}
	return m1
}

func (m1 Mat4) Mat4Sub(m2 Mat4) Mat4 {
	for m := 0; m < 4; m++ {
		for n := 0; n < 4; n++ {
			m1.Set(m, n, m1.At(m, n)-m2.At(m, n))
		}
	}
	return m1
}

func (mat Mat4) Index(m, n int) int {
	return 4*m + n
}

func (m1 Mat4) Mul(m2 Mat4) Mat4 {
	var res Mat4
	for m := 0; m < 4; m++ {
		for n := 0; n < 4; n++ {
			v := m1.At(m, 0)*m2.At(0, n)+
				m1.At(m, 1)*m2.At(1, n)+
				m1.At(m, 2)*m2.At(2, n)+
				m1.At(m, 3)*m2.At(3, n)
			res.Set(m, n, v)
		}
	}
	return res
}

func (mat Mat4) MulVec(vec Vec4) Vec4 {
	var res Vec4
	for m := 0; m < 4; m++ {
		res[m] = mat.At(m, 0)*vec[0]+
			mat.At(m, 1)*vec[1]+
			mat.At(m, 2)*vec[2]+
			mat.At(m, 3)*vec[3]
	}
	return res
}
