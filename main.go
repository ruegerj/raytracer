package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/render"
	"github.com/ruegerj/raytracing/scene/imprt"
)

func main() {
	// profiling:
	// defer profile.Start(profile.CPUProfile, profile.ProfilePath(".")).Stop()
	// defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	pathArg := flag.String("path", "", "path to a .gltf file to import")
	flag.Parse()
	if pathArg == nil || *pathArg == "" {
		fmt.Println("Please provide a valid path...")
		os.Exit(1)
	}

	log.Printf("importing %s...\n", *pathArg)
	img := image.NewRGBA(image.Rect(0, 0, int(config.WIDTH), int(config.HEIGHT)))

	world, err := imprt.FromGLTF(*pathArg)
	if err != nil {
		panic(err)
	}

	log.Println("imported world from: ", *pathArg)
	start := time.Now()
	render.Do(world, img)
	end := time.Now()
	log.Printf("total render time: %dms\n", end.UnixMilli()-start.UnixMilli())

	f, err := os.Create("out/out.jpeg")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)
}
