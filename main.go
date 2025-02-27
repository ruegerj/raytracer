package main

import (
	"image"
	"image/jpeg"
	"os"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/shape"
)

func main() {
	const height = 720
	const width = 1280

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	renderCircle(img)
	// render3Circles(img)

	f, err := os.Create("out/out.jpeg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)

}

func renderCircle(img *image.RGBA) {
	width, height := getDimensions(img)

	c := primitive.Vector{
		X: float64(width) / 2,
		Y: float64(height) / 2,
		Z: 0,
	}
	// r := 0.9 * min(float64(height), float64(width)) / float64(2)
	circle := shape.Circle{
		Origin: c,
		Radius: 200,
		Color:  primitive.ScalarColor{R: 0, G: 0, B: 1},
	}

	for y := range height {
		for x := range width {
			p := primitive.Vector{X: float64(x), Y: float64(y), Z: 0}
			if circle.Hits(p) {
				img.Set(x, y, circle.Color.ToRGBA())
			}
		}
	}
}

func render3Circles(img *image.RGBA) {
	width, height := getDimensions(img)
	var radius float64 = 200

	red := shape.NewCircle(primitive.Vector{X: 640, Y: 280, Z: 0}, radius, primitive.ScalarColor{R: 1, G: 0, B: 0})
	green := shape.NewCircle(primitive.Vector{X: 520, Y: 440, Z: 0}, radius, primitive.ScalarColor{R: 0, G: 1, B: 0})
	blue := shape.NewCircle(primitive.Vector{X: 760, Y: 440, Z: 0}, radius, primitive.ScalarColor{R: 0, G: 0, B: 1})

	for y := range height {
		for x := range width {
			p := primitive.Vector{X: float64(x), Y: float64(y), Z: 0}
			color := primitive.ScalarColor{R: 0, G: 0, B: 0}
			if red.Hits(p) {
				color = color.Add(red.Color)
			}
			if green.Hits(p) {
				color = color.Add(green.Color)
			}
			if blue.Hits(p) {
				color = color.Add(blue.Color)
			}

			img.Set(x, y, color.ToRGBA())
		}
	}
}

func getDimensions(img *image.RGBA) (int, int) {
	return img.Bounds().Dx(), img.Bounds().Dy()
}
