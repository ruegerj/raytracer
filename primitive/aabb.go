package primitive

import (
	"github.com/ruegerj/raytracing/common"
)

func MAX_AABB() AABB {
	return NewAABB(INFINITIY_VEC, NEG_INFINITY_VEC)
}

type AABB struct {
	Minimum Vec3
	Maximum Vec3
}

func NewAABB(minVec, maxVec Vec3) AABB {
	return AABB{
		Minimum: minVec,
		Maximum: maxVec,
	}
}

func (ab *AABB) Grow(vec Vec3) {
	ab.Minimum = ab.Minimum.Min(vec)
	ab.Maximum = ab.Maximum.Max(vec)
}

func (ab AABB) Area() float32 {
	extent := ab.Maximum.Sub(ab.Minimum)
	return extent.X*extent.Y + extent.Y*extent.Z + extent.Z*extent.X
}

func (ab AABB) Hit(ray Ray) float32 {
	t1 := (ab.Minimum.X - ray.Origin().X) * ray.DirectionInv().X
	t2 := (ab.Maximum.X - ray.Origin().X) * ray.DirectionInv().X

	tMin := min(t1, t2)
	tMax := max(t1, t2)

	t1 = (ab.Minimum.Y - ray.Origin().Y) * ray.DirectionInv().Y
	t2 = (ab.Maximum.Y - ray.Origin().Y) * ray.DirectionInv().Y

	tMin = max(tMin, min(t1, t2))
	tMax = min(tMax, max(t1, t2))

	t1 = (ab.Minimum.Z - ray.Origin().Z) * ray.DirectionInv().Z
	t2 = (ab.Maximum.Z - ray.Origin().Z) * ray.DirectionInv().Z

	tMin = max(tMin, min(t1, t2))
	tMax = min(tMax, max(t1, t2))

	if tMax >= max(tMin, 0.0) {
		return tMin
	}

	return common.F32_INF
}
