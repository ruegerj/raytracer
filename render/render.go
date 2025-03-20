package render

import (
	"image"
	"image/color"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/scene"
)

func Do(world *scene.World, img *image.RGBA, depth float64) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	cam := NewCamera(width, height, 1)

	for y := range height {
		for x := range width {
			r := cam.RayFrom(x, y)
			hit, hasHit := world.Hits(r)

			if !hasHit {
				img.Set(x, y, color.Black)
				continue
			}

			c := calcColor(hit, world.Lights()[0], false)
			img.Set(x, y, c.ToRGBA())
		}
	}
}

func calcColor(hit *scene.Hit, light scene.Light, ambient bool) primitive.ScalarColor {
	var ambientFactor float64 = 0
	if ambient {
		ambientFactor = 0.1
	}

	s := light.Origin.Sub(hit.Point).Normalize()
	intersectsLight := s.Dot(hit.Normal) >= 0

	if !intersectsLight {
		return hit.Color.MulScalar(ambientFactor)
	}

	ambientColor := hit.Color.AddScalar(ambientFactor)
	ambientLightColor := ambientColor.Mul(light.Color.MulScalar(float64(light.Intensity)))
	return ambientLightColor.MulScalar(s.Dot(hit.Normal))
}
