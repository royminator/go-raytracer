package pattern

import (
	"testing"

	"roytracer/color"
	m "roytracer/math"

	"github.com/stretchr/testify/assert"
)

func TestCreateStripePattern(t *testing.T) {
	p := NewStripePattern(color.White, color.Black)
	assert.Equal(t, color.White, p.PB.A)
	assert.Equal(t, color.Black, p.PB.B)
}

func TestStripePatternIsConstantInY(t *testing.T) {
	p := NewStripePattern(color.White, color.Black)
	assert := assert.New(t)
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 0)))
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 1, 0)))
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 2, 0)))
}

func TestStripePatternIsConstantInZ(t *testing.T) {
	p := NewStripePattern(color.White, color.Black)
	assert := assert.New(t)
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 0)))
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 1)))
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 2)))
}

func TestStripePatternAlternatesInX(t *testing.T) {
	p := NewStripePattern(color.White, color.Black)
	assert := assert.New(t)
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 0)))
	assert.Equal(color.White, p.SampleAt(m.Point4(0.9, 0, 0)))
	assert.Equal(color.Black, p.SampleAt(m.Point4(1, 0, 0)))
	assert.Equal(color.Black, p.SampleAt(m.Point4(-0.1, 0, 0)))
	assert.Equal(color.Black, p.SampleAt(m.Point4(-1, 0, 0)))
	assert.Equal(color.White, p.SampleAt(m.Point4(-1.1, 0, 0)))
}

func TestPatternGetTf(t *testing.T) {
	p := NewStripePattern(color.White, color.Black)
	assert.Equal(t, m.Mat4Ident(), p.GetTf())
}

func TestPatternSetTf(t *testing.T) {
	p := NewStripePattern(color.White, color.Black)
	tf := m.Trans(1, 2, 3)
	p.SetTf(tf)
	assert.Equal(t, tf, p.GetTf())
}

func TestGradientPatternLinearlyInterpolatesColors(t *testing.T) {
	p := NewGradientPattern(color.White, color.Black)
	assert := assert.New(t)
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 0)))
	assert.Equal(m.Vec4{0.75, 0.75, 0.75}, p.SampleAt(m.Point4(0.25, 0, 0)))
	assert.Equal(m.Vec4{0.5, 0.5, 0.5}, p.SampleAt(m.Point4(0.5, 0, 0)))
	assert.Equal(m.Vec4{0.25, 0.25, 0.25}, p.SampleAt(m.Point4(0.75, 0, 0)))
}

func TestRingPatternShouldExtendInXAndZ(t *testing.T) {
	p := NewRingPattern(color.White, color.Black)
	assert := assert.New(t)
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 0)))
	assert.Equal(color.Black, p.SampleAt(m.Point4(1, 0, 0)))
	assert.Equal(color.Black, p.SampleAt(m.Point4(0, 0, 1)))
	assert.Equal(color.Black, p.SampleAt(m.Point4(0.708, 0, 0.708)))
}

func TestCheckersPatternShouldRepeatInX(t *testing.T) {
	p := NewCheckersPattern(color.White, color.Black)
	assert := assert.New(t)
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 0)))
	assert.Equal(color.White, p.SampleAt(m.Point4(0.99, 0, 0)))
	assert.Equal(color.Black, p.SampleAt(m.Point4(1.01, 0, 0)))
}

func TestCheckersPatternShouldRepeatInY(t *testing.T) {
	p := NewCheckersPattern(color.White, color.Black)
	assert := assert.New(t)
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 0)))
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0.99, 0)))
	assert.Equal(color.Black, p.SampleAt(m.Point4(0, 1.01, 0)))
}

func TestCheckersPatternShouldRepeatInZ(t *testing.T) {
	p := NewCheckersPattern(color.White, color.Black)
	assert := assert.New(t)
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 0)))
	assert.Equal(color.White, p.SampleAt(m.Point4(0, 0, 0.99)))
	assert.Equal(color.Black, p.SampleAt(m.Point4(0, 0, 1.01)))
}
