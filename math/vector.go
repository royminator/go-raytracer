package math

type (
	Vec4 struct { X, Y, Z, W float64 }
)

func Point4(x, y, z float64) Vec4 {
	return Vec4{x, y, z, 1.0}
}

func Vector4(x, y, z float64) Vec4 {
	return Vec4{x, y, z, 0.0}
}
