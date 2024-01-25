package main

import (
	"fmt"
	"math"
	"os"
	"runtime/pprof"
	"time"

	"roytracer/camera"
	"roytracer/color"
	"roytracer/gfx"
	"roytracer/light"
	m "roytracer/math"
	"roytracer/mtl"
	"roytracer/shape"
	"roytracer/world"
)

const (
	Size      float64 = 25
	PosOffset         = -float64(Size) / 2.0
)

func main() {
	f, _ := os.Create("cpu.pprof")
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	mf, _ := os.Create("mem.pprof")
	defer pprof.WriteHeapProfile(mf)

	tstart := time.Now()
	w := GenerateWorld()
	camera := camera.NewCamera(1024, 512, math.Pi/3)
	camera.SetTf(m.View(
		m.Vec4{1.6 * Size, 1.2 * Size, -1.6 * Size},
		m.Vec4{0, -0.12 * Size, 0},
		m.Vec4{0, 1, 0},
	))
	image := camera.Render(w, 3)

	fmt.Println("Render took:", time.Since(tstart))

	writer := gfx.PPMWriter{MaxLineLength: 70}
	writer.Write(image)
	writer.SaveFile("scene.ppm")
	fmt.Println("Total:", time.Since(tstart))
}

func GenerateWorld() *world.World {
	w := world.World{
		Light: light.PointLight{
			Intensity: color.White,
			Pos:       m.Point4(Size*2, Size*2, -Size),
		},
	}

	mtl := mtl.DefaultMaterial()
	mtl.Reflective = 0.9
	mtl.Transparency = 1
	mtl.RefractiveIndex = 1.5

	for x := 0; x < int(Size); x++ {
		for y := 0; y < int(Size); y++ {
			for z := 0; z < int(Size); z++ {
				p := m.Vec4{float64(x), float64(y), float64(z)}
				mtl.Color = m.Vec4{
					p[0] / Size,
					p[1] / Size,
					p[2] / Size,
				}
				tf := m.Trans(
					p[0]+PosOffset,
					p[1]+PosOffset,
					p[2]+PosOffset,
				).Mul(m.Scale(0.33, 0.33, 0.33))
				s := shape.Sphere{
					O: shape.Object{
						Tf:       tf,
						InvTf:    tf.Inv(),
						Material: mtl,
					},
				}
				w.AddShape(&s)
			}
		}
	}
	return &w
}
