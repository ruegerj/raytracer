package imprt

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
	"github.com/ruegerj/raytracing/scene"
)

// TODO: import from glTF
var defaultColor = primitive.ScalarColor{R: 0, G: 1, B: 1}
var tempLight = scene.Light{
	Origin:    primitive.Vector{X: -2.6, Y: 3, Z: -2.1},
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

	world := scene.NewWorld(triangles, []scene.Light{tempLight}, cameras[0])

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
		rotation := resolveRotationOrDefaultOf(node)
		transform := createTransformMatrix(translation, rotation)

		cam := scene.NewCamera(aspectRatio, yFov, transform)

		cameras = append(cameras, cam)
	}

	return cameras, nil
}

func createVertex(coords, normals [3]float32) scene.Vertex {
	return scene.Vertex{
		Point:  primitive.Vector{X: coords[0], Y: coords[1], Z: coords[2]},
		Normal: primitive.Vector{X: normals[0], Y: normals[1], Z: normals[2]},
	}
}

func createTransformMatrix(translation [3]float64, rotation [4]float64) mgl32.Mat4 {
	// Convert quaternion (glTF stores it as [x, y, z, w]) to rotation matrix
	quat := mgl32.Quat{
		V: mgl32.Vec3{float32(rotation[0]), float32(rotation[1]), float32(rotation[2])},
		W: float32(rotation[3]),
	}
	rotationMat := quat.Mat4()

	translationMat := mgl32.Translate3D(float32(translation[0]), float32(translation[1]), float32(translation[2]))

	return translationMat.Mul4(rotationMat)
}

func resolveRotationOrDefaultOf(node *gltf.Node) [4]float64 {
	rotation := node.Rotation

	if isEmptyRotation(rotation) {
		rotation = [4]float64{0.707, 0, 0, -0.707} // align to -z axis
	}

	return rotation
}

func isEmptyRotation(rotation [4]float64) bool {
	return rotation[0] == 0 && rotation[1] == 0 && rotation[2] == 0 && rotation[3] == 1
}
