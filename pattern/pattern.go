package pattern

import (
	"math"

	m "roytracer/math"
)

type (
	Pattern interface {
		SampleAt(m.Vec4) m.Vec4
		SetTf(m.Mat4)
		GetTf() m.Mat4
		GetInvTf() m.Mat4
	}

	PatternBase struct {
		A, B      m.Vec4
		Tf, InvTf m.Mat4
	}

	StripePattern struct {
		PB PatternBase
	}

	GradientPattern struct {
		PB PatternBase
	}

	RingPattern struct {
		PB PatternBase
	}

	CheckersPattern struct {
		PB PatternBase
	}

	TestPattern struct {
		Tf, InvTf m.Mat4
	}
)

func defaultPatternBase(a, b m.Vec4) PatternBase {
	return PatternBase{
		A:     a,
		B:     b,
		Tf:    m.Mat4Ident(),
		InvTf: m.Mat4Ident(),
	}
}

func NewStripePattern(a, b m.Vec4) StripePattern {
	return StripePattern{defaultPatternBase(a, b)}
}

func (p *StripePattern) SampleAt(point m.Vec4) m.Vec4 {
	if int(math.Floor(point[0]))%2 == 0 {
		return p.PB.A
	}
	return p.PB.B
}

func (p *StripePattern) SetTf(tf m.Mat4) {
	p.PB.Tf = tf
	p.PB.InvTf = tf.Inv()
}

func (p *StripePattern) GetInvTf() m.Mat4 {
	return p.PB.InvTf
}

func (p *StripePattern) GetTf() m.Mat4 {
	return p.PB.Tf
}

func NewGradientPattern(a, b m.Vec4) GradientPattern {
	return GradientPattern{defaultPatternBase(a, b)}
}

func (p *GradientPattern) SampleAt(point m.Vec4) m.Vec4 {
	dist := p.PB.B
	dist.Sub(p.PB.A)
	frac := point[0] - math.Floor(point[0])
	res := p.PB.A
	res.Add(dist.Mul(frac))
	return res
}

func (p *GradientPattern) SetTf(tf m.Mat4) {
	p.PB.Tf = tf
	p.PB.InvTf = tf.Inv()
}

func (p *GradientPattern) GetTf() m.Mat4 {
	return p.PB.Tf
}

func (p *GradientPattern) GetInvTf() m.Mat4 {
	return p.PB.InvTf
}

func NewRingPattern(a, b m.Vec4) RingPattern {
	return RingPattern{defaultPatternBase(a, b)}
}

func (p *RingPattern) SampleAt(point m.Vec4) m.Vec4 {
	x := point[0]*point[0] + point[2]*point[2]
	if int(math.Floor(math.Sqrt(x)))%2 == 0 {
		return p.PB.A
	}
	return p.PB.B
}

func (p *RingPattern) SetTf(tf m.Mat4) {
	p.PB.Tf = tf
	p.PB.InvTf = tf.Inv()
}

func (p *RingPattern) GetTf() m.Mat4 {
	return p.PB.Tf
}

func (p *RingPattern) GetInvTf() m.Mat4 {
	return p.PB.InvTf
}

func NewCheckersPattern(a, b m.Vec4) CheckersPattern {
	return CheckersPattern{defaultPatternBase(a, b)}
}

func (p *CheckersPattern) SampleAt(point m.Vec4) m.Vec4 {
	x := math.Floor(point[0]) + math.Floor(point[1]) + math.Floor(point[2])
	if int(x)%2 == 0 {
		return p.PB.A
	}
	return p.PB.B
}

func (p *CheckersPattern) SetTf(tf m.Mat4) {
	p.PB.Tf = tf
	p.PB.InvTf = tf.Inv()
}

func (p *CheckersPattern) GetTf() m.Mat4 {
	return p.PB.Tf
}

func (p *CheckersPattern) GetInvTf() m.Mat4 {
	return p.PB.InvTf
}

func NewTestPattern() TestPattern {
	return TestPattern{
		Tf:    m.Mat4Ident(),
		InvTf: m.Mat4Ident(),
	}
}

func (tp *TestPattern) SampleAt(p m.Vec4) m.Vec4 {
	return m.Vec4{p[0], p[1], p[2]}
}

func (tp *TestPattern) SetTf(tf m.Mat4) {
	tp.Tf = tf
	tp.InvTf = tf.Inv()
}

func (tp *TestPattern) GetInvTf() m.Mat4 {
	return tp.InvTf
}

func (tp *TestPattern) GetTf() m.Mat4 {
	return tp.Tf
}

