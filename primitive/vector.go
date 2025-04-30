package primitive

import "math"

type Vector struct {
	X, Y, Z float32
}

var UnitVector = Vector{1, 1, 1}

func (v Vector) Length() float32 {
	value := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	return float32(math.Abs(math.Sqrt(float64(value))))
}

func (v Vector) Normalize() Vector {
	return v.DivScalar(float32(v.Length()))
}

func (v Vector) Inverse() Vector {
	return Vector{
		X: 1.0 / v.X,
		Y: 1.0 / v.Y,
		Z: 1.0 / v.Z,
	}
}

func (v Vector) Min(ov Vector) Vector {
	if v.X < ov.X {
		return v
	} else if ov.X < v.X {
		return ov
	}

	if v.Y < ov.Y {
		return v
	} else if ov.Y < v.Y {
		return ov
	}

	if v.Z < ov.Z {
		return v
	}

	return ov
}

func (v Vector) Max(ov Vector) Vector {
	if v.X > ov.X {
		return v
	} else if ov.X > v.X {
		return ov
	}

	if v.Y > ov.Y {
		return v
	} else if ov.Y > v.Y {
		return ov
	}

	if v.Z > ov.Z {
		return v
	}

	return ov
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

func (v Vector) AddScalar(scalar float32) Vector {
	return Vector{v.X + scalar, v.Y + scalar, v.Z + scalar}
}

func (v Vector) SubScalar(scalar float32) Vector {
	return Vector{v.X - scalar, v.Y - scalar, v.Z - scalar}
}

func (v Vector) MulScalar(scalar float32) Vector {
	return Vector{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vector) DivScalar(scalar float32) Vector {
	return Vector{v.X / scalar, v.Y / scalar, v.Z / scalar}
}

func (v Vector) Distance(ov Vector) float32 {
	return v.Sub(ov).Length()
}

func (v Vector) Abs() float32 {
	return float32(math.Sqrt(float64(v.Dot(v))))
}

// Computes the scalar product
func (v Vector) Dot(ov Vector) float32 {
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

func (v Vector) Reflect(normal Vector) Vector {
	return v.Sub(v.Mul(normal).MulScalar(2).Mul(normal))
}
