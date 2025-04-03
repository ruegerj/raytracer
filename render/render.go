package render

import (
	"image"
	"log"
	"sync"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/scene"
)

const epsilon = 1e-9

func Do(world *scene.World, img *image.RGBA) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	acceptAnyHit := func(_ *scene.Hit, _ scene.Hitable) bool { return true }

	imageBuffer := createImageBuffer(width, height)

	var wg sync.WaitGroup
	wg.Add(height)

	for y := range height {
		go func() {
			defer wg.Done()

			for x := range width {
				r := world.Camera().RayFrom(x, y)
				hit, hasHit := world.Hits(r, acceptAnyHit)

				if !hasHit {
					imageBuffer[y][x] = primitive.ScalarColor{R: 0, G: 0, B: 0}
					continue
				}

				c := calcColor(hit, world, true)
				imageBuffer[y][x] = c
			}
		}()
	}

	wg.Wait()

	exportBufferToImage(imageBuffer, img)

	log.Println("done")
}

func calcColor(hit *scene.Hit, world *scene.World, ambient bool) primitive.ScalarColor {
	var ambientFactor float32 = 0
	if ambient {
		ambientFactor = 0.1
	}

	lightFactors := []float32{}

	for _, light := range world.Lights() {
		var lightFactor float32 = 0
		lightVec := light.Origin.Sub(hit.Point)

		lightRay := primitive.Ray{
			Origin:    hit.Point.Add(hit.Normal.MulScalar(epsilon)),
			Direction: lightVec.Normalize(),
		}

		isValidShadowHit := func(elemHit *scene.Hit, elem scene.Hitable) bool {
			isNoSelfIntersection := elemHit.Distance > epsilon && lightRay.Direction.Dot(elemHit.Normal) <= 0
			isNotBehindLight := float64(elemHit.Distance) < lightVec.Length()
			return isNoSelfIntersection && isNotBehindLight
		}
		_, hasShadowHit := world.Hits(lightRay, isValidShadowHit)

		if !hasShadowHit {
			lightFactor = lightVec.Normalize().Dot(hit.Normal)
		}

		lightFactors = append(lightFactors, lightFactor)
	}

	lightFactor := avgLightFactor(lightFactors)
	shadedColor := hit.Color.MulScalar(lightFactor + ambientFactor)
	return shadedColor
}

func avgLightFactor(lightFactors []float32) float32 {
	var sum float32 = 0
	for _, v := range lightFactors {
		sum += v
	}

	return sum / float32(len(lightFactors))
}

func createImageBuffer(width, height int) [][]primitive.ScalarColor {
	imageBuffer := make([][]primitive.ScalarColor, height)
	for y := range imageBuffer {
		imageBuffer[y] = make([]primitive.ScalarColor, width)
	}

	return imageBuffer
}

func exportBufferToImage(imageBuffer [][]primitive.ScalarColor, img *image.RGBA) {
	for y := range imageBuffer {
		for x := range imageBuffer[y] {
			img.Set(x, y, imageBuffer[y][x].ToRGBA())
		}
	}
}
