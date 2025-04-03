package imprt

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/ext/lightspunctual"
	"github.com/qmuntal/gltf/modeler"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/scene"
)

// TODO: import from glTF
var defaultColor = primitive.ScalarColor{R: 0, G: 1, B: 1}
var defaultLight = scene.Light{
	Origin:    primitive.Vector{X: -2.5, Y: 3, Z: 2},
	Color:     primitive.ScalarColor{R: 1, G: 1, B: 1},
	Intensity: 1,
}

func FromGLTF(path string) (*scene.World, error) {
	doc, err := gltf.Open(path)
	if err != nil {
		return nil, err
	}

	triangles, err := loadTriangles(doc)
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

func loadTriangles(doc *gltf.Document) ([]scene.Hitable, error) {
	triangles := make([]scene.Hitable, 0)
	for _, node := range doc.Nodes {
		if node.Mesh == nil {
			continue
		}

		mesh := doc.Meshes[*node.Mesh]

		for _, prim := range mesh.Primitives {
			posAccessor := doc.Accessors[prim.Attributes["POSITION"]]
			normalAccessor := doc.Accessors[prim.Attributes["NORMAL"]]
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

			for i := 0; i < len(indices); i += 3 {
				p0 := positions[indices[i]]
				p1 := positions[indices[i+1]]
				p2 := positions[indices[i+2]]
				n0 := normals[indices[i]]
				n1 := normals[indices[i+1]]
				n2 := normals[indices[i+2]]

				triangle := scene.NewTriangle(
					createVertex(p0, n0),
					createVertex(p1, n1),
					createVertex(p2, n2),
					defaultColor,
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
			primitive.Vector{
				X: float32(lightTranslation[0]),
				Y: float32(lightTranslation[1]),
				Z: float32(lightTranslation[2]),
			},
			primitive.FromSlice(*lightData.Color),
			float32(*lightData.Intensity)/1000,
		)
		lightSources = append(lightSources, light)
	}

	return lightSources, nil
}

func createVertex(coords, normals [3]float32) scene.Vertex {
	return scene.Vertex{
		Point:  primitive.Vector{X: coords[0], Y: coords[1], Z: coords[2]},
		Normal: primitive.Vector{X: normals[0], Y: normals[1], Z: normals[2]},
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
