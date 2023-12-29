package main

import (
	"math"
	"roytracer/camera"
	"roytracer/gfx"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/shape"
	"roytracer/world"
)

func main() {
	w := world.World{
		Light: mtl.PointLight{
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

func arrangeObjects() []*shape.Sphere {
	floor := createFloor()
	leftWall := createLeftWall(floor.Material)
	rightWall := createRightWall(floor.Material)
	middle := createMiddleSphere()
	right := createRightSphere()
	left := createLeftSphere()
	return []*shape.Sphere{
		floor, leftWall, rightWall, middle, left, right,
	}
}

func createFloor() *shape.Sphere {
	floor := shape.NewSphere()
	floor.SetTf(m.Scale(m.Vec3{10, 0.01, 10}))
	floor.Material = mtl.DefaultMaterial()
	floor.Material.Color = m.Vec4{1, 0.9, 0.9, 0}
	floor.Material.Specular = 0
	return &floor
}

func createLeftWall(mat mtl.Material) *shape.Sphere {
	left := shape.NewSphere()
	tf := m.Trans(m.Vec3{0, 0, 5})
	tf = tf.Mul(m.RotY(-math.Pi/4.0))
	tf = tf.Mul(m.RotX(math.Pi/2.0))
	tf = tf.Mul(m.Scale(m.Vec3{10, 0.01, 10}))

	left.SetTf(tf)
	left.Material = mat
	return &left
}

func createRightWall(mat mtl.Material) *shape.Sphere {
	wall := shape.NewSphere()
	tf := m.Trans(m.Vec3{0, 0, 5})
	tf = tf.Mul(m.RotY(math.Pi/4.0))
	tf = tf.Mul(m.RotX(math.Pi/2.0))
	tf = tf.Mul(m.Scale(m.Vec3{10, 0.01, 10}))

	wall.SetTf(tf)
	wall.Material = mat
	return &wall
}

func createMiddleSphere() *shape.Sphere {
	s := shape.NewSphere()
	s.SetTf(m.Trans(m.Vec3{-0.5, 1, 0.5}))
	s.Material = mtl.DefaultMaterial()
	s.Material.Color = m.Vec4{0.1, 1, 0.5, 0}
	s.Material.Diffuse = 0.7
	s.Material.Specular = 0.3
	return &s
}

func createRightSphere() *shape.Sphere {
	s := shape.NewSphere()
	tf := m.Trans(m.Vec3{1.5, 1, -0.5})
	tf = tf.Mul(m.Scale(m.Vec3{0.5, 0.5, 0.5}))
	s.SetTf(tf)
	s.Material = mtl.DefaultMaterial()
	s.Material.Color = m.Vec4{0.5, 1, 0.1, 0}
	s.Material.Diffuse = 0.7
	s.Material.Specular = 0.3
	return &s
}

func createLeftSphere() *shape.Sphere {
	s := shape.NewSphere()
	tf := m.Trans(m.Vec3{-1.5, 0.33, -0.75})
	tf = tf.Mul(m.Scale(m.Vec3{0.33, 0.33, 0.33}))
	s.SetTf(tf)
	s.Material = mtl.DefaultMaterial()
	s.Material.Color = m.Vec4{1.0, 0.8, 0.1, 0}
	s.Material.Diffuse = 0.7
	s.Material.Specular = 0.3
	return &s
}
