package scene

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
)

var rayDirection = primitive.Vector{X: 0.0, Y: 0.0, Z: -1.0}

type Camera struct {
	halfWidth     float32
	halfHeight    float32
	meterPerPixel float32
	focalLength   float32
	transform     mgl32.Mat4
}

func NewCamera(aspectRatio, yFov float32, transform mgl32.Mat4) Camera {
	h := 1 / aspectRatio
	return Camera{
		halfWidth:     config.WIDTH / 2,
		halfHeight:    config.HEIGHT / 2,
		meterPerPixel: aspectRatio / config.HEIGHT,
		focalLength:   calcFocalLenght(h, yFov),
		transform:     transform,
	}
}

func (c Camera) RayFrom(x, y int) primitive.Ray {
	planeX := (float32(x) - c.halfWidth) * c.meterPerPixel
	planeY := (c.halfHeight - float32(y)) * c.meterPerPixel

	cameraPos := c.transform.Col(3).Vec3()
	localDir := primitive.Vector{X: planeX, Y: planeY, Z: -c.focalLength}.Normalize()

	return primitive.Ray{
		Origin:    vec3ToVector(cameraPos),
		Direction: localDir,
	}
}

func calcFocalLenght(height, yFov float32) float32 {
	return (height / 2) / float32(math.Tan(float64(yFov/2)))
}

func vec3ToVector(v mgl32.Vec3) primitive.Vector {
	return primitive.Vector{X: v.X(), Y: v.Y(), Z: v.Z()}
}
