package math

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestCreateMatrix(t *testing.T) {
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
