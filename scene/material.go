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

func (p *Phong) Scatter(ray primitive.Ray, hit *Hit, world *World) (primitive.Ray, bool, primitive.ScalarColor) {
	newColor := p.color.MulScalar(config.AMBIENT_FACTOR)
	reflectionDir := ray.Direction().Reflect(hit.Normal).Normalize()

	for _, light := range world.Lights() {
		lightVec := light.Origin.Sub(hit.Point)
		lightDistance := lightVec.Length()
		lightDir := lightVec.DivScalar(lightDistance).Normalize()

		shadowRay := primitive.NewRay(
			hit.Point.Add(lightDir.MulScalar(config.EPSILON)),
			lightDir,
		)

		var shadowDist float32 = common.F32_INF
		shadowHit := world.Hits(shadowRay)
		if shadowHit != nil {
			shadowDist = shadowHit.Distance
		}

		if shadowDist < lightDistance {
			continue
		}

		lightIntensity := min(light.Intensity, 1.0) / float32(len(world.Lights()))
		s := light.Origin.Sub(hit.Point)
		diffuse := p.color.
			MulScalar(max(s.Dot(hit.Normal), 0.0)).
			MulScalar(common.Recip(lightDistance)).
			Mul(light.Color).
			MulScalar(lightIntensity)

		specularExp := (1.0 - p.roughness) * 128.0
		specular := light.Color.
			MulScalar(1.0 - p.roughness).
			MulScalar(max(common.Pow(reflectionDir.Dot(lightDir), specularExp), 0.0))

		newColor = newColor.Add(diffuse.Add(specular))
	}

	return primitive.Ray{}, false, newColor.Clamp()
}

var _ Material = (*Metal)(nil)

type Metal struct {
	color primitive.ScalarColor
}

func NewMetal(color primitive.ScalarColor) *Metal {
	return &Metal{
		color: color,
	}
}

func (m *Metal) Scatter(ray primitive.Ray, hit *Hit, world *World) (primitive.Ray, bool, primitive.ScalarColor) {
	reflectionDir := ray.Direction().Reflect(hit.Normal).Normalize()
	reflectionRay := primitive.NewRay(
		hit.Point.Add(reflectionDir.MulScalar(config.EPSILON)),
		reflectionDir,
	)

	return reflectionRay, true, m.color
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

func calcDepthBasedLightIntensity(light Light, distance float32) float32 {
	intensityFactor := config.DEPTH_LIGHT_A_FACTOR*common.Pow(distance, 2) +
		config.DEPTH_LIGHT_B_FACTOR*distance +
		config.DEPTH_LIGHT_C_FACTOR

	return min(light.Intensity, 1.0) / intensityFactor
}
