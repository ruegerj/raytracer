package primitive

import "github.com/ruegerj/raytracing/common"

type AABB struct {
	minimum Vector
	maximum Vector
}

func NewAABB(minVec, maxVec Vector) AABB {
	return AABB{
		minimum: minVec,
		maximum: maxVec,
	}
}

func (ab AABB) Hit(ray Ray) common.Optional[float32] {
	t1 := (ab.minimum.X - ray.Origin().X) * ray.DirectionInv().X
	t2 := (ab.maximum.X - ray.Origin().X) * ray.DirectionInv().X

	tMin := min(t1, t2)
	tMax := max(t1, t2)

	t1 = (ab.minimum.Y - ray.Origin().Y) * ray.DirectionInv().Y
	t2 = (ab.maximum.Y - ray.Origin().Y) * ray.DirectionInv().Y

	tMin = max(tMin, min(t1, t2))
	tMax = min(tMax, max(t1, t2))

	t1 = (ab.minimum.Z - ray.Origin().Z) * ray.DirectionInv().Z
	t2 = (ab.maximum.Z - ray.Origin().Z) * ray.DirectionInv().Z

	tMin = max(tMin, min(t1, t2))
	tMax = min(tMax, max(t1, t2))

	if tMax >= max(tMin, 0.0) {
		return common.Some(tMin)
	}

	return common.Empty[float32]()
}

func (ab AABB) Grow(vec Vector) {
	ab.minimum = ab.minimum.Min(vec)
	ab.maximum = ab.maximum.Max(vec)
}

func (ab AABB) Area() float32 {
	extent := ab.maximum.Sub(ab.minimum)
	return extent.X*extent.Y + extent.Y*extent.Z + extent.Z*extent.X
}
