package scene

import (
	"github.com/ruegerj/raytracing/common"
	"github.com/ruegerj/raytracing/primitive"
)

type Hit struct {
	Distance  float32
	Point     primitive.Vec3
	Normal    primitive.Vec3
	UV        common.Optional[primitive.Vec2]
	FrontFace bool
	Material  Material
}

type Hitable interface {
	Hits(r primitive.Ray) (*Hit, bool)
}
