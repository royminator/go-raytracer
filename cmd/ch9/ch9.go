package main

import (
	"math"
	"roytracer/camera"
	"roytracer/gfx"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/shape"
	"roytracer/world"
	"roytracer/light"
)

func main() {
	w := world.World{
		Light: light.PointLight{
			Pos: m.Point4(-10, 10, -10),
			Intensity: m.Vec4{1, 1, 1, 0},
		},
		Objects: arrangeObjects(),
	}
	camera := camera.NewCamera(600, 300, math.Pi/3.0)
	from := m.Point4(0, 1.5, -5)
	to := m.Point4(0, 1, 0)
	up := m.Vector4(0, 1, 0)
	camera.SetTf(m.View(from, to ,up))
	image := camera.Render(&w)
	writer := gfx.PPMWriter{MaxLineLength: 70}
	writer.Write(image)
	writer.SaveFile("scene.ppm")
}

func arrangeObjects() []shape.Shape {
	floor := createFloor()
	middle := createMiddleSphere()
	right := createRightSphere()
	left := createLeftSphere()
	return []shape.Shape{
		floor, middle, left, right,
	}
}

func createFloor() *shape.Sphere {
	floor := shape.NewSphere()
	floor.SetTf(m.Scale(10, 0.01, 10))
	floor.O.Material = mtl.DefaultMaterial()
	floor.O.Material.Color = m.Vec4{1, 0.9, 0.9, 0}
	floor.O.Material.Specular = 0
	return &floor
}

func createMiddleSphere() *shape.Sphere {
	s := shape.NewSphere()
	s.SetTf(m.Trans(-0.5, 1, 0.5))
	s.O.Material = mtl.DefaultMaterial()
	s.O.Material.Color = m.Vec4{0.1, 1, 0.5, 0}
	s.O.Material.Diffuse = 0.7
	s.O.Material.Specular = 0.3
	return &s
}

func createRightSphere() *shape.Sphere {
	s := shape.NewSphere()
	tf := m.Trans(1.5, 0.5, -0.5)
	tf = tf.Mul(m.Scale(0.5, 0.5, 0.5))
	s.SetTf(tf)
	s.O.Material = mtl.DefaultMaterial()
	s.O.Material.Color = m.Vec4{0.5, 1, 0.1, 0}
	s.O.Material.Diffuse = 0.7
	s.O.Material.Specular = 0.3
	return &s
}

func createLeftSphere() *shape.Sphere {
	s := shape.NewSphere()
	tf := m.Trans(-1.5, 0.33, -0.75)
	tf = tf.Mul(m.Scale(0.33, 0.33, 0.33))
	s.SetTf(tf)
	s.O.Material = mtl.DefaultMaterial()
	s.O.Material.Color = m.Vec4{1.0, 0.8, 0.1, 0}
	s.O.Material.Diffuse = 0.7
	s.O.Material.Specular = 0.3
	return &s
}
