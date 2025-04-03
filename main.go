package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/render"
	"github.com/ruegerj/raytracing/scene"
	"github.com/ruegerj/raytracing/scene/imprt"
)

func main() {

	pathArg := flag.String("path", "", "path to a .gltf file to import")
	flag.Parse()
	if pathArg == nil || *pathArg == "" {
		fmt.Println("Please provide a valid path...")
		os.Exit(1)
	}

	img := image.NewRGBA(image.Rect(0, 0, int(config.WIDTH), int(config.HEIGHT)))

	world, err := imprt.FromGLTF(*pathArg)
	if err != nil {
		panic(err)
	}
	log.Println("imported world from: ", *pathArg)
	// world := create3dCircleWorld()
	render.Do(world, img)

	f, err := os.Create("out/out.jpeg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)
}

func create3dCircleWorld() *scene.World {
	light := scene.NewLight(
		primitive.Vector{X: -3, Y: 2, Z: -1.0}, // 7, 18
		primitive.ScalarColor{R: 1, G: 1, B: 1},
		1.0,
	)
	cam := scene.NewCamera(config.WIDTH/config.HEIGHT, config.DEFAULT_FOV, mgl32.Ident4())
	world := scene.NewWorld([]scene.Hitable{}, []scene.Light{light}, cam)
	var radius float32 = 0.25

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
