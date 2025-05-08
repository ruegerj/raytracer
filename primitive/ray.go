package primitive

type Ray struct {
	origin       Vec3
	direction    Vec3
	directionInv Vec3
}

func NewRay(origin, direction Vec3) Ray {
	return Ray{
		origin:       origin,
		direction:    direction,
		directionInv: direction.Inverse(),
	}
}

func (r Ray) Origin() Vec3 {
	return r.origin
}

func (r Ray) Direction() Vec3 {
	return r.direction
}

func (r Ray) DirectionInv() Vec3 {
	return r.directionInv
}

func (r Ray) Point(t float32) Vec3 {
	b := r.direction.MulScalar(t)
	return r.origin.Add(b)
}
