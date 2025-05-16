package scene

import (
	"github.com/ruegerj/raytracing/common"
	"github.com/ruegerj/raytracing/config"
	"github.com/ruegerj/raytracing/primitive"
)

const ROOT_INDEX uint = 0

type Bvh struct {
	nodes     []BvhNode
	triangles []Triangle
	usedNodes uint
}

func NewBvh(triangles []Triangle) *Bvh {
	root := NewBvhNode(0, uint(len(triangles)), triangles)
	nodes := []BvhNode{root}

	bvh := &Bvh{
		nodes:     nodes,
		triangles: triangles,
		usedNodes: 1,
	}

	bvh.Subdivide(ROOT_INDEX)

	return bvh
}

func (b *Bvh) Intersects(ray primitive.Ray) *Hit {
	node := &b.nodes[ROOT_INDEX]
	stack := [64]*BvhNode{node}
	stackPointer := 0

	var nearestDist float32 = common.F32_INF
	var nearestTriangle *Triangle = nil

	for {
		if node.IsLeaf() {
			for i := node.firstTri; i < node.firstTri+node.triCount; i++ {
				tri := &b.triangles[i]
				hit := tri.Hits(ray)
				if hit == nil {
					continue
				}

				if hit.Distance < nearestDist {
					nearestDist = hit.Distance
					nearestTriangle = tri
				}

			}

			if stackPointer == 0 {
				break
			}

			stackPointer--
			node = stack[stackPointer]
			continue
		}

		child1 := &b.nodes[node.leftChild]
		child2 := &b.nodes[node.leftChild+1]

		dist1 := common.F32_INF
		if hitDist := child1.aabb.Hit(ray); hitDist < nearestDist {
			dist1 = hitDist
		}

		dist2 := common.F32_INF
		if hitDist := child2.aabb.Hit(ray); hitDist < nearestDist {
			dist2 = hitDist
		}

		if dist1 > dist2 {
			dist1, dist2 = dist2, dist1
			child1, child2 = child2, child1
		}

		if dist1 == common.F32_INF {
			if stackPointer == 0 {
				break
			}

			stackPointer--
			node = stack[stackPointer]
		} else {
			node = child1
			if dist2 != common.F32_INF {
				stack[stackPointer] = child2
				stackPointer++
			}
		}
	}

	if nearestTriangle == nil {
		return nil
	}

	hit := nearestTriangle.CreateHitFor(ray, nearestDist)
	return hit
}

func (b *Bvh) Subdivide(nodeIndex uint) {
	node := &b.nodes[nodeIndex]

	var bestAxis uint = 3
	var bestPos float32 = 0.0
	var bestCost float32 = common.F32_INF

	for axis := range 3 {
		boundsMin := node.aabb.Minimum.Axis(uint(axis))
		boundsMax := node.aabb.Maximum.Axis(uint(axis))

		if boundsMin == boundsMax {
			continue
		}

		scale := (boundsMax - boundsMin) / config.BVH_SPACES
		for i := 1; i < config.BVH_SPACES; i++ {
			candidatePos := boundsMin + float32(i)*scale
			cost := node.EvaluateSAH(uint(axis), candidatePos, b.triangles)
			if cost < bestCost {
				bestPos = candidatePos
				bestAxis = uint(axis)
				bestCost = cost
			}
		}
	}

	extent := node.aabb.Maximum.Sub(node.aabb.Minimum)
	parentArea := extent.X*extent.Y + extent.Y*extent.Z + extent.Z*extent.X
	parentCost := float32(node.triCount) * parentArea

	if bestCost >= parentCost {
		return
	}

	axis := bestAxis
	splitPos := bestPos

	i := node.firstTri
	j := i + node.triCount - 1

	for i <= j {
		if b.triangles[i].Centroid.Axis(axis) < splitPos {
			i++
		} else {
			common.Swap(b.triangles, i, j)
			j--
		}
	}

	leftCount := i - node.firstTri
	if leftCount == 0 || leftCount == node.triCount {
		return
	}

	leftChildIndex := b.usedNodes
	rightChildIndex := b.usedNodes + 1
	b.usedNodes += 2

	firstPrim := node.firstTri
	primCount := node.triCount

	node.leftChild = leftChildIndex
	node.triCount = 0

	leftChild := NewBvhNode(firstPrim, leftCount, b.triangles)
	b.nodes = append(b.nodes, leftChild)
	rightChild := NewBvhNode(i, primCount-leftCount, b.triangles)
	b.nodes = append(b.nodes, rightChild)

	b.Subdivide(leftChildIndex)
	b.Subdivide(rightChildIndex)
}
