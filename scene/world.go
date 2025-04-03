package scene

import (
	"math"

	"github.com/ruegerj/raytracing/primitive"
)

type ValidHitPredicate = func(*Hit, Hitable) bool

type World struct {
	elements []Hitable
	lights   []Light
	camera   Camera
}

func NewWorld(elems []Hitable, lights []Light, camera Camera) *World {
	return &World{
		elements: elems,
		lights:   lights,
		camera:   camera,
	}
}

func (w *World) Camera() Camera {
	return w.camera
}

func (w *World) Lights() []Light {
	return w.lights
}

func (w *World) Add(elem Hitable) {
	w.elements = append(w.elements, elem)
}

func (w *World) AddAll(elems ...Hitable) {
	for _, elem := range elems {
		w.Add(elem)
	}
}

func (w *World) AddLight(light Light) {
	w.lights = append(w.lights, light)
}

func (w *World) Color() primitive.ScalarColor {
	return primitive.ScalarColor{R: 0, G: 0, B: 0}
}

func (w *World) Hits(r primitive.Ray, isValidHit ValidHitPredicate) (*Hit, bool) {
	var closestHit *Hit = nil
	closestDist := float32(math.MaxFloat32)

	for _, elem := range w.elements {
		hit, hits := elem.Hits(r)
		if !hits || !isValidHit(hit, elem) {
			continue
		}

		if hit.Distance < closestDist {
			closestDist = hit.Distance
			closestHit = hit
		}
	}

	return closestHit, closestHit != nil
}
