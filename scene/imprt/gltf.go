package imprt

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/ext/lightspunctual"
	"github.com/qmuntal/gltf/modeler"
	"github.com/ruegerj/raytracing/common/optional"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/scene"
)

var defaultLight = scene.Light{
	Origin:    primitive.Vec3{X: -2.5, Y: 3, Z: 2},
	Color:     primitive.ScalarColor{R: 1, G: 1, B: 1},
	Intensity: 1,
}

func FromGLTF(path string) (*scene.World, error) {
	doc, err := gltf.Open(path)
	if err != nil {
		return nil, err
	}

	materials := loadMaterials(doc)

	triangles, err := loadTriangles(doc, materials)
	if err != nil {
		return nil, err
	}

	cameras, err := loadCameras(doc)
	if err != nil {
		return nil, err
	}

	lightSources, err := loadLightSources(doc)
	if err != nil {
		return nil, err
	}

	world := scene.NewWorld(triangles, lightSources, cameras[0])

	return world, nil
}

func loadTriangles(doc *gltf.Document, materials []scene.Material) ([]scene.Triangle, error) {
	triangles := make([]scene.Triangle, 0)
	for _, node := range doc.Nodes {
		if node.Mesh == nil {
			continue
		}

		mesh := doc.Meshes[*node.Mesh]

		for _, prim := range mesh.Primitives {
			posAccessor := doc.Accessors[prim.Attributes["POSITION"]]
			normalAccessor := doc.Accessors[prim.Attributes["NORMAL"]]
			texCoordsAccessor := doc.Accessors[prim.Attributes["TEXCOORD_0"]]
			indicesAccessor := doc.Accessors[*prim.Indices]

			positions, err := modeler.ReadPosition(doc, posAccessor, nil)
			if err != nil {
				return nil, err
			}
			normals, err := modeler.ReadNormal(doc, normalAccessor, nil)
			if err != nil {
				return nil, err
			}
			indices, err := modeler.ReadIndices(doc, indicesAccessor, nil)
			if err != nil {
				return nil, err
			}
			texCoords, err := modeler.ReadTextureCoord(doc, texCoordsAccessor, nil)
			if err != nil {
				return nil, err
			}

			var material scene.Material
			if prim.Material != nil {
				material = materials[*prim.Material]
			}

			for i := 0; i < len(indices); i += 3 {
				triangle := scene.NewTriangle(
					createVertex(uint(i), indices, positions, normals, texCoords),
					createVertex(uint(i+1), indices, positions, normals, texCoords),
					createVertex(uint(i+2), indices, positions, normals, texCoords),
					material,
				)

				triangles = append(triangles, triangle)
			}
		}
	}
	return triangles, nil
}

func loadCameras(doc *gltf.Document) ([]scene.Camera, error) {
	cameras := []scene.Camera{}

	for _, node := range doc.Nodes {
		if node.Camera == nil {
			continue
		}

		aspectRatio := config.DEFAULT_ASPECT_RATIO
		yFov := config.DEFAULT_FOV

		camInfo := doc.Cameras[*node.Camera]
		if camInfo.Perspective != nil {
			if camInfo.Perspective.AspectRatio != nil {
				aspectRatio = float32(*camInfo.Perspective.AspectRatio)
			}
			yFov = float32(camInfo.Perspective.Yfov)
		}

		translation := node.TranslationOrDefault()
		rotation := node.RotationOrDefault()
		transform := createTransformMatrix(translation, rotation)

		cam := scene.NewCamera(aspectRatio, yFov, transform)

		cameras = append(cameras, cam)
	}

	return cameras, nil
}

func loadLightSources(doc *gltf.Document) ([]scene.Light, error) {
	lightSources := []scene.Light{}

	rawLightData, hasLightData := doc.Extensions[lightspunctual.ExtensionName]
	if !hasLightData {
		lightSources = append(lightSources, defaultLight)
		return lightSources, nil
	}

	lights := rawLightData.(lightspunctual.Lights)

	for _, node := range doc.Nodes {
		rawExtensionData, isLight := node.Extensions[lightspunctual.ExtensionName]
		if !isLight {
			continue
		}

		lightIdx := rawExtensionData.(lightspunctual.LightIndex)
		lightData := lights[lightIdx]

		if lightData.Type != lightspunctual.TypePoint {
			continue
		}

		lightTranslation := node.TranslationOrDefault()
		light := scene.NewLight(
			primitive.Vec3{
				X: float32(lightTranslation[0]),
				Y: float32(lightTranslation[1]),
				Z: float32(lightTranslation[2]),
			},
			primitive.FromSlice(*lightData.Color),
			float32(*lightData.Intensity),
		)
		lightSources = append(lightSources, light)
	}

	return lightSources, nil
}

func loadMaterials(doc *gltf.Document) []scene.Material {
	materials := make([]scene.Material, len(doc.Materials))

	for i, m := range doc.Materials {
		var roughness float32 = 1.0
		var metallicFactor float32 = 0
		color := [4]float64{1, 1, 1}
		if m.PBRMetallicRoughness != nil {
			color = m.PBRMetallicRoughness.BaseColorFactorOrDefault()
			roughness = float32(*m.PBRMetallicRoughness.RoughnessFactor)
			metallicFactor = float32(m.PBRMetallicRoughness.MetallicFactorOrDefault())
		}

		scalarColor := primitive.ScalarColor{
			R: float32(color[0]),
			G: float32(color[1]),
			B: float32(color[2]),
		}

		if metallicFactor >= 1 {
			materials[i] = scene.NewMetal(scalarColor)
			continue
		}

		materials[i] = scene.NewPhong(scalarColor, roughness)
	}

	return materials
}

func createVertex(idx uint, indices []uint32, positions, normals [][3]float32, texCoords [][2]float32) scene.Vertex {
	edgeCoords := positions[indices[idx]]
	edgeNormals := normals[indices[idx]]
	uv := optional.None[primitive.Vec2]()

	if len(texCoords) > int(idx) {
		uvCoords := texCoords[idx]
		uv = optional.Some(primitive.Vec2{X: uvCoords[0], Y: uvCoords[1]})
	}

	return scene.Vertex{
		Point:  primitive.Vec3{X: edgeCoords[0], Y: edgeCoords[1], Z: edgeCoords[2]},
		Normal: primitive.Vec3{X: edgeNormals[0], Y: edgeNormals[1], Z: edgeNormals[2]}.Normalize(),
		UV:     uv,
	}
}

func createTransformMatrix(translation [3]float64, rotation [4]float64) primitive.AffineTransformation {
	rotationMat := mgl32.Ident3()

	if !isEmptyRotation(rotation) {
		quat := mgl32.Quat{
			V: mgl32.Vec3{float32(rotation[0]), float32(rotation[1]), float32(rotation[2])},
			W: float32(rotation[3]),
		}.Normalize()
		rotationMat = quat.Mat4().Mat3()
	}

	return primitive.AffineTransformation{
		Rotation:    rotationMat,
		Translation: mgl32.Vec3{float32(translation[0]), float32(translation[1]), float32(translation[2])},
	}
}

func isEmptyRotation(rotation [4]float64) bool {
	return rotation[0] == 0 && rotation[1] == 0 && rotation[2] == 0 && rotation[3] == 1
}
