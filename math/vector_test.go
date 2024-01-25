package math

import (
	"math"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreatPointWShouldBe1(t *testing.T) {
	actual := Point4(4, -4, 3)
	expected := Vec4{4, -4, 3, 1}
	assert.Equal(t, expected, actual)
}

func TestCreateVectorWShouldBe0(t *testing.T) {
	actual := Vector4(4, -4, 3)
	expected := Vec4{4, -4, 3, 0}
	assert.Equal(t, expected, actual)
}

func TestAddVec4(t *testing.T) {
	a1 := Vec4{3, -2, 5, 1}
	a2 := Vec4{-2, 3, 1, 0}
	expected := Vec4{1, 1, 6, 1}
	actual := a1
	actual.Add(a2)
	assert.Equal(t, expected, actual)
}

func TestSubPoints(t *testing.T) {
	a1 := Point4(3, 2, 1)
	a2 := Point4(5, 6, 7)
	expected := Vector4(-2, -4, -6)
	actual := a1
	actual.Sub(a2)
	assert.Equal(t, expected, actual)
}

func TestSubVecFromPoint(t *testing.T) {
	a1 := Point4(3, 2, 1)
	a2 := Vector4(5, 6, 7)
	expected := Point4(-2, -4, -6)
	actual := a1
	actual.Sub(a2)
	assert.Equal(t, expected, actual)
}

func TestSubVectors(t *testing.T) {
	a1 := Vector4(3, 2, 1)
	a2 := Vector4(5, 6, 7)
	expected := Vector4(-2, -4, -6)
	actual := a1
	actual.Sub(a2)
	assert.Equal(t, expected, actual)
}

func TestSubFromZeroVec(t *testing.T) {
	zero := Vec4With(0)
	v := Vector4(1, -2, 3)
	expected := Vector4(-1, 2, -3)
	actual := zero
	actual.Sub(v)
	assert.Equal(t, expected, actual)
}

func TestNegate(t *testing.T) {
	expected := Vec4{1, -2, 3, -4}
	actual := (&Vec4{-1, 2, -3, 4}).Negate()
	assert.Equal(t, expected, actual)
}

func TestScalarMul(t *testing.T) {
	expected := Vec4{3.5, -7, 10.5, -14}
	actual := (&Vec4{1, -2, 3, -4}).Mul(3.5)
	assert.Equal(t, expected, actual)
}

func TestScalarMulFrac(t *testing.T) {
	expected := Vec4{0.5, -1, 1.5, -2}
	actual := Vec4{1, -2, 3, -4}.Mul(0.5)
	assert.Equal(t, expected, actual)
}

func TestScalarDiv(t *testing.T) {
	expected := Vec4{0.5, -1, 1.5, -2}
	actual := Vec4{1, -2, 3, -4}.Div(2)
	assert.Equal(t, expected, actual)
}

func TestMagn(t *testing.T) {
	type testCase struct {vec Vec4; res float64}
	cases := []testCase{
		{Vector4(1, 0, 0), 1},
		{Vector4(0, 1, 0), 1},
		{Vector4(0, 0, 1), 1},
		{Vector4(1, 2, 3), math.Sqrt(14)},
		{Vector4(-1, -2, -3), math.Sqrt(14)},
	}
	for _, c := range cases {
		assert.Equal(t, c.vec.Magnitude(), c.res)
	}
}

func TestVec4Normalize(t *testing.T) {
	type testCase struct {vec Vec4; res Vec4}
	cases := []testCase{
		{Vector4(4, 0, 0), Vector4(1, 0, 0)},
		{Vector4(1, 2, 3), Vector4(0.26726, 0.53452, 0.80178)},
	}

	for _, c := range cases {
		actual := c.vec.Normalize()
		assert.True(t, actual.ApproxEqual(c.res))
	}
}

func TestVec4MagnitudeOfNormalizedVecIs1(t *testing.T) {
	actual := Vector4(1, 2, 3).Normalize().Magnitude()
	assert.Equal(t, 1.0, actual)
}

func TestDotProduct(t *testing.T) {
	v := Vector4(1, 2, 3)
	assert.Equal(t, 20.0, v.Dot(Vector4(2, 3, 4)))
}

func TestCrossProduct(t *testing.T) {
	a := Vector4(1, 2, 3)
	b := Vector4(2, 3, 4)
	assert.Equal(t, Vector4(-1, 2, -1), a.Cross(b))
	assert.Equal(t, Vector4(1, -2, 1), b.Cross(a))
}

func TestMulElem(t *testing.T) {
	expected := Vector4(0.9, 0.2, 0.04)
	actual := Vector4(1, 0.2, 0.4).MulElem(Vector4(0.9, 1, 0.1))
	assert.True(t, expected.ApproxEqual(actual))
}

func TestReflectVec4ApproachingAt45(t *testing.T) {
	expected := Vector4(1, 1, 0)
	v := Vector4(1, -1, 0)
	n := Vector4(0, 1, 0)
	assert.Equal(t, expected, v.Reflect(n))
}

func TestReflectVec4OffSlantedSurface(t *testing.T) {
	expected := Vector4(1, 0, 0)
	v := Vector4(0, -1, 0)
	n := Vector4(math.Sqrt2/2.0, math.Sqrt2/2.0, 0)
	assert.True(t, expected.ApproxEqual(v.Reflect(n)))
}

func TestVectorMulMat(t *testing.T) {
	expected := Vec4{18, 24, 33, 1}
	a := Mat4FromRows(
		Vec4{1, 2, 3, 4},
		Vec4{2, 4, 4, 2},
		Vec4{8, 6, 4, 1},
		Vec4{0, 0, 0, 1},
	)
	b := Vec4{1, 2, 3, 1}
	b.MulMat(a)
	assert.Equal(t, expected, b)
}
