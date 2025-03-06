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
	const height = 720
	const width = 1280
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

func create3dCircleWorld() *shape.World {
	world := &shape.World{}
	var radius float64 = 200

	red := shape.NewSphere(primitive.Vector{X: 640, Y: 280, Z: -500}, radius, primitive.ScalarColor{R: 1, G: 0, B: 0})
	green := shape.NewSphere(primitive.Vector{X: 520, Y: 440, Z: -600}, radius, primitive.ScalarColor{R: 0, G: 1, B: 0})
	blue := shape.NewSphere(primitive.Vector{X: 760, Y: 440, Z: -700}, radius, primitive.ScalarColor{R: 0, G: 0, B: 1})

	world.AddAll(red, green, blue)

	return world
}
