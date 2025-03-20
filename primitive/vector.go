package primitive

import "math"

type Vector struct {
	X, Y, Z float64
}

var UnitVector = Vector{1, 1, 1}

func (v Vector) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

func (v Vector) Normalize() Vector {
	return v.DivScalar(v.Length())
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

func (v Vector) Distance(ov Vector) float64 {
	return v.Sub(ov).Length()
}

func (v Vector) Abs() float64 {
	return math.Sqrt(v.Dot(v))
}

// Computes the scalar product
func (v Vector) Dot(ov Vector) float64 {
	return v.X*ov.X + v.Y*ov.Y + v.Z*ov.Z
}

// Calculate the cross product of two vectors
func (v Vector) Cross(ov Vector) Vector {
	return Vector{
		X: v.Y*ov.Z - v.Z*ov.Y,
		Y: v.Z*ov.X - v.X*ov.Z,
		Z: v.X*ov.Y - v.Y*ov.X,
	}
}
