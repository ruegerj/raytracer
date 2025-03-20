package main

import (
	"image"
	"image/jpeg"
	"os"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/render"
	"github.com/ruegerj/raytracing/shape"
)

func main() {
	const height = 1080
	const width = 1920
	const depth float64 = 1000

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	world := create3dCircleWorld(img.Bounds().Dx(), img.Bounds().Dy())
	render.Do(world, img, depth)

	f, err := os.Create("out/out.jpeg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)

}

func create3dCircleWorld(width, height int) *shape.World {
	world := &shape.World{}
	var radius float64 = 0.25

	redSphere := shape.NewSphere(
		primitive.Vector{X: -0.575, Y: 0, Z: -1.0},
		radius,
		primitive.ScalarColor{R: 1, G: 0, B: 0},
	)
	greenSphere := shape.NewSphere(
		primitive.Vector{X: 0, Y: 0, Z: -1.0},
		radius,
		primitive.ScalarColor{R: 0, G: 1, B: 0},
	)
	blueSphere := shape.NewSphere(
		primitive.Vector{X: 0.575, Y: 0, Z: -1.0},
		radius,
		primitive.ScalarColor{R: 0, G: 0, B: 1},
	)

	world.AddAll(redSphere, greenSphere, blueSphere)

	return world
}

// epsilon: 10^-9
// start ray: hitpoint + vector to light * epsilon
