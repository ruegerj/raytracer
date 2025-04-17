package scene

import (
	"github.com/ruegerj/raytracing/primitive"
)

type Hit struct {
	Distance float32
	Point    primitive.Vector
	Normal   primitive.Vector
	Material Material
}

type Hitable interface {
	Hits(r primitive.Ray) (*Hit, bool)
}
