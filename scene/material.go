package scene

import (
	"github.com/ruegerj/raytracing/primitive"
)

type Material interface {
	Color() primitive.ScalarColor
}

type Phong struct {
	color     primitive.ScalarColor
	roughness float32
}

func NewPhong(color primitive.ScalarColor, roughness float32) *Phong {
	return &Phong{
		color:     color,
		roughness: roughness,
	}
}

func (p *Phong) Color() primitive.ScalarColor {
	return p.color
}

func (p *Phong) Roughness() float32 {
	return p.roughness
}

type Metallic struct {
	color       primitive.ScalarColor
	metalicness float32
}

func NewMetallic(color primitive.ScalarColor) *Metallic {
	return &Metallic{
		color: color,
	}
}

func (m *Metallic) Color() primitive.ScalarColor {
	return m.color
}
