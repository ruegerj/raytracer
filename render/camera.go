package render

import "github.com/ruegerj/raytracing/primitive"

var rayDirection = primitive.Vector{X: 0.0, Y: 0.0, Z: -1.0}

type Camera struct {
	halfWidth     float32
	halfHeight    float32
	meterPerPixel float32
}

func NewCamera(width, height, size int) Camera {
	return Camera{
		halfWidth:     float32(width) / 2.0,
		halfHeight:    float32(height) / 2.0,
		meterPerPixel: float32(size) / float32(height),
	}
}

func (c Camera) RayFrom(x, y int) primitive.Ray {
	worldX := (float32(x) - c.halfWidth) * c.meterPerPixel
	worldY := (c.halfHeight - float32(y)) * c.meterPerPixel
	return primitive.Ray{
		Origin:    primitive.Vector{X: worldX, Y: worldY, Z: 0.0},
		Direction: rayDirection,
	}
}
