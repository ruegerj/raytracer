package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

func main() {
	const HEIGHT = 400
	const WIDTH = 600

	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	for row := 0; row < HEIGHT; row++ {
		for col := 0; col < WIDTH; col++ {
			renderCircle(col, row, WIDTH, HEIGHT, img)
		}
	}

	f, err := os.Create("out/out.jpeg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)

}

func renderCircle(x, y, width, height int, img *image.RGBA) {
	const r = 100
	cx := x - width/2
	cy := y - height/2

	if cx*cx+cy*cy <= r*r {
		img.Set(x, y, color.RGBA{0, 0, 255, 1})
	}
}
