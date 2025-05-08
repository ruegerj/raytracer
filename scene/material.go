package scene

import (
	"github.com/ruegerj/raytracing/common"
	"github.com/ruegerj/raytracing/common/optional"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
)

var _ Scaterable = (Material)(nil)

type Material interface {
	Scaterable
	Color() primitive.ScalarColor
}

var _ Material = (*Phong)(nil)

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

func (p *Phong) Scatter(ray primitive.Ray, hit *Hit, world *World) (optional.Optional[primitive.Ray], primitive.ScalarColor) {
	newColor := p.color.MulScalar(config.AMBIENT_FACTOR)
	reflectionDir := ray.Direction().Reflect(hit.Normal).Normalize()

	for _, light := range world.Lights() {
		lightVec := light.Origin.Sub(hit.Point)
		lightDistance := lightVec.Length()
		lightDir := lightVec.Normalize()

		shadowRay := primitive.NewRay(
			hit.Point.Add(lightDir.MulScalar(config.EPSILON)),
			lightDir,
		)

		isValidShadowHit := func(elemHit *Hit, elem Hitable) bool {
			isNoSelfIntersection := elemHit.Distance > config.EPSILON &&
				shadowRay.Direction().Dot(elemHit.Normal) <= 0
			isNotBehindLight := elemHit.Distance < lightVec.Length()
			return isNoSelfIntersection && isNotBehindLight
		}
		_, hasShadowHit := world.Hits(shadowRay, isValidShadowHit)

		if hasShadowHit {
			continue
		}

		lightIntensity := calcDepthBasedLightIntensity(light, lightDistance)

		s := light.Origin.Sub(hit.Point)
		diffuse := p.color.
			MulScalar(max(s.Dot(hit.Normal), 0.0)).
			MulScalar(common.Recip(common.Pow(lightDistance, 2))).
			Mul(light.Color.MulScalar(lightIntensity))

		specularExp := (1.0 - p.roughness) * 128.0
		specular := light.Color.
			MulScalar(1.0 - p.roughness).
			MulScalar(max(common.Pow(reflectionDir.Dot(lightDir), specularExp), 0.0))

		newColor = newColor.Add(diffuse.Add(specular))
	}

	return optional.None[primitive.Ray](), newColor
}

var _ Material = (*Metal)(nil)

type Metal struct {
	color       primitive.ScalarColor
	metalicness float32
}

func NewMetal(color primitive.ScalarColor) *Metal {
	return &Metal{
		color: color,
	}
}

func (m *Metal) Color() primitive.ScalarColor {
	return m.color
}

func (m *Metal) Scatter(ray primitive.Ray, hit *Hit, world *World) (optional.Optional[primitive.Ray], primitive.ScalarColor) {
	reflectionDir := ray.Direction().Reflect(hit.Normal).Normalize()
	reflectionRay := primitive.NewRay(
		hit.Point.Add(reflectionDir.MulScalar(config.EPSILON)),
		reflectionDir,
	)

	return optional.Some(reflectionRay), m.color
}

func calcDepthBasedLightIntensity(light Light, distance float32) float32 {
	intensityFactor := config.DEPTH_LIGHT_A_FACTOR*common.Pow(distance, 2) +
		config.DEPTH_LIGHT_B_FACTOR*distance +
		config.DEPTH_LIGHT_C_FACTOR

	return light.Intensity / intensityFactor
}
