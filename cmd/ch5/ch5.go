package main

import (
	"fmt"
	"sync"
	"time"

	. "roytracer/gfx"
	. "roytracer/math"
	. "roytracer/shape"
)

var (
	LightPos            = Vec4{0, 0, -3}
	SpherePos           = Vec3{0, 0, -1}
	SphereRadius        = 2.0
	WallWidth           = 4.0
	WallHeight          = 4.0
	CanvasWidth  uint32 = 300
	CanvasHeight uint32 = 300
)

type SyncCanvas struct {
	C *Canvas
	sync.Mutex
	wg sync.WaitGroup
}

func main() {
	sphere := NewSphere()
	sphere.SetTf(Trans(Vec3{0, 0, 1}).Mul(Scale(Vec3{SphereRadius, SphereRadius, SphereRadius})))
	canvas := SyncCanvas{C: NewCanvas(CanvasWidth, CanvasHeight, ColorBlack)}
	fmt.Printf("seq: %v\n", sequential(sphere, &canvas))
	fmt.Printf("par: %v\n", parallel(sphere, &canvas))
	p := PPMWriter{MaxLineLength: 70}
	p.Write(canvas.C)
	p.SaveFile("scene.ppm")
}

func sequential(sphere Sphere, canvas *SyncCanvas) time.Duration {
	tStart := time.Now()
	for m := uint32(0); m < canvas.C.Height; m++ {
		for n := uint32(0); n < canvas.C.Width; n++ {
			work(m, n, sphere, canvas)
		}
	}
	return time.Since(tStart)
}

func parallel(sphere Sphere, canvas *SyncCanvas) time.Duration {
	tStart := time.Now()
	for m := uint32(0); m < canvas.C.Height; m++ {
		canvas.wg.Add(1)
		go func(mm uint32, canvas *SyncCanvas) {
			for n := uint32(0); n < canvas.C.Width; n++ {
				work(mm, n, sphere, canvas)
			}
			canvas.wg.Done()
		}(m, canvas)
	}
	canvas.wg.Wait()
	return time.Since(tStart)
}

func work(m, n uint32, sphere Sphere, canvas *SyncCanvas) {
	dir := computeDir(m, n, WallWidth, WallHeight, CanvasWidth, CanvasHeight)
	ray := Ray{Origin: LightPos, Dir: dir}
	castRay(sphere, ray, canvas, m, n)
}

func castRay(s Sphere, ray Ray, canvas *SyncCanvas, m, n uint32) {
	isects := s.Intersect(ray)
	if _, isHit := Hit(isects); isHit {
		canvas.C.WritePixel(n, m, ColorRed)
	}
}

func computeDir(m, n uint32, wallW, wallH float64, canvasW, canvasH uint32) Vec4 {
	halfWallW := wallW / 2.0
	halfWallH := wallH / 2.0
	xw := wallW/float64(canvasW)*float64(n) - halfWallW
	yw := wallH/float64(canvasH)*float64(m) - halfWallH
	wallPoint := Point4(xw, yw, 0)
	return wallPoint.Sub(LightPos).Normalize()
}
