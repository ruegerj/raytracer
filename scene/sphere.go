package scene

import (
	"math"

	"github.com/ruegerj/raytracing/primitive"
)

var _ Hitable = (*Sphere)(nil)

type Sphere struct {
	Center   primitive.Vec3
	Radius   float32
	Material Material
}

func NewSphere(origin primitive.Vec3, radius float32, color primitive.ScalarColor) Sphere {
	return Sphere{
		Center:   origin,
		Radius:   radius,
		Material: NewPhong(color, 1),
	}
}

func (s Sphere) HitsVector(p primitive.Vec3) bool {
	return p.Distance(s.Center) <= s.Radius
}

func (s Sphere) Hits(r primitive.Ray) (*Hit, bool) {
	u := r.Direction()
	v := r.Origin().Sub(s.Center)

	ma := float64(u.Dot(u))
	mb := float64(2 * u.Dot(v))
	mc := float64(v.Dot(v) - s.Radius*s.Radius)

	discriminant := mb*mb - 4*ma*mc

	if discriminant < 0 {
		return nil, false
	}

	lambda := (-mb - math.Sqrt(discriminant)) / 2 * ma

	q := r.Point(float32(lambda))
	n := q.Sub(s.Center).Normalize()

	hit := &Hit{
		Distance: float32(lambda),
		Point:    q,
		Normal:   n,
		Material: s.Material,
	}

	return hit, true
}
