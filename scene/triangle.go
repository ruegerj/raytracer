package scene

import (
	"github.com/ruegerj/raytracing/common"
	"github.com/ruegerj/raytracing/primitive"
)

const epsilon = 0.0000001
const centroid_factor = 0.3333333

var _ Hitable = (*Triangle)(nil)

type Triangle struct {
	V0, V1, V2 Vertex
	Centroid   primitive.Vec3
	Material   Material
}

type Vertex struct {
	Point  primitive.Vec3
	Normal primitive.Vec3
	UV     common.Optional[primitive.Vec2]
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
		Normal:   tr.V0.Normal,
		Material: tr.Material,
	}, true
}

func (tr Triangle) CreateHitFor(ray primitive.Ray, dist float32) Hit {
	pointVec := ray.Point(dist)
	barycentric := tr.barycentricCoordinats(pointVec)

	uv := common.Empty[primitive.Vec2]()
	if tr.V0.UV.IsPresent() && tr.V1.UV.IsPresent() && tr.V2.UV.IsPresent() {
		uv0 := tr.V0.UV.Get()
		uv1 := tr.V1.UV.Get()
		uv2 := tr.V2.UV.Get()

		result := uv0.MulScalar(barycentric.X).
			Add(uv1.MulScalar(barycentric.Y)).
			Add(uv2.MulScalar(barycentric.Z))

		uv = common.Some(result)
	}

	normal := tr.V0.Normal.MulScalar(barycentric.X).
		Add(tr.V1.Normal.MulScalar(barycentric.Y)).
		Add(tr.V2.Normal.MulScalar(barycentric.Z))

	hitsFront := true
	if ray.Direction().Dot(normal) > 0.0 {
		normal = normal.MulScalar(-1)
		hitsFront = false
	}

	return Hit{
		Distance:  dist,
		Point:     pointVec,
		Normal:    normal,
		UV:        uv,
		FrontFace: hitsFront,
		Material:  tr.Material,
	}
}

func (tr Triangle) barycentricCoordinats(p primitive.Vec3) primitive.Vec3 {
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

	return primitive.Vec3{X: v, Y: w, Z: u}
}
