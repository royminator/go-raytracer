package math

import (
	"math"
	"testing"
)

func TestCreatPointWShouldBe1(t *testing.T) {
	actual := Point4(4, -4, 3)
	expected := Vec4{4, -4, 3, 1}
	checkResult(expected, actual, t)
}

func TestCreateVectorWShouldBe0(t *testing.T) {
	actual := Vector4(4, -4, 3)
	expected := Vec4{4, -4, 3, 0}
	checkResult(expected, actual, t)
}

func TestAddVec4(t *testing.T) {
	a1 := Vec4{3, -2, 5, 1}
	a2 := Vec4{-2, 3, 1, 0}
	expected := Vec4{1, 1, 6, 1}
	actual := a1.Add(a2)
	checkResult(expected, actual, t)
}

func TestSubPoints(t *testing.T) {
	a1 := Point4(3, 2, 1)
	a2 := Point4(5, 6, 7)
	expected := Vector4(-2, -4, -6)
	actual := a1.Sub(a2)
	checkResult(expected, actual, t)
}

func TestSubVecFromPoint(t *testing.T) {
	a1 := Point4(3, 2, 1)
	a2 := Vector4(5, 6, 7)
	expected := Point4(-2, -4, -6)
	actual := a1.Sub(a2)
	checkResult(expected, actual, t)
}

func TestSubVectors(t *testing.T) {
	a1 := Vector4(3, 2, 1)
	a2 := Vector4(5, 6, 7)
	expected := Vector4(-2, -4, -6)
	actual := a1.Sub(a2)
	checkResult(expected, actual, t)
}

func TestSubFromZeroVec(t *testing.T) {
	zero := Vec4With(0)
	v := Vector4(1, -2, 3)
	expected := Vector4(-1, 2, -3)
	actual := zero.Sub(v)
	checkResult(expected, actual, t)
}

func TestNegate(t *testing.T) {
	expected := Vec4{1, -2, 3, -4}
	actual := (&Vec4{-1, 2, -3, 4}).Negate()
	checkResult(expected, actual, t)
}

func TestScalarMul(t *testing.T) {
	expected := Vec4{3.5, -7, 10.5, -14}
	actual := (&Vec4{1, -2, 3, -4}).Mul(3.5)
	checkResult(expected, actual, t)
}

func TestScalarMulFrac(t *testing.T) {
	expected := Vec4{0.5, -1, 1.5, -2}
	actual := Vec4{1, -2, 3, -4}.Mul(0.5)
	checkResult(expected, actual, t)
}

func TestScalarDiv(t *testing.T) {
	expected := Vec4{0.5, -1, 1.5, -2}
	actual := Vec4{1, -2, 3, -4}.Div(2)
	checkResult(expected, actual, t)
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
		checkResult(c.res, c.vec.Magnitude(), t)
	}
}

func TestNormalize(t *testing.T) {
	type testCase struct {vec Vec4; res Vec4}
	cases := []testCase{
		{Vector4(4, 0, 0), Vector4(1, 0, 0)},
		{Vector4(1, 2, 3), Vector4(0.26726, 0.53452, 0.80178)},
	}

	for _, c := range cases {
		actual := c.vec.Normalize()
		if !actual.Approx(c.res) {
			t.Errorf("Expected approx %v, got %v", c.res, actual)
		}
	}
}

func TestMagnitudeOfNormalizedVecIs1(t *testing.T) {
	norm := Vec4{1, 2, 3, 0}.Normalize()
	actual := norm.Magnitude()
	if actual != 1.0 {
		t.Errorf("Expected %f, got %f", 1.0, actual) 
	}
}

/////////////////////// HELPERS /////////////////////// 
func checkResult(expected, actual interface{}, t *testing.T) {
	if expected != actual {
		t.Errorf("Expected %v, got %v", expected, actual)
	}
}
