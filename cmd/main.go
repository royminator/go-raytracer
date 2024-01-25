package main

import (
	"fmt"
	"math"
	"time"

	"os"
	"roytracer/camera"
	"roytracer/gfx"
	"roytracer/light"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/pattern"
	"roytracer/render"
	"roytracer/shape"
	"roytracer/world"
	"runtime/trace"
)

var (
	CameraPos     = m.Point4(-10, 9, -25)
	LightPos      = m.Point4(8, 30, -16)
	SpherePos     = m.Vec4{-0.5, 2, 6.0, 1}
	CameraViewPos = m.Vec4{-0.5, 5, -0.5, 1}
)

func main() {
	f, _ := os.Create("trace.out")
	trace.Start(f)
	defer trace.Stop()

	w := world.World{
		Light: light.PointLight{
			Pos:       LightPos,
			Intensity: m.Vec4{1, 1, 1, 0},
		},
		Objects: arrangeObjects(),
	}
	camera := setupCamera()
	renderer := render.NewRenderer(&camera, &w, camera.Vsize, 3)
	tstart := time.Now().UTC()
	renderer.RenderParallel()
	fmt.Println("Render took:", time.Since(tstart))
	writer := gfx.PPMWriter{MaxLineLength: 70}
	writer.Write(renderer.Canvas)
	writer.SaveFile("scene.ppm")
}

func arrangeObjects() []shape.Shape {
	bigSphere := createBigSphere()
	smallSphere := createSmallSphere()
	floor := createFloor()
	cube := createCube()
	cylinder := createCylinder()
	cone := createCone()
	return []shape.Shape{
		bigSphere, smallSphere, floor, cube, cylinder, cone,
	}
}

func createBigSphere() *shape.Sphere {
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
	s.SetTf(m.Trans(SpherePos[0], 10, SpherePos[2]).Mul(m.Scale(10, 10, 10)))
	return &s
}

func createSmallSphere() *shape.Sphere {
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
	s.SetTf(m.Trans(SpherePos[0]+9, 3, SpherePos[2]-12).Mul(m.Scale(3, 3, 3)))
	return &s
}

func createCube() *shape.Cube {
	c := shape.NewCube()
	mtl := mtl.Material{
		Transparency:    0.0,
		Reflective:      0.0,
		Ambient:         0.2,
		Diffuse:         0.7,
		Color:           m.Vec4{0.8, 0.8, 0.4},
		Specular:        0.3,
		Shininess:       50,
		RefractiveIndex: 0.0,
	}
	c.SetMat(mtl)
	tf := m.Trans(SpherePos[0]-4, 2, SpherePos[2]-13)
	tf = tf.Mul(m.Scale(2, 2, 2))
	tf = tf.Mul(m.RotY(math.Pi/3.2))
	c.SetTf(tf)
	return &c
}

func createCylinder() *shape.Cylinder {
	c := shape.NewCylinder()
	c.Min = 0
	c.Max = 1
	c.Closed = true

	mtl := mtl.Material{
		Transparency:    0.0,
		Reflective:      0.1,
		Ambient:         0.2,
		Diffuse:         0.7,
		Color:           m.Vec4{0.73, 0.2, 0.36},
		Specular:        0.3,
		Shininess:       300,
		RefractiveIndex: 1.0,
	}
	c.SetMat(mtl)
	tf := m.Trans(SpherePos[0]-11, 0, SpherePos[2]-8)
	tf = tf.Mul(m.Scale(3, 8, 3))
	c.SetTf(tf)
	return &c
}

func createCone() *shape.Cone {
	c := shape.NewCone()
	c.Min = 0
	c.Max = 2
	c.Closed = true

	mtl := mtl.Material{
		Transparency:    0.0,
		Reflective:      0.0,
		Ambient:         0.2,
		Diffuse:         0.7,
		Color:           m.Vec4{0.2, 0.73, 0.73},
		Specular:        0.3,
		Shininess:       50,
		RefractiveIndex: 0.0,
	}
	c.SetMat(mtl)
	tf := m.Trans(SpherePos[0]+3, 10, SpherePos[2]-13)
	tf = tf.Mul(m.Scale(2, 5, 2))
	tf = tf.Mul(m.RotZ(math.Pi))
	c.SetTf(tf)
	return &c
}

func createFloor() *shape.Plane {
	c1 := m.Vec4{0.878, 0.525, 0.275}
	c2 := m.Vec4{0.357, 0.478, 0.71}
	pattern := pattern.NewRingPattern(c1, c2)

	floor := shape.NewPlane()
	floor.SetTf(m.Scale(3, 0.01, 3))
	floor.O.Material = mtl.DefaultMaterial()
	floor.O.Material.Color = m.Vec4{1, 0.9, 0.9, 0}
	floor.O.Material.Specular = 0
	floor.O.Material.Pattern = &pattern
	floor.O.Material.Shininess = 300
	floor.O.Material.Reflective = 0.8
	return &floor
}

func setupCamera() camera.Camera {
	camera := camera.NewCamera(1400, 1200, math.Pi/2.0)
	from := CameraPos
	to := CameraViewPos
	up := m.Vector4(0, 1, 0)
	camera.SetTf(m.View(from, to, up))
	return camera
}
