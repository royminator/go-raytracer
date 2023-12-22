package math

type (
	Mat2 [4]float64
	Mat3 [9]float64
	Mat4 [16]float64
	Mat2x3 [6]float64
	Mat3x4 [12]float64
)

// //////////////////////////// MAT2 //////////////////////////////
func Mat2FromRows(m0, m1 Vec2) Mat2 {
	return Mat2{
		m0[0], m0[1],
		m1[0], m1[1],
	}
}


func Mat2FromCols(col1, col2 Vec2) Mat2 {
	return Mat2{
		col1[0], col2[0],
		col1[1], col2[1],
	}
}

func (mat Mat2) At(m, n int) float64 {
	return mat[2*m+n]
}

func (mat Mat2) Det() float64 {
	return mat[0]*mat[3]-mat[1]*mat[2]
}

func (mat *Mat2) Set(m, n int, v float64) {
	mat[2*m+n] = v
}

// //////////////////////////// MAT3 //////////////////////////////
func Mat3FromRows(m0, m1, m2 Vec3) Mat3 {
	return Mat3{
		m0[0], m0[1], m0[2],
		m1[0], m1[1], m1[2],
		m2[0], m2[1], m2[2],
	}
}

func Mat3FromCols(col0, col1, col2 Vec3) Mat3 {
	return Mat3{
		col0[0], col1[0], col2[0],
		col0[1], col1[1], col2[1],
		col0[2], col1[2], col2[2],
	}
}

func (mat Mat3) At(m, n int) float64 {
	return mat[3*m+n]
}

func (mat *Mat3) Set(m, n int, v float64) {
	mat[3*m+n] = v
}

func (mat Mat3) SubMat(row, col int) Mat2 {
	return mat.DeleteRow(row).DeleteCol(col)
}

func (mat Mat3) DeleteRow(row int) Mat2x3 {
	vecs := make([]Vec3, 2)
	var i = 0
	for m := 0; m < 2; m++ {
		if m == row {
			i++
		}
		vecs[m] = mat.Row(i)
		i++
	}
	return Mat2x3{
		vecs[0][0], vecs[0][1], vecs[0][2],
		vecs[1][0], vecs[1][1], vecs[1][2],
	}
}

func (mat Mat3) Row(m int) Vec3 {
	return Vec3{mat.At(m, 0), mat.At(m, 1), mat.At(m, 2)}
}

func (mat Mat3) Minor(m, n int) float64 {
	return mat.SubMat(m, n).Det()
}

func (mat Mat3) Cofactor(m, n int) float64 {
	var sign float64 = -1
	if (m + n) % 2 == 0 {
		sign = 1
	}
	return mat.Minor(m, n)*sign
}

func (mat Mat3) Det() float64 {
	return mat.At(0, 0)*mat.Cofactor(0, 0)+
		mat.At(0, 1)*mat.Cofactor(0, 1)+
		mat.At(0, 2)*mat.Cofactor(0, 2)
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

func (mat Mat4) Tpose() Mat4 {
	var res Mat4
	for m := 0; m < 4; m++ {
		for n := 0; n < 4; n++ {
			res.Set(m, n, mat.At(n, m))
		}
	}
	return res
}

func (mat Mat4) SubMat(row, col int) Mat3 {
	return mat.DeleteRow(row).DeleteCol(col)
}

func (mat Mat4) DeleteRow(row int) Mat3x4 {
	rows := make([]Vec4, 3)
	var i = 0
	for m := 0; m < 3; m++ {
		if m == row {
			i++
		}
		rows[m] = mat.Row(i)
		i++
	}
	return Mat3x4{
		rows[0][0], rows[0][1], rows[0][2], rows[0][3],
		rows[1][0], rows[1][1], rows[1][2], rows[1][3],
		rows[2][0], rows[2][1], rows[2][2], rows[2][3],
	}
}

func (mat Mat4) Row(m int) Vec4 {
	return Vec4{mat.At(m, 0), mat.At(m, 1), mat.At(m, 2), mat.At(m, 3)}
}

func (mat Mat4) Cofactor(m, n int) float64 {
	var sign = -1.0
	if (m + n) % 2 == 0 {
		sign = 1.0
	}
	return sign*mat.SubMat(m, n).Det()
}

func (mat Mat4) Det() float64 {
	return mat.At(0, 0)*mat.Cofactor(0, 0)+
		mat.At(0, 1)*mat.Cofactor(0, 1)+
		mat.At(0, 2)*mat.Cofactor(0, 2)+
		mat.At(0, 3)*mat.Cofactor(0, 3)
}

func (mat Mat4) IsInvertible() (bool, float64) {
	det := mat.Det()
	return !eqApprox(det, 0.0), det
}

func (mat Mat4) Inv() Mat4 {
	isInvertible, det := mat.IsInvertible()
	if !isInvertible {
		panic("matrix not invertible")
	}

	var res Mat4
	for m := 0; m < 4; m++ {
		for n := 0; n < 4; n++ {
			res.Set(n, m, mat.Cofactor(m, n)/det)
		}
	}
	return res
}

func (m1 Mat4) ApproxEqual(m2 Mat4) bool {
	for m := 0; m < 4; m++ {
		for n := 0; n < 4; n++ {
			if !eqApprox(m1.At(m, n), m2.At(m, n)) {
				return false
			}
		}
	}
	return true
}

// //////////////////////////// MAT2x3 //////////////////////////////
func (mat Mat2x3) DeleteCol(col int) Mat2 {
	cols := make([]Vec2, 2)
	var i = 0
	for n := 0; n < 2; n++ {
		if n == col {
			i++
		}
		cols[n] = mat.Col(i)
		i++
	}
	return Mat2FromCols(cols[0], cols[1])
}

func (mat Mat2x3) Col(n int) Vec2 {
	return Vec2{mat[n], mat[n+3]}
}

// //////////////////////////// MAT3x4 //////////////////////////////
func (mat Mat3x4) DeleteCol(col int) Mat3 {
	cols := make([]Vec3, 3)
	var i = 0
	for n := 0; n < 3; n++ {
		if n == col {
			i++
		}
		cols[n] = mat.Col(i)
		i++
	}
	return Mat3FromCols(cols[0], cols[1], cols[2])
}

func (mat Mat3x4) Col(n int) Vec3 {
	return Vec3{mat[n], mat[n+4], mat[n+8]}
}
