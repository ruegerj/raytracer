package render

import (
	"fmt"
	"image"
	"image/color"
	"math"

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

			c := calcColor(hit, world, false)
			img.Set(x, y, c.ToRGBA())
		}
	}
}

// epsilon: 10^-9
// start ray: hitpoint + vector to light * epsilon
func calcColor(hit *scene.Hit, world *scene.World, ambient bool) primitive.ScalarColor {
	var ambientFactor float64 = 0
	if ambient {
		ambientFactor = 0.1
	}

	lightFactors := []float64{}

	for _, l := range world.Lights() {
		lightVec := l.Origin.Sub(hit.Point)
		s := lightVec.Normalize()
		lightFactor := s.Dot(hit.Normal)

		if lightFactor < 0 {
			return hit.Color.MulScalar(ambientFactor)
		}

		lightRay := primitive.Ray{
			Origin: hit.Point.Add(lightVec.MulScalar(math.Pow(10, -9))),
			// Origin:    hit.Point,
			Direction: lightVec.Normalize(),
		}

		blockHit, hasBlockHit := world.Hits(lightRay)
		if hasBlockHit && blockHit.Distance < lightVec.Length() {
			lightFactor = 0
		}

		lightFactors = append(lightFactors, lightFactor)
	}

	avg := avgValue(lightFactors)
	fmt.Println(avg, lightFactors)
	shadedColor := hit.Color.MulScalar(avg + ambientFactor)
	return shadedColor
}

func avgValue(values []float64) float64 {
	var sum float64 = 0
	for _, v := range values {
		sum += v
	}

	return sum / float64(len(values))
}
