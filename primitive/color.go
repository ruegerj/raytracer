package primitive

import "image/color"

type ScalarColor struct {
	R float32
	G float32
	B float32
}

func (sc ScalarColor) ToRGBA() color.RGBA {
	return color.RGBA{
		R: uint8(sc.R) * 255,
		G: uint8(sc.G) * 255,
		B: uint8(sc.B) * 255,
	}
}

func (sc ScalarColor) Add(osc ScalarColor) ScalarColor {
	return ScalarColor{sc.R + osc.R, sc.G + osc.G, sc.B + osc.B}
}

func FromRGBAToScalar(base color.RGBA) ScalarColor {
	return ScalarColor{
		R: float32(base.R) / 255,
		G: float32(base.G) / 255,
		B: float32(base.B) / 255,
	}
}
