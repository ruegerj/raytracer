package shape

import "github.com/ruegerj/raytracing/primitive"

type Circle struct {
	Origin primitive.Vector
	Radius float64
	Color  primitive.ScalarColor
}

func NewCircle(origin primitive.Vector, radius float64, color primitive.ScalarColor) Circle {
	return Circle{
		Origin: origin,
		Radius: radius,
		Color:  color,
	}
}

func (c Circle) Hits(p primitive.Vector) bool {
	return p.Distance(c.Origin) <= c.Radius
}
