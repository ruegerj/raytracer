package primitive

type Ray struct {
	Origin, Direction Vector
}

func (r Ray) Point(t float32) Vector {
	b := r.Direction.MulScalar(t)
	return r.Origin.Add(b)
}
