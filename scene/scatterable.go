package scene

import (
	"github.com/ruegerj/raytracing/common"
	"github.com/ruegerj/raytracing/primitive"
)

type Scaterable interface {
	Scatter(ray primitive.Ray, hit *Hit, world *World) (common.Optional[primitive.Ray], primitive.ScalarColor)
}
