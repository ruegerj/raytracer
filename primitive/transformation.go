package primitive

import "github.com/go-gl/mathgl/mgl32"

type AffineTransformation struct {
	Translation mgl32.Vec3
	Rotation    mgl32.Mat3
}
