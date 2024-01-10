package mtl

import (
	"testing"

	m "roytracer/math"

	"github.com/stretchr/testify/assert"
)

func TestDefaultMaterial(t *testing.T) {
	material := DefaultMaterial()

	assert := assert.New(t)
	assert.Equal(m.Color4(1, 1, 1, 0), material.Color)
	assert.Equal(0.1, material.Ambient)
	assert.Equal(0.9, material.Diffuse)
	assert.Equal(0.9, material.Specular)
	assert.Equal(200.0, material.Shininess)
	assert.Equal(0.0, material.Reflective)
}
