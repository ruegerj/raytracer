package render

import (
	"fmt"
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
	imageBuffer := make([][]primitive.ScalarColor, height)
	log.Println(fmt.Sprintf("rendering image: %dx%d", img.Bounds().Dx(), img.Bounds().Dy()))

	var wg sync.WaitGroup
	wg.Add(height)

	for y := range height {
		go func() {
			defer wg.Done()
			imageBuffer[y] = renderLine(y, width, world)
		}()
	}

	wg.Wait()

	log.Println("traced")
	exportBufferToImage(imageBuffer, img)
	log.Println("done")
}

func renderLine(y int, width int, world *scene.World) []primitive.ScalarColor {
	line := make([]primitive.ScalarColor, width)

	for x := range width {
		color := primitive.BLACK

		ray := world.Camera().RayFrom(x, y)
		colorPart := trace(ray, config.MAX_DEPTH, world)
		color = color.Add(colorPart)

		line[x] = color.GammaCorrect()
	}

	return line
}

func trace(ray primitive.Ray, depth float32, world *scene.World) primitive.ScalarColor {
	if depth < config.EPSILON {
		return primitive.BLACK
	}

	potentialHit := world.Hits(ray)
	if potentialHit.IsEmpty() {
		return primitive.BLACK
	}

	hit := potentialHit.Get()
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

func exportBufferToImage(imageBuffer [][]primitive.ScalarColor, img *image.RGBA) {
	for y := range imageBuffer {
		for x := range imageBuffer[y] {
			img.Set(x, y, imageBuffer[y][x].ToRGBA())
		}
	}
}
