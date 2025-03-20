package render

import "github.com/ruegerj/raytracing/primitive"

var rayDirection = primitive.Vector{X: 0.0, Y: 0.0, Z: -1.0}

type Camera struct {
	halfWidth     float64
	halfHeight    float64
	meterPerPixel float64
}

func NewCamera(width, height, size int) Camera {
	return Camera{
		halfWidth:     float64(width) / 2.0,
		halfHeight:    float64(height) / 2.0,
		meterPerPixel: float64(size) / float64(height),
	}
}

func (c Camera) RayFrom(x, y int) primitive.Ray {
	worldX := (float64(x) - c.halfWidth) * c.meterPerPixel
	worldY := (c.halfHeight - float64(y)) * c.meterPerPixel
	return primitive.Ray{
		Origin:    primitive.Vector{X: worldX, Y: worldY, Z: 0.0},
		Direction: rayDirection,
	}
}
