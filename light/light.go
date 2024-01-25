package light

import (
	"math"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/shape"
)

type (
	PointLight struct {
		Pos   m.Vec4
		Intensity m.Vec4
	}
)

var (
	ColorBlack = m.Vec4{0, 0, 0, 0}
)

func Lighting(mat mtl.Material, shape shape.Shape, light PointLight, pos, eye, normal m.Vec4, inShadow bool) m.Vec4 {
	color := mat.Color
	if shape.GetMat().Pattern != nil {
		color = shape.SamplePatternAt(pos)
	}

	effColor := color.MulElem(light.Intensity)
	ambient := effColor.Mul(mat.Ambient)

	if inShadow {
		return ambient
	}

	lightv := light.Pos
	lightv.Sub(pos)
	lightv = lightv.Normalize()
	lightDotNorm := lightv.Dot(normal)
	var diffuse, specular m.Vec4

	if lightDotNorm < 0.0 {
		diffuse = ColorBlack
		specular = ColorBlack 
	} else {
		diffuse = effColor.Mul(mat.Diffuse*lightDotNorm)
		reflectv := lightv.Negate().Reflect(normal)
		reflectDotEye := reflectv.Dot(eye)

		if reflectDotEye <= 0.0 {
			specular = ColorBlack
		} else {
			factor := math.Pow(reflectDotEye, mat.Shininess)
			specular = light.Intensity.Mul(mat.Specular*factor)
		}
	}
	res := ambient
	res.Add(diffuse)
	res.Add(specular)
	return res
}
