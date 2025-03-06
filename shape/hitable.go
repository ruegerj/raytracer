package shape

import "github.com/ruegerj/raytracing/primitive"

type Hit struct {
	Distance float64
	Element  Hitable
}

type Hitable interface {
	Color() primitive.ScalarColor
	Hits(r primitive.Ray) (*Hit, bool)
}
