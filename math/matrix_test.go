package math

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

// //////////////////////////// MAT2 //////////////////////////////
func TestMat2FromRows(t *testing.T) {
	mat := Mat2FromRows(
		Vec2{-3, 5},
		Vec2{1, -2},
	)

	assert := assert.New(t)
	assert.Equal(-3.0, mat.At(0, 0))
	assert.Equal(5.0, mat.At(0, 1))
	assert.Equal(1.0, mat.At(1, 0))
	assert.Equal(-2.0, mat.At(1, 1))
}

func TestMat2Det(t *testing.T) {
	expected := 17.0
	a := Mat2FromRows(
		Vec2{1, 5},
		Vec2{-3, 2},
	)
	assert.Equal(t, expected, a.Det())
}

// //////////////////////////// MAT3 //////////////////////////////
func TestMat3FromRows(t *testing.T) {
	mat := Mat3FromRows(
		Vec3{-3, 5, 0},
		Vec3{1, -2, -7},
		Vec3{0, 1, 1},
	)

	assert := assert.New(t)
	assert.Equal(-3.0, mat.At(0, 0))
	assert.Equal(-2.0, mat.At(1, 1))
	assert.Equal(1.0, mat.At(2, 2))
}

func TestMat3SubMat(t *testing.T) {
	type testData struct { i, j int; res Mat2 }
	td := []testData{
		{0, 0, Mat2{2, 7, 6, -3}},
		{0, 1, Mat2{-3, 7, 0, -3}},
		{2, 2, Mat2{1, 5, -3, 2}},
	}
	a := Mat3FromRows(
		Vec3{1, 5, 0},
		Vec3{-3, 2, 7},
		Vec3{0, 6, -3},
	)
	assert := assert.New(t)
	for _, d := range td {
		assert.Equal(d.res, a.SubMat(d.i, d.j))
	}
}

func TestMat3Minor(t *testing.T) {
	a := Mat3FromRows(
		Vec3{3, 5, 0},
		Vec3{2, -1, -7},
		Vec3{6, -1, 5},
	)
	actual := a.Minor(1, 0)
	assert.Equal(t, a.SubMat(1, 0).Det(), actual)
	assert.Equal(t, 25.0, actual)
}

func TestMat3Cofactor(t *testing.T) {
	a := Mat3FromRows(
		Vec3{3, 5, 0},
		Vec3{2, -1, -7},
		Vec3{6, -1, 5},
	)
	assert := assert.New(t)
	assert.Equal(-12.0, a.Minor(0, 0))
	assert.Equal(-12.0, a.Cofactor(0, 0))
	assert.Equal(25.0, a.Minor(1, 0))
	assert.Equal(-25.0, a.Cofactor(1, 0))
}

func TestMat3Det(t *testing.T) {
	a := Mat3FromRows(
		Vec3{1, 2, 6},
		Vec3{-5, 8, -4},
		Vec3{2, 6, 4},
	)
	assert := assert.New(t)
	assert.Equal(56.0, a.Cofactor(0, 0))
	assert.Equal(12.0, a.Cofactor(0, 1))
	assert.Equal(-46.0, a.Cofactor(0, 2))
	assert.Equal(-196.0, a.Det())
}

// //////////////////////////// MAT4 //////////////////////////////
func TestMat4FromRows(t *testing.T) {
	mat := Mat4FromRows(
		Vec4{1, 2, 3, 4},
		Vec4{5, 6, 7, 8},
		Vec4{9, 10, 11, 12},
		Vec4{13, 14, 15, 16},
	)

	assert := assert.New(t)
	assert.Equal(1.0, mat.At(0, 0));
	assert.Equal(15.0, mat.At(3, 2));
	assert.Equal(12.0, mat.At(2, 3));
	assert.Equal(7.0, mat.At(1, 2));
}

func TestMat4Identity(t *testing.T) {
	expected := Mat4FromRows(
		Vec4{1, 0, 0, 0},
		Vec4{0, 1, 0, 0},
		Vec4{0, 0, 1, 0},
		Vec4{0, 0, 0, 1},
	)
	actual := Mat4Ident()

	assert.Equal(t, expected, actual)
}

func TestMat4Diag(t *testing.T) {
	expected := Mat4Ident()
	actual := Mat4Diag(Vec4{1, 1, 1, 1})
	assert.Equal(t, expected, actual)
}

func TestMat4Add(t *testing.T) {
	expected := Mat4FromRows(
		Vec4{1, 8, 8, 9},
		Vec4{0, -2, 0, 3},
		Vec4{-12, 0, 1, 0},
		Vec4{0, 17, 0, 1},
	)
	m1 := Mat4FromRows(
		Vec4{0, 7, 1, 1},
		Vec4{0, -2, 0, 2},
		Vec4{1, -2, 9, 0},
		Vec4{0, 5, 6, -7},
	)
	m2 := Mat4FromRows(
		Vec4{1, 1, 7, 8},
		Vec4{0, 0, 0, 1},
		Vec4{-13, 2, -8, 0},
		Vec4{0, 12, -6, 8},
	)
	actual := m1.Add(m2)
	assert.Equal(t, expected, actual)
}

func TestMat4Sub(t *testing.T) {
	expected := Mat4FromRows(
		Vec4{-1, 6, -6, -7},
		Vec4{0, -2, 0, 1},
		Vec4{14, -4, 17, 0},
		Vec4{0, -7, 12, -15},
	)
	m1 := Mat4FromRows(
		Vec4{0, 7, 1, 1},
		Vec4{0, -2, 0, 2},
		Vec4{1, -2, 9, 0},
		Vec4{0, 5, 6, -7},
	)
	m2 := Mat4FromRows(
		Vec4{1, 1, 7, 8},
		Vec4{0, 0, 0, 1},
		Vec4{-13, 2, -8, 0},
		Vec4{0, 12, -6, 8},
	)
	actual := m1.Mat4Sub(m2)
	assert.Equal(t, expected, actual)
}

func TestMat4Mul(t *testing.T) {
	expected := Mat4FromRows(
		Vec4{20, 22, 50, 48},
		Vec4{44, 54, 114, 108},
		Vec4{40, 58, 110, 102},
		Vec4{16, 26, 46, 42},
	)
	a := Mat4FromRows(
		Vec4{1, 2, 3, 4},
		Vec4{5, 6, 7, 8},
		Vec4{9, 8, 7, 6},
		Vec4{5, 4, 3, 2},
	)
	b := Mat4FromRows(
		Vec4{-2, 1, 2, 3},
		Vec4{3, 2, 1, -1},
		Vec4{4, 3, 6, 5},
		Vec4{1, 2, 7, 8},
	)
	assert.Equal(t, expected, a.Mul(b))
}

func TestMat4MulVec(t *testing.T) {
	expected := Vec4{18, 24, 33, 1}
	a := Mat4FromRows(
		Vec4{1, 2, 3, 4},
		Vec4{2, 4, 4, 2},
		Vec4{8, 6, 4, 1},
		Vec4{0, 0, 0, 1},
	)
	b := Vec4{1, 2, 3, 1}
	assert.Equal(t, expected, a.MulVec(b))
}

func TestMat4MulIdentity(t *testing.T) {
	expected := Mat4FromRows(
		Vec4{-2, 1, 2, 3},
		Vec4{3, 2, 1, -1},
		Vec4{4, 3, 6, 5},
		Vec4{1, 2, 7, 8},
	)
	assert.Equal(t, expected, expected.Mul(Mat4Ident()))
}

func TestMat4Tpose(t *testing.T) {
	expected := Mat4FromRows(
		Vec4{0, 9, 1, 0},
		Vec4{9, 8, 8, 0},
		Vec4{3, 0, 5, 5},
		Vec4{0, 8, 3, 8},
	)
	a := Mat4FromRows(
		Vec4{0, 9, 3, 0},
		Vec4{9, 8, 0, 8},
		Vec4{1, 8, 5, 3},
		Vec4{0, 0, 5, 8},
	)
	assert.Equal(t, expected, a.Tpose())
}

func TestMat4TposeIdentity(t *testing.T) {
	assert.Equal(t, Mat4Ident(), Mat4Ident().Tpose())
}

func TestMat4SubMat(t *testing.T) {
	type testData struct { i, j int; res Mat3 }
	td := []testData{
		{0, 0, Mat3FromRows(
			Vec3{2, 7, 0},
			Vec3{6, -3, 4},
			Vec3{5, 6, -7},
		)},
		{0, 1, Mat3FromRows(
			Vec3{-3, 7, 0},
			Vec3{0, -3, 4},
			Vec3{0, 6, -7},
		)},
		{3, 3, Mat3FromRows(
			Vec3{1, 5, 0},
			Vec3{-3, 2, 7},
			Vec3{0, 6, -3},
		)},
	}
	a := Mat4FromRows(
		Vec4{1, 5, 0, 6},
		Vec4{-3, 2, 7, 0},
		Vec4{0, 6, -3, 4},
		Vec4{0, 5, 6, -7},
	)
	assert := assert.New(t)
	for _, d := range td {
		assert.Equal(d.res, a.SubMat(d.i, d.j))
	}
}

func TestMat4DeleteRow(t *testing.T) {
	type testData struct { i int; res Mat3x4 }
	td := []testData{
		{0, Mat3x4{
			-3, 2, 7, 0,
			0, 6, -3, 4,
			0, 5, 6, -7,
		}},
		{1, Mat3x4{
			1, 5, 0, 6,
			0, 6, -3, 4,
			0, 5, 6, -7,
		}},
		{3, Mat3x4{
			1, 5, 0, 6,
			-3, 2, 7, 0,
			0, 6, -3, 4,
		}},
	}
	a := Mat4FromRows(
		Vec4{1, 5, 0, 6},
		Vec4{-3, 2, 7, 0},
		Vec4{0, 6, -3, 4},
		Vec4{0, 5, 6, -7},
	)
	assert := assert.New(t)
	for _, d := range td {
		assert.Equal(d.res, a.DeleteRow(d.i))
	}
}

func TestMat4Det(t *testing.T) {
	a := Mat4FromRows(
		Vec4{-2, -8, 3, 5},
		Vec4{-3, 1, 7, 3},
		Vec4{1, 2, -9, 6},
		Vec4{-6, 7, 7, -9},
	)
	assert := assert.New(t)
	assert.Equal(690.0, a.Cofactor(0, 0))
	assert.Equal(447.0, a.Cofactor(0, 1))
	assert.Equal(210.0, a.Cofactor(0, 2))
	assert.Equal(51.0, a.Cofactor(0, 3))
	assert.Equal(-4071.0, a.Det())
}

func TestMat4IsInvertible(t *testing.T) {
	type testData struct {mat Mat4; res bool}
	td := []testData{
		{ Mat4FromRows(
			Vec4{6, 4, 4, 4},
			Vec4{5, 5, 7, 6},
			Vec4{4, -9, 3, -7},
			Vec4{9, 1, 7, -6},
		), true },
		{ Mat4FromRows(
			Vec4{-4, 2, -2, -3},
			Vec4{9, 6, 2, 6},
			Vec4{0, -5, 1, -5},
			Vec4{0, 0, 0, 0},
		), false },
	}
	for _, d := range td {
		isInv, _ := d.mat.IsInvertible()
		assert.Equal(t, d.res, isInv)
	}
}

func TestMat4Inverse(t *testing.T) {
	a := Mat4FromRows(
		Vec4{-5, 2, 6, -8},
		Vec4{1, -5, 1, 8},
		Vec4{7, 7, -6, -7},
		Vec4{1, -3, 7, 4},
	)
	b := a.Inv()
	bExpected := Mat4FromRows(
		Vec4{0.21805, 0.45113, 0.24060, -0.04511},
		Vec4{-0.80827, -1.45677, -0.44361, 0.52068},
		Vec4{-0.07895, -0.22368, -0.05263, 0.19737},
		Vec4{-0.52256, -0.81391, -0.30075, 0.30639},
	)

	assert := assert.New(t)
	assert.Equal(532.0, a.Det())
	assert.Equal(-160.0, a.Cofactor(2, 3))
	assert.Equal(-160.0/532.0, b.At(3, 2))
	assert.Equal(105.0, a.Cofactor(3, 2))
	assert.Equal(105.0/532.0, b.At(2, 3))
	assert.True(bExpected.ApproxEqual(b))
}

func TestMat4MoreInverse(t *testing.T) {
	type testData struct {mat Mat4; res Mat4}
	td := []testData{
		{ Mat4FromRows(
			Vec4{8 ,-5, 9, 2},
			Vec4{7, 5, 6, 1},
			Vec4{-6, 0, 9, 6},
			Vec4{-3, 0, -9, -4},
		), Mat4FromRows(
			Vec4{-0.15385, -0.15385, -0.28205, -0.53846},
			Vec4{-0.07692, 0.12308, 0.02564, 0.03077},
			Vec4{0.35897, 0.35897, 0.43590, 0.92308},
			Vec4{-0.69231, -0.69231, -0.76923, -1.92308},
		)},
		{ Mat4FromRows(
			Vec4{9, 3, 0, 9},
			Vec4{-5, -2, -6, -3},
			Vec4{-4, 9, 6, 4},
			Vec4{-7, 6, 6, 2},
		), Mat4FromRows(
			Vec4{-0.04074, -0.07778, 0.14444, -0.22222},
			Vec4{-0.07778, 0.03333, 0.36667, -0.33333},
			Vec4{-0.02901, -0.14630, -0.10926, 0.12963},
			Vec4{0.17778, 0.06667, -0.26667, 0.33333},
		)},
	}

	for _, d := range td {
		assert.True(t, d.mat.Inv().ApproxEqual(d.res))
	}
}

func TestMat4InverseMultiplication(t *testing.T) {
	a := Mat4FromRows(
		Vec4{3, -9, 7, 3},
		Vec4{3, -8, 2, -9},
		Vec4{-4, 4, 4, 1},
		Vec4{-6, 5, -1, 1},
	)
	b := Mat4FromRows(
		Vec4{8, 2, 2, 2},
		Vec4{3, -1, 7, 0},
		Vec4{7, 0, 5, 4},
		Vec4{6, -2, 0, 5},
	)
	c := a.Mul(b)
	assert.True(t, c.Mul(b.Inv()).ApproxEqual(a))
}

// //////////////////////////// MAT3x4 //////////////////////////////
func TestMat3x4DeleteCol(t *testing.T) {
	type testData struct { i int; res Mat3 }
	td := []testData{
		{0, Mat3FromRows(
			Vec3{5, 0, 6},
			Vec3{2, 7, 0},
			Vec3{6, -3, 4},
		)},
		{1, Mat3FromRows(
			Vec3{1, 0, 6},
			Vec3{-3, 7, 0},
			Vec3{0, -3, 4},
		)},
		{3, Mat3FromRows(
			Vec3{1, 5, 0},
			Vec3{-3, 2, 7},
			Vec3{0, 6, -3},
		)},
	}
	a := Mat3x4{
		1, 5, 0, 6,
		-3, 2, 7, 0,
		0, 6, -3, 4,
	}
	assert := assert.New(t)
	for _, d := range td {
		assert.Equal(d.res, a.DeleteCol(d.i))
	}
}
