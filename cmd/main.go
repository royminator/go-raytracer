package main

import (
	"fmt"
	"math"
	"roytracer/camera"
	"roytracer/gfx"
	"roytracer/light"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/pattern"
	"roytracer/shape"
	"roytracer/world"
	"time"
)

func main() {
	w := world.World{
		Light: light.PointLight{
			Pos: m.Point4(-10, 10, -10),
			Intensity: m.Vec4{1, 1, 1, 0},
		},
		Objects: arrangeObjects(),
	}
	camera := setupCamera()
	tstart := time.Now().UTC()
	image := camera.Render(&w)
	fmt.Println("Render took:", time.Since(tstart))
	writer := gfx.PPMWriter{MaxLineLength: 70}
	writer.Write(image)
	writer.SaveFile("scene.ppm")
}

func arrangeObjects() []shape.Shape {
	middle := createMiddleSphere()
	floor := createFloor()
	return []shape.Shape{
		middle, floor,
	}
}

func createMiddleSphere() *shape.Sphere {
	s := shape.NewSphere()
	s.SetTf(m.Trans(-0.5, 2, -0.5).Mul(m.Scale(2, 2, 2)))
	c1 := m.Vec4{0.812, 0.376, 0.702}
	c2 := m.Vec4{0.835, 0.922, 0.09, 0.349}
	p := pattern.NewStripePattern(c1, c2)
	s.O.Material.Color = m.Vec4{0.1, 1, 0.5, 0}
	s.O.Material.Diffuse = 0.7
	s.O.Material.Specular = 0.3
	s.O.Material.Reflective = 0.5
	s.O.Material.Shininess = 50
	s.SetPattern(&p)
	return &s
}

func createFloor() *shape.Plane{
	c1 := m.Vec4{0.878, 0.525, 0.275}
	c2 := m.Vec4{0.357, 0.478, 0.71}
	pattern := pattern.NewCheckersPattern(c1, c2)

	floor := shape.NewPlane()
	floor.SetTf(m.Scale(10, 0.01, 10))
	floor.O.Material = mtl.DefaultMaterial()
	floor.O.Material.Color = m.Vec4{1, 0.9, 0.9, 0}
	floor.O.Material.Specular = 0
	floor.O.Material.Pattern = &pattern
	return &floor
}

func setupCamera() camera.Camera {
	camera := camera.NewCamera(800, 400, math.Pi/2.0)
	from := m.Point4(0, 4.5, -6)
	to := m.Point4(0, 1, 0)
	up := m.Vector4(0, 1, 0)
	camera.SetTf(m.View(from, to ,up))
	return camera
}
