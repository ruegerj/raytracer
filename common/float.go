package common

import "math"

func Min(value, min float32) float32 {
	if value < min {
		return min
	}
	return value
}

func Recip(value float32) float32 {
	return 1 / value
}

func Pow(base, exp float32) float32 {
	return float32(math.Pow(float64(base), float64(exp)))
}
