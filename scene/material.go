package scene

import (
	"math"
	"math/rand"

	"github.com/ruegerj/raytracing/common"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
)

const glass_ior float32 = 1.52

type Material interface {
	Scatter(ray primitive.Ray, hit *Hit, world *World) (primitive.Ray, bool, primitive.ScalarColor)
}

var _ Material = (*Diffuse)(nil)

type Diffuse struct {
	color primitive.ScalarColor
}

func NewDiffuse(color primitive.ScalarColor) *Diffuse {
	return &Diffuse{color: color}
}

func (d *Diffuse) Scatter(ray primitive.Ray, hit *Hit, world *World) (primitive.Ray, bool, primitive.ScalarColor) {
	rayDir := (hit.Normal.Add(primitive.RandomUnitVector())).Normalize()
	rayOrigin := hit.Point.Add(rayDir.MulScalar(config.EPSILON))
	return primitive.NewRay(rayOrigin, rayDir), true, d.color
}

var _ Material = (*Metal)(nil)

type Metal struct {
	color     primitive.ScalarColor
	roughness float32
}

func NewMetal(color primitive.ScalarColor, roughness float32) *Metal {
	return &Metal{
		color:     color,
		roughness: roughness,
	}
}

func (m *Metal) Scatter(ray primitive.Ray, hit *Hit, world *World) (primitive.Ray, bool, primitive.ScalarColor) {
	randomReflection := primitive.RandomOnHemisphere(hit.Normal).MulScalar(m.roughness)
	reflectionDir := ray.Direction().Reflect(hit.Normal).Add(randomReflection).Normalize()
	reflectionOrigin := hit.Point.Add(reflectionDir.MulScalar(config.EPSILON))
	return primitive.NewRay(reflectionOrigin, reflectionDir), true, m.color
}

var _ Material = (*Glass)(nil)

type Glass struct {
	color primitive.ScalarColor
}

func NewGlass(color primitive.ScalarColor) *Glass {
	return &Glass{color: color}
}

func (g *Glass) Scatter(ray primitive.Ray, hit *Hit, world *World) (primitive.Ray, bool, primitive.ScalarColor) {
	eta := glass_ior
	if hit.FrontFace {
		eta = common.Recip(eta)
	}

	cosTheta := min(ray.Direction().Negate().Dot(hit.Normal), 1.0)
	sinTheta := float32(math.Sqrt(float64(1.0 - cosTheta*cosTheta)))
	cannotRefract := eta*sinTheta > 1.0

	var targetDir primitive.Vec3
	doesReflect := cannotRefract || reflectanceSchlick(cosTheta, eta) > rand.Float32()
	if doesReflect {
		targetDir = ray.Direction().Reflect(hit.Normal)
	} else {
		targetDir = ray.Direction().Refract(hit.Normal, eta)
	}

	rayOrigin := hit.Point.Add(targetDir.MulScalar(config.EPSILON))
	return primitive.NewRay(rayOrigin, targetDir), true, g.color
}

var _ Material = (*Emissive)(nil)

type Emissive struct {
	color primitive.ScalarColor
}

func NewEmissive(color primitive.ScalarColor) *Emissive {
	return &Emissive{color: color}
}

func (e *Emissive) Scatter(ray primitive.Ray, hit *Hit, world *World) (primitive.Ray, bool, primitive.ScalarColor) {
	return primitive.Ray{}, false, e.color
}

func reflectanceSchlick(cosine, ior float32) float32 {
	r0 := common.Pow((1.0-ior)/(1.0+ior), 2.0)
	approx := common.Pow(r0+(1.0-r0)*(1.0-cosine), 5.0)
	return approx
}
