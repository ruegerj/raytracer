package render

import (
	"image"
	"log"
	"sync"

	"github.com/ruegerj/raytracing/common"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/scene"
)

var DEFAULT_COLOR = primitive.ScalarColor{R: 0, G: 1, B: 1}

func Do(world *scene.World, img *image.RGBA) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	imageBuffer := createImageBuffer(width, height)

	var wg sync.WaitGroup
	wg.Add(height)

	for y := range height {
		go func() {
			defer wg.Done()

			for x := range width {
				ray := world.Camera().RayFrom(x, y)
				color := trace(ray, config.MAX_DEPTH, world)
				imageBuffer[y][x] = color.GammaCorrect()
			}
		}()
	}

	wg.Wait()

	exportBufferToImage(imageBuffer, img)
	log.Println("done")
}

func trace(ray primitive.Ray, depth float32, world *scene.World) primitive.ScalarColor {
	if depth < config.EPSILON {
		return primitive.BLACK
	}

	acceptAnyHit := func(_ *scene.Hit, _ scene.Hitable) bool { return true }
	hit, hasHit := world.Hits(ray, acceptAnyHit)
	if !hasHit {
		return primitive.BLACK
	}

	if hit.Material == nil {
		return DEFAULT_COLOR
	}

	reflectedRay, color := hit.Material.Scatter(ray, hit, world)
	if reflectedRay.IsEmpty() {
		return color
	}

	nextColor := trace(reflectedRay.Get(), depth-1, world)
	return color.Mul(correctColorForDepth(nextColor, depth))
}

func correctColorForDepth(color primitive.ScalarColor, depth float32) primitive.ScalarColor {
	correctionFactor := common.Pow(config.DEPTH_COLOR_DEGRADING_FACTOR, config.MAX_DEPTH-depth)
	return color.MulScalar(correctionFactor)
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
