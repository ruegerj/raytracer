package primitive

type Material struct {
	color ScalarColor
}

func NewMaterial(color ScalarColor) *Material {
	return &Material{
		color: color,
	}
}

func (m *Material) Color() ScalarColor {
	return m.color
}
