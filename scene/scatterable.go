package scene

import (
	"github.com/ruegerj/raytracing/common/optional"
	"github.com/ruegerj/raytracing/primitive"
)

type Scaterable interface {
	Scatter(ray primitive.Ray, hit *Hit, world *World) (optional.Optional[primitive.Ray], primitive.ScalarColor)
}
