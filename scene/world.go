package scene

import (
	"github.com/ruegerj/raytracing/common/optional"
	"github.com/ruegerj/raytracing/primitive"
)

type ValidHitPredicate = func(*Hit, Hitable) bool

type World struct {
	lights []Light
	camera Camera
	bvh    Bvh
}

func NewWorld(triangles []Triangle, lights []Light, camera Camera) *World {
	return &World{
		lights: lights,
		camera: camera,
		bvh:    NewBvh(triangles),
	}
}

func (w *World) Camera() Camera {
	return w.camera
}

func (w *World) Lights() []Light {
	return w.lights
}

func (w *World) Hits(r primitive.Ray) optional.Optional[Hit] {
	return w.bvh.Intersects(r)
}
