package scene

import "github.com/ruegerj/raytracing/primitive"

type Light struct {
	Origin    primitive.Vector
	Color     primitive.ScalarColor
	Intensity float32
}

func NewLight(origin primitive.Vector, color primitive.ScalarColor, intensity float32) Light {
	return Light{
		Origin:    origin,
		Color:     color,
		Intensity: intensity,
	}
}
