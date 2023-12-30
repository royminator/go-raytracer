package mtl

import (
	"math"
	m "roytracer/math"
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

func Lighting(mat Material, light PointLight, pos, eye, normal m.Vec4, inShadow bool) m.Vec4 {
	effColor := mat.Color.MulElem(light.Intensity)
	ambient := effColor.Mul(mat.Ambient)

	if inShadow {
		return ambient
	}

	lightv := light.Pos.Sub(pos).Normalize()
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
	return ambient.Add(diffuse).Add(specular)
}
