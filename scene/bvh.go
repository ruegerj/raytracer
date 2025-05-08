package scene

import (
	"math"

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

func NewBvh(triangles []Triangle) Bvh {
	root := NewBvhNode(0, uint(len(triangles)), triangles)
	nodes := []BvhNode{root}

	bvh := Bvh{
		nodes:     nodes,
		triangles: triangles,
		usedNodes: 1,
	}

	bvh.Subdivide(ROOT_INDEX)
	return bvh
}

func (b Bvh) Intersects(ray primitive.Ray) common.Optional[Hit] {
	node := b.nodes[ROOT_INDEX]
	stack := [65]BvhNode{node}
	stackPointer := 0

	var nearestDist float32 = math.MaxFloat32
	nearestTriangle := common.Empty[Triangle]()

	for {
		if node.IsLeaf() {
			for _, tri := range node.GetOwnTriangles(b.triangles) {
				hit, hasHit := tri.Hits(ray)
				if !hasHit {
					continue
				}

				if hit.Distance < nearestDist {
					nearestDist = hit.Distance
					nearestTriangle = common.Some(tri)
				}
			}

			if stackPointer == 0 {
				break
			}

			stackPointer--
			node = stack[stackPointer]
			continue
		}

		child1 := b.nodes[node.leftChild]
		child2 := b.nodes[node.leftChild+1]

		dist1 := float32(math.Inf(1))
		if hit := child1.aabb.Hit(ray); hit.IsPresent() && hit.Get() < nearestDist {
			dist1 = hit.Get()
		}

		dist2 := float32(math.Inf(1))
		if hit := child2.aabb.Hit(ray); hit.IsPresent() && hit.Get() < nearestDist {
			dist2 = hit.Get()
		}

		nearDist := common.Empty[float32]()
		farDist := common.Empty[float32]()
		var nearChild, farChild BvhNode

		if dist1 > dist2 {
			nearDist = common.Some(dist2)
			nearChild = child2
			farDist = common.Some(dist1)
			farChild = child1
		} else {
			nearDist = common.Some(dist1)
			nearChild = child1
			farDist = common.Some(dist2)
			farChild = child2
		}

		if nearDist.IsEmpty() {
			if stackPointer == 0 {
				break
			}

			stackPointer--
			node = stack[stackPointer]
		} else {
			node = nearChild
			if farDist.IsPresent() {
				stack[stackPointer] = farChild
				stackPointer++
			}
		}
	}

	if nearestTriangle.IsEmpty() {
		return common.Empty[Hit]()
	}

	hit := nearestTriangle.Get().CreateHitFor(ray, nearestDist)
	return common.Some(hit)
}

func (b Bvh) Subdivide(nodeIndex uint) {
	node := b.nodes[nodeIndex]

	var bestAxis uint = 3
	var bestPos float32 = 0.0
	var bestCost float32 = math.MaxFloat32

	for axis := range 3 {
		boundsMin := node.aabb.Minimum.Axis(uint(axis))
		boundsMax := node.aabb.Maximum.Axis(uint(axis))

		if boundsMin == boundsMax {
			continue
		}

		scale := (boundsMax - boundsMin) / config.BVH_SPACES
		for i := range config.BVH_SPACES {
			candidatePos := boundsMin + float32(i+1)*scale
			cost := node.EvaluateSH(uint(axis), candidatePos, b.triangles)
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

	leftChild := NewBvhNode(firstPrim, primCount, b.triangles)
	b.nodes = append(b.nodes, leftChild)
	rightChild := NewBvhNode(i, primCount-leftCount, b.triangles)
	b.nodes = append(b.nodes, rightChild)

	b.Subdivide(leftChildIndex)
	b.Subdivide(rightChildIndex)
}
