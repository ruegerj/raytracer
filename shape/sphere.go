package shape

import (
	"math"

	"github.com/ruegerj/raytracing/primitive"
)

type Sphere struct {
	Center primitive.Vector
	Radius float64
	color  primitive.ScalarColor
}

func NewSphere(origin primitive.Vector, radius float64, color primitive.ScalarColor) Sphere {
	return Sphere{
		Center: origin,
		Radius: radius,
		color:  color,
	}
}

func (c Sphere) Color() primitive.ScalarColor {
	return c.color
}

func (c Sphere) HitsVector(p primitive.Vector) bool {
	return p.Distance(c.Center) <= c.Radius
}

func (c Sphere) Hits(r primitive.Ray) (*Hit, bool) {
	u := r.Direction
	v := c.Center.Sub(r.Origin)

	ma := u.Length() * u.Length()
	mb := 2 * u.Dot(v)
	mc := v.Length()*v.Length() - c.Radius*c.Radius

	discriminant := mb*mb - 4*ma*mc

	if discriminant < 0 {
		return nil, false
	}

	lambda := (-mb - math.Sqrt(discriminant)) / 2 * ma

	return &Hit{Distance: lambda, Element: c}, true
}
