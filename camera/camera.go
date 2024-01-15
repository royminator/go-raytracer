package camera

import (
	"math"
	"sync"

	"roytracer/gfx"
	m "roytracer/math"
	"roytracer/shape"
	"roytracer/world"
)

type Camera struct {
	Hsize, Vsize int
	Fov          float64
	PixelSize    float64
	HalfW, HalfH float64
	Tf           m.Mat4
	InvTf        m.Mat4
}

func NewCamera(hsize, vsize int, fov float64) Camera {
	halfView := math.Tan(fov / 2.0)
	aspect := float64(hsize) / float64(vsize)

	var halfW, halfH float64

	if aspect >= 1.0 {
		halfW = halfView
		halfH = halfView / aspect
	} else {
		halfW = halfView * aspect
		halfH = halfView
	}
	pixelSize := halfW * 2.0 / float64(hsize)

	return Camera{
		Hsize:     hsize,
		Vsize:     vsize,
		Fov:       fov,
		PixelSize: pixelSize,
		HalfW:     halfW,
		HalfH:     halfH,
		Tf:        m.Mat4Ident(),
		InvTf:     m.Mat4Ident(),
	}
}

func (c *Camera) SetTf(tf m.Mat4) {
	c.Tf = tf
	c.InvTf = tf.Inv()
}

func (c *Camera) RayForPixel(px, py int) shape.Ray {
	xoffset := (float64(px) + 0.5) * c.PixelSize
	yoffset := (float64(py) + 0.5) * c.PixelSize
	worldx := c.HalfW - xoffset
	worldy := c.HalfH - yoffset
	pixel := c.InvTf.MulVec(m.Point4(worldx, worldy, -1))
	origin := c.InvTf.MulVec(m.Point4(0, 0, 0))
	dir := pixel.Sub(origin).Normalize()
	return shape.Ray{Origin: origin, Dir: dir}
}

func (c *Camera) Render(w *world.World) *gfx.Canvas {
	canvas := gfx.NewCanvas(uint32(c.Hsize), uint32(c.Vsize), gfx.ColorBlack)
	// c.renderSequential(w, canvas)
	c.renderParallel(w, canvas)
	return canvas
}

func (c *Camera) renderSequential(w *world.World, canvas *gfx.Canvas) {
	for m := 0; m < c.Vsize; m++ {
		for n := 0; n < c.Hsize; n++ {
			c.computePixel(m, n, w, canvas)
		}
	}
}

func (c *Camera) renderParallel(w *world.World, canvas *gfx.Canvas) {
	var wg sync.WaitGroup
	for m := 0; m < c.Vsize; m++ {
		wg.Add(1)
		go func(m int, w *world.World, canvas *gfx.Canvas) {
			for n := 0; n < c.Hsize; n++ {
				c.computePixel(m, n, w, canvas)
			}
			wg.Done()
		}(m, w, canvas)
	}
	wg.Wait()
}

func (c *Camera) computePixel(m, n int, w *world.World, canvas *gfx.Canvas) {
	ray := c.RayForPixel(n, m)
	color := w.ColorAt(ray, 10)
	canvas.WritePixel(uint32(n), uint32(m), color)
}
