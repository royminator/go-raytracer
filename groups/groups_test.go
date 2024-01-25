package groups

import (
	"testing"

	m "roytracer/math"
	"roytracer/shape"

	"github.com/stretchr/testify/assert"
)

func TestCreateGroup(t *testing.T) {
	g := NewGroup()
	assert.Equal(t, m.Mat4Ident(), g.Tf)
	assert.Empty(t, g.Shapes)
}

func TestAddChildToGroup(t *testing.T) {
	g := NewGroup()
	s := shape.NewTestShape()
	g.AddChild(s)
	assert.Len(t, g.Shapes, 1)
	assert.True(t, g.Contains(s))
}

func TestIntersectRayWithEmptyGroup(t *testing.T) {
	g := NewGroup()
	r := shape.Ray{
		Origin: m.Point4(0, 0, 0),
		Dir:    m.Vector4(0, 0, 1),
	}
	xs := g.localIntersect(r)
	assert.Empty(t, xs)
}

func TestIntersectRayWithNonEmptyGroup(t *testing.T) {
	g := NewGroup()
	s1 := shape.NewSphere()
	s2 := shape.NewSphere()
	s2.SetTf(m.Trans(0, 0, -3))
	s3 := shape.NewSphere()
	g.AddChild(&s1)
	g.AddChild(&s2)
	g.AddChild(&s3)
	r := shape.Ray{
		Origin: m.Point4(0, 0, -5),
		Dir:    m.Vector4(0, 0, 1),
	}

	xs := g.localIntersect(r)

	assert := assert.New(t)
	assert.Equal(&s2, xs[0].S)
	assert.Equal(&s2, xs[1].S)
	assert.Equal(&s1, xs[2].S)
	assert.Equal(&s1, xs[3].S)
}
