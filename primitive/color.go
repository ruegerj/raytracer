package primitive

import (
	"image/color"
)

type ScalarColor struct {
	R float64
	G float64
	B float64
}

func (sc ScalarColor) ToRGBA() color.RGBA {
	return color.RGBA{
		R: uint8(sc.R * 255),
		G: uint8(sc.G * 255),
		B: uint8(sc.B * 255),
	}
}

func (sc ScalarColor) Add(osc ScalarColor) ScalarColor {
	return ScalarColor{
		R: clamp(sc.R + osc.R),
		G: clamp(sc.G + osc.G),
		B: clamp(sc.B + osc.B),
	}
}

func (sc ScalarColor) AddScalar(t float64) ScalarColor {
	return ScalarColor{
		R: clamp(sc.R + t),
		G: clamp(sc.G + t),
		B: clamp(sc.B + t),
	}
}

func (sc ScalarColor) Mul(osc ScalarColor) ScalarColor {
	return ScalarColor{
		R: clamp(sc.R * osc.R),
		G: clamp(sc.G * osc.G),
		B: clamp(sc.B * osc.B),
	}
}

func (sc ScalarColor) MulScalar(t float64) ScalarColor {
	return ScalarColor{
		R: clamp(sc.R * t),
		G: clamp(sc.G * t),
		B: clamp(sc.B * t),
	}
}

func FromRGBAToScalar(base color.RGBA) ScalarColor {
	return ScalarColor{
		R: float64(base.R) / 255,
		G: float64(base.G) / 255,
		B: float64(base.B) / 255,
	}
}

func clamp(v float64) float64 {
	if v < 0 {
		return 0
	}

	if v > 1 {
		return 1
	}

	return v
}
