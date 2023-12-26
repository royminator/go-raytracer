package mtl

import m "roytracer/math"

type Material struct {
	Color     m.Vec4
	Ambient   float64
	Diffuse   float64
	Specular  float64
	Shininess float64
}

func DefaultMaterial() Material {
	return Material{
		Color:     m.Color4(1, 1, 1, 0),
		Ambient:   0.1,
		Diffuse:   0.9,
		Specular:  0.9,
		Shininess: 200.0,
	}
}
