package scene

import (
	"math"

	"github.com/ruegerj/raytracing/primitive"
)

type BvhNode struct {
	aabb      *primitive.AABB
	leftChild uint
	firstTri  uint
	triCount  uint
}

func NewBvhNode(firstPrim, primCount uint, triangles []Triangle) BvhNode {
	aabb := primitive.MAX_AABB()

	for _, tri := range triangles {
		aabb.Grow(tri.V0.Point)
		aabb.Grow(tri.V1.Point)
		aabb.Grow(tri.V2.Point)
	}

	return BvhNode{
		aabb:      aabb,
		leftChild: 0,
		firstTri:  firstPrim,
		triCount:  primCount,
	}
}

func (n BvhNode) IsLeaf() bool {
	return n.triCount > 0
}

func (n BvhNode) GetOwnTriangles(triangles []Triangle) []Triangle {
	return triangles[n.firstTri : n.firstTri+n.triCount]
}

func (n BvhNode) EvaluateSAH(axis uint, pos float32, triangles []Triangle) float32 {
	leftBox := primitive.MAX_AABB()
	rightBox := primitive.MAX_AABB()

	leftCount := 0
	rightCount := 0

	for _, tri := range n.GetOwnTriangles(triangles) {
		if tri.Centroid.Axis(axis) < pos {
			leftCount++
			leftBox.Grow(tri.V0.Point)
			leftBox.Grow(tri.V1.Point)
			leftBox.Grow(tri.V2.Point)
		} else {
			rightCount++
			rightBox.Grow(tri.V0.Point)
			rightBox.Grow(tri.V1.Point)
			rightBox.Grow(tri.V2.Point)
		}
	}

	cost := float32(leftCount)*leftBox.Area() + float32(rightCount)*rightBox.Area()
	if cost <= 0.0 {
		return float32(math.Inf(1))
	}

	return cost
}
