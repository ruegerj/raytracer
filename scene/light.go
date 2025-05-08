package scene

import "github.com/ruegerj/raytracing/primitive"

type Light struct {
	Origin    primitive.Vec3
	Color     primitive.ScalarColor
	Intensity float32
}

func NewLight(origin primitive.Vec3, color primitive.ScalarColor, intensity float32) Light {
	return Light{
		Origin:    origin,
		Color:     color,
		Intensity: intensity,
	}
}
