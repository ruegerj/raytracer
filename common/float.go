package common

import "math"

var F32_INF = float32(math.Inf(1))

func Recip(value float32) float32 {
	return 1 / value
}

func Pow(base, exp float32) float32 {
	return float32(math.Pow(float64(base), float64(exp)))
}
