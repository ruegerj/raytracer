package primitive

import "math"

type Vector struct {
	X, Y, Z float64
}

var UnitVector = Vector{1, 1, 1}

func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector) Add(ov Vector) Vector {
	return Vector{v.X + ov.X, v.Y + ov.Y, v.Z + ov.Z}
}

func (v Vector) Sub(ov Vector) Vector {
	return Vector{v.X - ov.X, v.Y - ov.Y, v.Z - ov.Z}
}

func (v Vector) Mul(ov Vector) Vector {
	return Vector{v.X * ov.X, v.Y * ov.Y, v.Z * ov.Z}
}

func (v Vector) Div(ov Vector) Vector {
	return Vector{v.X / ov.X, v.Y / ov.Y, v.Z / ov.Z}
}

func (v Vector) AddScalar(scalar float64) Vector {
	return Vector{v.X + scalar, v.Y + scalar, v.Z + scalar}
}

func (v Vector) SubScalar(scalar float64) Vector {
	return Vector{v.X - scalar, v.Y - scalar, v.Z - scalar}
}

func (v Vector) MulScalar(scalar float64) Vector {
	return Vector{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vector) DivScalar(scalar float64) Vector {
	return Vector{v.X / scalar, v.Y / scalar, v.Z / scalar}
}
