package shape

import (
	"math"

	"github.com/ruegerj/raytracing/primitive"
)

type Sphere struct {
	Center primitive.Vector
	Radius float64
	Color  primitive.ScalarColor
}

func NewSphere(origin primitive.Vector, radius float64, color primitive.ScalarColor) Sphere {
	return Sphere{
		Center: origin,
		Radius: radius,
		Color:  color,
	}
}

func (s Sphere) HitsVector(p primitive.Vector) bool {
	return p.Distance(s.Center) <= s.Radius
}

func (s Sphere) Hits(r primitive.Ray) (*Hit, bool) {
	u := r.Direction
	v := s.Center.Sub(r.Origin)

	ma := u.Dot(u)
	mb := 2 * u.Dot(v)
	mc := v.Dot(v) - s.Radius*s.Radius

	discriminant := mb*mb - 4*ma*mc

	if discriminant < 0 {
		return nil, false
	}

	lambda := (-mb - math.Sqrt(discriminant)) / 2 * ma

	q := r.Point(lambda)
	n := q.Sub(s.Center).Normalize()

	hit := &Hit{
		Distance: lambda,
		Point:    q,
		Normal:   n,
		Color:    s.Color,
	}

	return hit, true
}
