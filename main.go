package main

import (
	"image"
	"image/jpeg"
	"math"
	"os"

	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/shape"
)

func main() {
	const height = 720
	const width = 1280

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// renderCircle(img)
	// render3Circles(img)
	render3dSpheres(img)

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
	circle := shape.Sphere{
		Center: c,
		Radius: 200,
		Color:  primitive.ScalarColor{R: 0, G: 0, B: 1},
	}

	for y := range height {
		for x := range width {
			p := primitive.Vector{X: float64(x), Y: float64(y), Z: 0}
			if circle.HitsVector(p) {
				img.Set(x, y, circle.Color.ToRGBA())
			}
		}
	}
}

func render3Circles(img *image.RGBA) {
	width, height := getDimensions(img)
	var radius float64 = 200

	red := shape.NewSphere(primitive.Vector{X: 640, Y: 280, Z: 0}, radius, primitive.ScalarColor{R: 1, G: 0, B: 0})
	green := shape.NewSphere(primitive.Vector{X: 520, Y: 440, Z: 0}, radius, primitive.ScalarColor{R: 0, G: 1, B: 0})
	blue := shape.NewSphere(primitive.Vector{X: 760, Y: 440, Z: 0}, radius, primitive.ScalarColor{R: 0, G: 0, B: 1})

	for y := range height {
		for x := range width {
			p := primitive.Vector{X: float64(x), Y: float64(y), Z: 0}
			color := primitive.ScalarColor{R: 0, G: 0, B: 0}
			if red.HitsVector(p) {
				color = color.Add(red.Color)
			}
			if green.HitsVector(p) {
				color = color.Add(green.Color)
			}
			if blue.HitsVector(p) {
				color = color.Add(blue.Color)
			}

			img.Set(x, y, color.ToRGBA())
		}
	}
}

func render3dSpheres(img *image.RGBA) {
	width, height := getDimensions(img)
	var radius float64 = 200

	red := shape.NewSphere(primitive.Vector{X: 640, Y: 280, Z: -100}, radius, primitive.ScalarColor{R: 1, G: 0, B: 0})
	green := shape.NewSphere(primitive.Vector{X: 520, Y: 440, Z: -200}, radius, primitive.ScalarColor{R: 0, G: 1, B: 0})
	blue := shape.NewSphere(primitive.Vector{X: 760, Y: 440, Z: -50}, radius, primitive.ScalarColor{R: 0, G: 0, B: 1})

	spheres := []shape.Sphere{green, blue, red}

	for y := range height {
		for x := range width {
			ray := primitive.Ray{
				Origin:    primitive.Vector{X: float64(x), Y: float64(y), Z: 0},
				Direction: primitive.Vector{X: 0, Y: 0, Z: 1},
			}

			minDist := math.MaxFloat64
			color := primitive.ScalarColor{R: 0, G: 0, B: 0}
			for _, sphere := range spheres {
				dist, hits := sphere.Hits(ray)
				if !hits {
					continue
				}

				if dist < minDist {
					minDist = dist
					color = sphere.Color
				}
			}

			if minDist < math.MaxFloat64 {
				img.Set(x, y, color.ToRGBA())
			}
		}
	}
}

func getDimensions(img *image.RGBA) (int, int) {
	return img.Bounds().Dx(), img.Bounds().Dy()
}
