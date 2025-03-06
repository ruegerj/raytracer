package render

import (
	"image"
	"image/color"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/shape"
)

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

			depthScalar := 1 - (hit.Distance*-1)/depth
			effectiveColor := hit.Element.Color().MulScalar(float32(depthScalar))

			img.Set(x, y, effectiveColor.ToRGBA())
		}
	}
}
