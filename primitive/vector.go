package primitive

import "math"

type Vec3 struct {
	X, Y, Z float32
}

var UnitVector = Vec3{1, 1, 1}
var INFINTIY_VEC = Vec3{float32(math.Inf(1)), float32(math.Inf(1)), float32(math.Inf(1))}
var NEG_INFINITY_VEC = Vec3{float32(math.Inf(-1)), float32(math.Inf(-1)), float32(math.Inf(-1))}

func (v Vec3) Axis(index uint) float32 {
	if index > 2 {
		return math.MaxFloat32
	}

	if index == 0 {
		return v.X
	}

	if index == 1 {
		return v.Y
	}

	return v.Z
}

func (v Vec3) Length() float32 {
	value := v.X*v.X + v.Y*v.Y + v.Z*v.Z
	return float32(math.Abs(math.Sqrt(float64(value))))
}

func (v Vec3) Normalize() Vec3 {
	return v.DivScalar(float32(v.Length()))
}

func (v Vec3) Inverse() Vec3 {
	return Vec3{
		X: 1.0 / v.X,
		Y: 1.0 / v.Y,
		Z: 1.0 / v.Z,
	}
}

func (v Vec3) Min(ov Vec3) Vec3 {
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

func (v Vec3) Max(ov Vec3) Vec3 {
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

func (v Vec3) Add(ov Vec3) Vec3 {
	return Vec3{v.X + ov.X, v.Y + ov.Y, v.Z + ov.Z}
}

func (v Vec3) Sub(ov Vec3) Vec3 {
	return Vec3{v.X - ov.X, v.Y - ov.Y, v.Z - ov.Z}
}

func (v Vec3) Mul(ov Vec3) Vec3 {
	return Vec3{v.X * ov.X, v.Y * ov.Y, v.Z * ov.Z}
}

func (v Vec3) Div(ov Vec3) Vec3 {
	return Vec3{v.X / ov.X, v.Y / ov.Y, v.Z / ov.Z}
}

func (v Vec3) AddScalar(scalar float32) Vec3 {
	return Vec3{v.X + scalar, v.Y + scalar, v.Z + scalar}
}

func (v Vec3) SubScalar(scalar float32) Vec3 {
	return Vec3{v.X - scalar, v.Y - scalar, v.Z - scalar}
}

func (v Vec3) MulScalar(scalar float32) Vec3 {
	return Vec3{v.X * scalar, v.Y * scalar, v.Z * scalar}
}

func (v Vec3) DivScalar(scalar float32) Vec3 {
	return Vec3{v.X / scalar, v.Y / scalar, v.Z / scalar}
}

func (v Vec3) Distance(ov Vec3) float32 {
	return v.Sub(ov).Length()
}

func (v Vec3) Abs() float32 {
	return float32(math.Sqrt(float64(v.Dot(v))))
}

// Computes the scalar product
func (v Vec3) Dot(ov Vec3) float32 {
	return v.X*ov.X + v.Y*ov.Y + v.Z*ov.Z
}

// Calculate the cross product of two vectors
func (v Vec3) Cross(ov Vec3) Vec3 {
	return Vec3{
		X: v.Y*ov.Z - v.Z*ov.Y,
		Y: v.Z*ov.X - v.X*ov.Z,
		Z: v.X*ov.Y - v.Y*ov.X,
	}
}

func (v Vec3) Reflect(normal Vec3) Vec3 {
	return v.Sub(v.Mul(normal).MulScalar(2).Mul(normal))
}

type Vec2 struct {
	X, Y float32
}

func (v Vec2) Add(ov Vec2) Vec2 {
	return Vec2{
		X: v.X + ov.X,
		Y: v.Y + ov.Y,
	}
}

func (v Vec2) MulScalar(t float32) Vec2 {
	return Vec2{
		X: v.X * t,
		Y: v.Y * t,
	}
}
