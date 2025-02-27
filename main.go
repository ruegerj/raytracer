package main

import (
	"image"
	"image/jpeg"
	"math"
	"os"

	"github.com/ruegerj/raytracing/primitive"
)

func main() {
	const height = 400
	const width = 600

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	c := primitive.Vector{
		X: width / 2,
		Y: height / 2,
		Z: 0,
	}
	r := 0.9 * min(height, width) / 2

	for row := range height {
		for col := range width {
			p := primitive.Vector{X: float64(row) / 2, Y: float64(col) / 2, Z: 0}
			if math.Abs(p.Sub(c).Length()) <= r {
				blue := primitive.ScalarColor{R: 0, G: 0, B: 1}
				img.Set(col, row, blue.ToRGBA())
			}
		}
	}

	f, err := os.Create("out/out.jpeg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)

}
