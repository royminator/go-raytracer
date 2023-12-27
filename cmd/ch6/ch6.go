package main

import (
	. "roytracer/gfx"
	. "roytracer/math"
	. "roytracer/ray"
	"roytracer/mtl"
	"time"
)

var (
	SpherePos = Vec3{0, 0, -1}
	ViewPos = Vec4{0, 0, 5, 1}
	SphereRadius = 2.0
	WallWidth = 5.0
	WallHeight = 5.0
	CanvasWidth uint32 = 900
	CanvasHeight uint32 = 900
	Light = mtl.PointLight{Pos: Point4(-10, -10, 10), Intensity: Color4(1, 1, 1, 0)}
)

func main() {
	sphere := NewSphere()
	sphere.Material.Color = Vec4{1, 0.2, 1, 0}
	sphere.Tf = Trans(SpherePos).Mul(Scale(Vec3{SphereRadius, SphereRadius, SphereRadius}))
	canvas := NewCanvas(CanvasWidth, CanvasHeight, ColorBlack)
	sequential(sphere, canvas)
	p := PPMWriter{MaxLineLength: 70}
	p.Write(canvas)
	p.SaveFile("scene.ppm")
}

func sequential(sphere Sphere, canvas *Canvas) time.Duration {
	tStart := time.Now()
	for m := uint32(0); m < canvas.Height; m++ {
		for n := uint32(0); n < canvas.Width; n++ {
			work(m, n, sphere, canvas)
		}
	}
	return time.Since(tStart)
}

func work(m, n uint32, sphere Sphere, canvas *Canvas) {
	dir := computeDir(m, n, WallWidth, WallHeight, CanvasWidth, CanvasHeight)
	ray := Ray{Origin: ViewPos, Dir: dir}
	castRay(sphere, ray, canvas, m, n)
}

func castRay(s Sphere, ray Ray, canvas *Canvas, m, n uint32) {
	isects := s.Intersect(ray)
	if hit, isHit := Hit(isects); isHit {
		point := ray.Pos(hit.T)
		normal := s.NormalAt(point)
		eye := ray.Dir.Negate()
		color := mtl.Lighting(s.Material, Light, point, eye, normal)
		canvas.WritePixel(n, m, color)
	}
}

func computeDir(m, n uint32, wallW, wallH float64, canvasW, canvasH uint32) Vec4 {
	halfWallW := wallW/2.0
	halfWallH := wallH/2.0
	xw := wallW/float64(canvasW)*float64(n) - halfWallW
	yw := wallH/float64(canvasH)*float64(m) - halfWallH
	wallPoint := Point4(xw, yw, 0)
	return wallPoint.Sub(ViewPos).Normalize()
}
