package scene

import (
	"github.com/ruegerj/raytracing/primitive"
)

const epsilon = 0.0000001

var _ Hitable = (*Triangle)(nil)

type Triangle struct {
	V0, V1, V2 Vertex
	Normal     primitive.Vector
	Material   Material
}

type Vertex struct {
	Point  primitive.Vector
	Normal primitive.Vector
}

func NewTriangle(v0, v1, v2 Vertex, material Material) Triangle {
	triangle := Triangle{
		V0:       v0,
		V1:       v1,
		V2:       v2,
		Material: material,
	}
	triangle.Normal = v0.Normal

	return triangle
}

// Möller-Trumbore algorithm
func (tr Triangle) Hits(r primitive.Ray) (*Hit, bool) {
	edge1 := tr.V1.Point.Sub(tr.V0.Point)
	edge2 := tr.V2.Point.Sub(tr.V0.Point)

	h := r.Direction().Cross(edge2)
	a := edge1.Dot(h)

	if a > -epsilon && a < epsilon {
		return nil, false // Ray is parallel to the triangle
	}

	f := 1.0 / a
	s := r.Origin().Sub(tr.V0.Point)

	u := f * s.Dot(h)
	if u < 0.0 || u > 1.0 {
		return nil, false
	}

	q := s.Cross(edge1)
	v := f * r.Direction().Dot(q)
	if v < 0.0 || u+v > 1.0 {
		return nil, false
	}

	t := f * edge2.Dot(q)
	if t <= epsilon {
		return nil, false // Line intersection but not a ray intersection
	}

	intersection := r.Origin().Add(r.Direction().MulScalar(t))
	return &Hit{
		Distance: t,
		Point:    intersection,
		Normal:   tr.Normal,
		Material: tr.Material,
	}, true
}
