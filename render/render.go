package render

import (
	"image"
	"image/color"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/shape"
)

var light = primitive.Vector{X: 300, Y: 300, Z: 300}

func Do(target shape.Hitable, img *image.RGBA, depth float64) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	for y := range height {
		for x := range width {
			r := primitive.Ray{
				Origin:    primitive.Vector{X: float64(x), Y: float64(y), Z: 0},
				Direction: primitive.Vector{X: 0, Y: 0, Z: 1},
			}

			hit, hasHit := target.Hits(r)

			if !hasHit {
				img.Set(x, y, color.Black)
				continue
			}

			s := hit.Point.Sub(light).Normalize()
			shadedColor := hit.Color.MulScalar(hit.Normal.Mul(s).Length())
			shadedAmbientColor := shadedColor.AddScalar(0.1)

			img.Set(x, y, shadedAmbientColor.ToRGBA())
		}
	}
}
