package scene

import (
	"log"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/schollz/progressbar/v3"
)

type World struct {
	lights []Light
	camera Camera
	bvh    *Bvh
}

func NewWorld(triangles []Triangle, lights []Light, camera Camera) *World {
	spinner := progressbar.Default(-1, "building bvh tree")
	bvh := NewBvh(triangles)
	_ = spinner.Close()
	log.Printf("bvh node count: %d\n", len(bvh.nodes))

	return &World{
		lights: lights,
		camera: camera,
		bvh:    bvh,
	}
}

func (w *World) Camera() Camera {
	return w.camera
}

func (w *World) Lights() []Light {
	return w.lights
}

func (w *World) Hits(r primitive.Ray) *Hit {
	return w.bvh.Intersects(r)
}
