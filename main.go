package main

import (
	"image"
	"image/jpeg"
	"os"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/render"
	"github.com/ruegerj/raytracing/scene"
)

func main() {
	const height = 1080
	const width = 1920
	const depth float64 = 1000

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	world := create3dCircleWorld()
	render.Do(world, img, depth)

	f, err := os.Create("out/out.jpeg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)
}

func create3dCircleWorld() *scene.World {
	light := scene.NewLight(
		primitive.Vector{X: -10, Y: 7, Z: 7.5}, // 7, 18
		primitive.ScalarColor{R: 1, G: 1, B: 1},
		1.0,
	)
	world := scene.NewWorld([]scene.Hitable{}, []scene.Light{light})
	var radius float64 = 0.25

	redSphere := scene.NewSphere(
		primitive.Vector{X: -0.575, Y: 0, Z: -1.0},
		radius,
		primitive.ScalarColor{R: 1, G: 0, B: 0},
	)
	greenSphere := scene.NewSphere(
		primitive.Vector{X: 0, Y: 0, Z: -1.0},
		radius,
		primitive.ScalarColor{R: 0, G: 1, B: 0},
	)
	blueSphere := scene.NewSphere(
		primitive.Vector{X: 0.575, Y: 0, Z: -1.0},
		radius,
		primitive.ScalarColor{R: 0, G: 0, B: 1},
	)

	world.AddAll(redSphere, greenSphere, blueSphere)

	return world
}
