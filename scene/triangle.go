package scene

import (
	"github.com/ruegerj/raytracing/primitive"
)

const epsilon = 0.0000001
const centroid_factor = 0.3333333

type Triangle struct {
	V0, V1, V2 Vertex
	Centroid   primitive.Vec3
	Material   Material
}

type Vertex struct {
	Point  primitive.Vec3
	Normal primitive.Vec3
	UV     *primitive.Vec2
}

func NewTriangle(v0, v1, v2 Vertex, material Material) Triangle {
	centroid := v0.Point.Add(v1.Point).Add(v2.Point).MulScalar(centroid_factor)

	triangle := Triangle{
		V0:       v0,
		V1:       v1,
		V2:       v2,
		Centroid: centroid,
		Material: material,
	}

	return triangle
}

// Möller-Trumbore algorithm
func (tr Triangle) Hits(r primitive.Ray) *Hit {
	edge1 := tr.V1.Point.Sub(tr.V0.Point)
	edge2 := tr.V2.Point.Sub(tr.V0.Point)

	h := r.Direction().Cross(edge2)
	a := edge1.Dot(h)

	if a > -epsilon && a < epsilon {
		return nil // Ray is parallel to the triangle
	}

	f := 1.0 / a
	s := r.Origin().Sub(tr.V0.Point)

	u := f * s.Dot(h)
	if u < 0.0 || u > 1.0 {
		return nil
	}

	q := s.Cross(edge1)
	v := f * r.Direction().Dot(q)
	if v < 0.0 || u+v > 1.0 {
		return nil
	}

	t := f * edge2.Dot(q)
	if t <= epsilon {
		return nil // Line intersection but not a ray intersection
	}

	intersection := r.Origin().Add(r.Direction().MulScalar(t))
	return &Hit{
		Distance: t,
		Point:    intersection,
		Normal:   tr.V0.Normal,
		Material: tr.Material,
	}
}

func (tr Triangle) CreateHitFor(ray primitive.Ray, dist float32) *Hit {
	pointVec := ray.Point(dist)
	barycentric := tr.barycentricCoordinates(pointVec)

	var uv *primitive.Vec2
	if tr.V0.UV != nil && tr.V1.UV != nil && tr.V2.UV != nil {
		uv0 := tr.V0.UV
		uv1 := tr.V1.UV
		uv2 := tr.V2.UV

		result := uv0.MulScalar(barycentric.X).
			Add(uv1.MulScalar(barycentric.Y)).
			Add(uv2.MulScalar(barycentric.Z))

		uv = &result
	}

	normal := tr.V0.Normal.MulScalar(barycentric.X).
		Add(tr.V1.Normal.MulScalar(barycentric.Y)).
		Add(tr.V2.Normal.MulScalar(barycentric.Z))

	hitsFront := true
	if ray.Direction().Dot(normal) > 0.0 {
		normal = normal.MulScalar(-1)
		hitsFront = false
	}

	return &Hit{
		Distance:  dist,
		Point:     pointVec,
		Normal:    normal,
		UV:        uv,
		FrontFace: hitsFront,
		Material:  tr.Material,
	}
}

func (tr Triangle) barycentricCoordinates(p primitive.Vec3) primitive.Vec3 {
	v0v1 := tr.V1.Point.Sub(tr.V0.Point)
	v0v2 := tr.V2.Point.Sub(tr.V0.Point)
	v0p := p.Sub(tr.V0.Point)

	d11 := v0v1.Dot(v0v1)
	d12 := v0v1.Dot(v0v2)
	d22 := v0v2.Dot(v0v2)
	d31 := v0p.Dot(v0v1)
	d32 := v0p.Dot(v0v2)

	denom := d11*d22 - d12*d12
	v := (d22*d31 - d12*d32) / denom
	w := (d11*d32 - d12*d31) / denom
	u := 1.0 - v - w

	return primitive.Vec3{X: u, Y: v, Z: w}
}
