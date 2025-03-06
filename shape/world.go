package shape

import (
	"math"

	"github.com/ruegerj/raytracing/primitive"
)

type World []Hitable

func (w *World) Add(elem Hitable) {
	*w = append(*w, elem)
}

func (w *World) AddAll(elems ...Hitable) {
	for _, elem := range elems {
		w.Add(elem)
	}
}

func (w *World) Color() primitive.ScalarColor {
	return primitive.ScalarColor{R: 0, G: 0, B: 0}
}

func (w *World) Hits(r primitive.Ray) (*Hit, bool) {
	var closestHit *Hit = nil
	closestDist := math.MaxFloat64

	for _, elem := range *w {
		hit, hits := elem.Hits(r)
		if !hits {
			continue
		}

		if hit.Distance < closestDist {
			closestDist = hit.Distance
			closestHit = hit
		}
	}

	return closestHit, closestHit != nil
}
