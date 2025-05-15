package primitive

import (
	"image/color"

	"github.com/ruegerj/raytracing/common"
)

var BLACK = ScalarColor{0, 0, 0}

type ScalarColor struct {
	R float32
	G float32
	B float32
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

func (sc ScalarColor) AddScalar(t float32) ScalarColor {
	return ScalarColor{
		R: clamp(sc.R + t),
		G: clamp(sc.G + t),
		B: clamp(sc.B + t),
	}
}

func (sc ScalarColor) Sub(osc ScalarColor) ScalarColor {
	return ScalarColor{
		R: clamp(sc.R - osc.R),
		G: clamp(sc.G - osc.G),
		B: clamp(sc.B - osc.B),
	}
}

func (sc ScalarColor) Mul(osc ScalarColor) ScalarColor {
	return ScalarColor{
		R: clamp(sc.R * osc.R),
		G: clamp(sc.G * osc.G),
		B: clamp(sc.B * osc.B),
	}
}

func (sc ScalarColor) MulScalar(t float32) ScalarColor {
	return ScalarColor{
		R: clamp(sc.R * t),
		G: clamp(sc.G * t),
		B: clamp(sc.B * t),
	}
}

func (sc ScalarColor) DivScalar(t float32) ScalarColor {
	return ScalarColor{
		R: clamp(sc.R / t),
		G: clamp(sc.G / t),
		B: clamp(sc.B / t),
	}
}

func (sc ScalarColor) GammaCorrect() ScalarColor {
	convert := func(c float32) float32 {
		if c <= 0.0031308 {
			return 12.92 * c
		}
		return 1.055*common.Pow(c, common.Recip(2.4)) - 0.055
	}

	return ScalarColor{
		R: convert(sc.R),
		G: convert(sc.G),
		B: convert(sc.B),
	}
}

func AvgColor(values []ScalarColor) ScalarColor {
	var sumR, sumG, sumB float32
	for _, v := range values {
		sumR += v.R
		sumG += v.G
		sumB += v.B
	}

	colorCount := float32(len(values))
	return ScalarColor{
		R: clamp(sumR / colorCount),
		G: clamp(sumG / colorCount),
		B: clamp(sumB / colorCount),
	}
}

func FromRGBAToScalar(base color.RGBA) ScalarColor {
	return ScalarColor{
		R: float32(base.R) / 255,
		G: float32(base.G) / 255,
		B: float32(base.B) / 255,
	}
}

func FromSlice(slice [3]float64) ScalarColor {
	return ScalarColor{
		R: clamp(float32(slice[0])),
		G: clamp(float32(slice[1])),
		B: clamp(float32(slice[2])),
	}
}

func clamp(v float32) float32 {
	if v < 0 {
		return 0
	}

	if v > 1 {
		return 1
	}

	return v
}
