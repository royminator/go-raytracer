package math

type (
	Mat4 [16]float64
)

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
		Vec4{v[0], 0, 0 ,0},
		Vec4{0, v[1], 0 ,0},
		Vec4{0, 0, v[2] ,0},
		Vec4{0, 0, 0 ,v[3]},
	)
}

func (mat Mat4) At(m, n int) float64 {
	return mat[4*m+n]
}

func (m1 Mat4) Mat4Add(m2 Mat4) Mat4 {
	for m := 0; m < 4; m++ {
		for n := 0; n < 4; n++ {
			m1[mat4Index(m, n)] = m1.At(m, n)+m2.At(m, n)
		}
	}
	return m1
}

func (m1 Mat4) Mat4Sub(m2 Mat4) Mat4 {
	for m := 0; m < 4; m++ {
		for n := 0; n < 4; n++ {
			m1[mat4Index(m, n)] = m1.At(m, n)-m2.At(m, n)
		}
	}
	return m1
}

func mat4Index(m, n int) int {
	return 4*m+n
}
