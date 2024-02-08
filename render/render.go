package render

import (
	"roytracer/camera"
	"roytracer/gfx"
	"roytracer/world"
	"sync"
)

type Renderer struct {
	Camera         *camera.Camera
	World          *world.World
	Canvas         *gfx.Canvas
	NProc          int
	RecursiveDepth int
}

func NewRenderer(camera *camera.Camera, world *world.World, nproc, recursiveDepth int) Renderer {
	return Renderer{
		Camera:         camera,
		World:          world,
		Canvas:         gfx.NewCanvas(uint32(camera.Hsize), uint32(camera.Vsize), gfx.ColorBlack),
		NProc:          nproc,
		RecursiveDepth: recursiveDepth,
	}
}

func (r *Renderer) RenderSequential() {
	for m := 0; m < r.Camera.Vsize; m++ {
		for n := 0; n < r.Camera.Hsize; n++ {
			r.computePixel(m, n)
		}
	}
}

func (r *Renderer) RenderParallel() {
	var wg sync.WaitGroup
	wg.Add(r.Camera.Vsize)
	for m := 0; m < r.Camera.Vsize; m++ {
		go func(m int) {
			for n := 0; n < r.Camera.Hsize; n++ {
				r.computePixel(m, n)
			}
			wg.Done()
		}(m)
	}
	wg.Wait()
}

func (r *Renderer) RenderHorizontalChunks(nChunks int) {
	var wg sync.WaitGroup
	wg.Add(nChunks)
	chunkSize := r.Camera.Vsize / nChunks
	for i := 0; i < nChunks; i++ {
		startRow := i*chunkSize
		endRow := startRow+chunkSize
		go func(startRow, endRow int) {
			for m := startRow; m < endRow; m++ {
				for n := 0; n < r.Camera.Hsize; n++ {
					r.computePixel(m, n)
				}
			}
			wg.Done()
		}(startRow, endRow)
	}
	wg.Wait()
}

func (r *Renderer) computePixel(m, n int) {
	ray := r.Camera.RayForPixel(n, m)
	color := r.World.ColorAt(ray, r.RecursiveDepth)
	r.Canvas.WritePixel(uint32(n), uint32(m), color)
}
