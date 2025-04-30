package primitive

type Ray struct {
	origin       Vector
	direction    Vector
	directionInv Vector
}

func NewRay(origin, direction Vector) Ray {
	return Ray{
		origin:       origin,
		direction:    direction,
		directionInv: direction.Inverse(),
	}
}

func (r Ray) Origin() Vector {
	return r.origin
}

func (r Ray) Direction() Vector {
	return r.direction
}

func (r Ray) DirectionInv() Vector {
	return r.directionInv
}

func (r Ray) Point(t float32) Vector {
	b := r.direction.MulScalar(t)
	return r.origin.Add(b)
}
