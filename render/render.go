package render

import (
	"image"
	"image/color"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/shape"
)

var light = primitive.Vector{X: -10, Y: 7, Z: 18}

func Do(target shape.Hitable, img *image.RGBA, depth float64) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	cam := NewCamera(width, height, 1)

	for y := range height {
		for x := range width {
			r := cam.RayFrom(x, y)
			hit, hasHit := target.Hits(r)

			if !hasHit {
				img.Set(x, y, color.Black)
				continue
			}

			s := light.Sub(hit.Point).Normalize()

			if s.Dot(hit.Normal) < 0 {
				img.Set(x, y, hit.Color.MulScalar(0.1).ToRGBA())
				continue
			}

			shadedColor := hit.Color.MulScalar(s.Dot(hit.Normal))
			shadedAmbientColor := shadedColor.AddScalar(0.1)
			img.Set(x, y, shadedAmbientColor.ToRGBA())
		}
	}
}
