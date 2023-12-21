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

func (mat Mat4) At(m, n int) float64 {
	return mat[4*m+n]
}
