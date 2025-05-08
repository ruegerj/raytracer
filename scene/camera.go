package scene

import (
	"math"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/ruegerj/raytracing/common"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
)

var rayDirection = primitive.Vec3{X: 0.0, Y: 0.0, Z: -1.0}

type Camera struct {
	halfWidth     float32
	halfHeight    float32
	meterPerPixel float32
	focalLength   float32
	transform     primitive.AffineTransformation
}

func NewCamera(aspectRatio, yFov float32, transform primitive.AffineTransformation) Camera {
	h := common.Recip(aspectRatio)

	pixelHeight := config.HEIGHT
	pixelWidth := config.WIDTH

	return Camera{
		halfWidth:     pixelWidth / 2,
		halfHeight:    pixelHeight / 2,
		meterPerPixel: h / config.HEIGHT,
		focalLength:   calcFocalLenght(h, yFov),
		transform:     transform,
	}
}

func (c Camera) RayFrom(x, y int) primitive.Ray {
	planeX := (float32(x) - c.halfWidth) * c.meterPerPixel
	planeY := (c.halfHeight - float32(y)) * c.meterPerPixel

	direction := mgl32.Vec3{planeX, planeY, -c.focalLength}.Normalize()
	rotatedDirection := c.transform.Rotation.Mul3x1(direction)

	return primitive.NewRay(
		vec3ToVector(c.transform.Translation),
		vec3ToVector(rotatedDirection).Normalize(),
	)
}

func calcFocalLenght(height, yFov float32) float32 {
	return (height / 2) / float32(math.Tan(float64(yFov/2)))
}

func vec3ToVector(v mgl32.Vec3) primitive.Vec3 {
	return primitive.Vec3{X: v.X(), Y: v.Y(), Z: v.Z()}
}
