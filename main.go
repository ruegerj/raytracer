package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/shape"
)

func main() {
	const height = 720
	const width = 1280

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	world := create3dCircleWorld()
	render(world, img)

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

	red := shape.NewSphere(primitive.Vector{X: 640, Y: 280, Z: -100}, radius, primitive.ScalarColor{R: 1, G: 0, B: 0})
	green := shape.NewSphere(primitive.Vector{X: 520, Y: 440, Z: -200}, radius, primitive.ScalarColor{R: 0, G: 1, B: 0})
	blue := shape.NewSphere(primitive.Vector{X: 760, Y: 440, Z: -50}, radius, primitive.ScalarColor{R: 0, G: 0, B: 1})

	world.AddAll(red, green, blue)

	return world
}

func render(target shape.Hitable, img *image.RGBA) {
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

			img.Set(x, y, hit.Element.Color().ToRGBA())
		}
	}
}
