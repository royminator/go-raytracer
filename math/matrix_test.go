package math

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

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
	actual := m1.Mat4Add(m2)
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
