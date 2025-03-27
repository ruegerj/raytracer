package scene

import (
	"github.com/qmuntal/gltf"
	"github.com/qmuntal/gltf/modeler"
	"github.com/ruegerj/raytracing/primitive"
)

// TODO: import from glTF
var defaultColor = primitive.ScalarColor{R: 0, G: 1, B: 1}
var tempLight = Light{
	Origin:    primitive.Vector{X: 3, Y: 1, Z: -1},
	Color:     primitive.ScalarColor{R: 1, G: 1, B: 1},
	Intensity: 1,
}

func ImportFromGLTF(path string) (*World, error) {
	doc, err := gltf.Open(path)
	if err != nil {
		return nil, err
	}

	triangles, err := loadTriangles(doc)
	if err != nil {
		return nil, err
	}
	world := NewWorld(triangles, []Light{tempLight})

	return world, nil
}

func loadTriangles(doc *gltf.Document) ([]Hitable, error) {
	triangles := make([]Hitable, 0)
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

				triangle := NewTriangle(
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

func createVertex(coords, normals [3]float32) Vertex {
	return Vertex{
		Point:  primitive.Vector{X: coords[0], Y: coords[1], Z: coords[2]},
		Normal: primitive.Vector{X: normals[0], Y: normals[1], Z: normals[2]},
	}
}
