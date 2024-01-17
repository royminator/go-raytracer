package main

import (
	"fmt"
	"math"
	"time"

	"roytracer/camera"
	"roytracer/gfx"
	"roytracer/light"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/pattern"
	"roytracer/shape"
	"roytracer/world"
)

var (
	CameraPos     = m.Point4(-12, 17, -12)
	LightPos      = m.Point4(8, 30, -16)
	SpherePos     = m.Vec4{-0.5, 2, 6.0, 1}
	// SpherePos     = m.Vec4{CameraPos[0], 2, CameraPos[2]}
	CameraViewPos = m.Vec4{-0.5, 2, -0.5, 1}
)

func main() {
	w := world.World{
		Light: light.PointLight{
			Pos:       LightPos,
			Intensity: m.Vec4{1, 1, 1, 0},
		},
		Objects: arrangeObjects(),
	}
	camera := setupCamera()
	tstart := time.Now().UTC()
	image := camera.Render(&w, 5)
	fmt.Println("Render took:", time.Since(tstart))
	writer := gfx.PPMWriter{MaxLineLength: 70}
	writer.Write(image)
	writer.SaveFile("scene.ppm")
}

func arrangeObjects() []shape.Shape {
	middle := createMiddleSphere()
	other := createOtherSphere()
	floor := createFloor()
	return []shape.Shape{
		middle, other, floor,
	}
}

func createMiddleSphere() *shape.Sphere {
	s := shape.NewSphere()
	mtl := mtl.Material{
		Transparency:    1.0,
		Reflective:      0.9,
		Ambient:         0.1,
		Diffuse:         0.1,
		Color:           m.Vec4{0.3, 0.1, 0.6},
		Specular:        1.0,
		Shininess:       400,
		RefractiveIndex: 2.0,
	}
	s.SetMat(mtl)
	s.SetTf(m.Trans(SpherePos[0], SpherePos[1], SpherePos[2]).Mul(m.Scale(10, 10, 10)))
	return &s
}

func createOtherSphere() *shape.Sphere {
	s := shape.NewSphere()
	mtl := mtl.Material{
		Transparency:    0.0,
		Reflective:      0.3,
		Ambient:         0.3,
		Diffuse:         0.4,
		Color:           m.Vec4{0.2, 0.7, 0.4},
		Specular:        1.0,
		Shininess:       300,
		RefractiveIndex: 0.9,
	}
	s.SetMat(mtl)
	s.SetTf(m.Trans(SpherePos[0]+2, SpherePos[1], SpherePos[2]-15).Mul(m.Scale(3, 3, 3)))
	return &s
}

func createFloor() *shape.Plane {
	c1 := m.Vec4{0.878, 0.525, 0.275}
	c2 := m.Vec4{0.357, 0.478, 0.71}
	pattern := pattern.NewCheckersPattern(c1, c2)

	floor := shape.NewPlane()
	floor.SetTf(m.Scale(5, 0.01, 5))
	floor.O.Material = mtl.DefaultMaterial()
	floor.O.Material.Color = m.Vec4{1, 0.9, 0.9, 0}
	floor.O.Material.Specular = 0
	floor.O.Material.Pattern = &pattern
	return &floor
}

func setupCamera() camera.Camera {
	camera := camera.NewCamera(1200, 700, math.Pi/2.0)
	from := CameraPos
	to := CameraViewPos
	up := m.Vector4(0, 1, 0)
	camera.SetTf(m.View(from, to, up))
	return camera
}
