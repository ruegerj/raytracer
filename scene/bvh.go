package scene

import (
	"fmt"
	"log"
	"math"

	"github.com/ruegerj/raytracing/common"
	"github.com/ruegerj/raytracing/common/optional"
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

	log.Println("start building bvh...")
	bvh.Subdivide(ROOT_INDEX)
	log.Println("bvh log count ", len(bvh.nodes))

	return bvh
}

func (b *Bvh) Intersects(ray primitive.Ray) optional.Optional[Hit] {
	node := b.nodes[ROOT_INDEX]
	stack := [64]BvhNode{node}
	stackPointer := 0

	var nearestDist float32 = float32(math.Inf(1))
	nearestTriangle := optional.None[Triangle]()

	for {
		if node.IsLeaf() {
			for _, tri := range node.GetOwnTriangles(b.triangles) {
				potentialHit := tri.Hits(ray)
				if potentialHit.IsEmpty() {
					continue
				}

				hit := potentialHit.Get()
				if hit.Distance < nearestDist {
					nearestDist = hit.Distance
					nearestTriangle = optional.Some(tri)
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

		nearDist := optional.None[float32]()
		farDist := optional.None[float32]()
		var nearChild, farChild BvhNode

		if dist1 > dist2 {
			nearDist = optional.Some(dist2)
			nearChild = child2
			farDist = optional.Some(dist1)
			farChild = child1
		} else {
			nearDist = optional.Some(dist1)
			nearChild = child1
			farDist = optional.Some(dist2)
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
		return optional.None[Hit]()
	}

	hit := nearestTriangle.Get().CreateHitFor(ray, nearestDist)
	return optional.Some(hit)
}

func (b *Bvh) Subdivide(nodeIndex uint) {
	node := &b.nodes[nodeIndex]

	var bestAxis uint = 3
	var bestPos float32 = 0.0
	var bestCost float32 = float32(math.Inf(1))

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
		fmt.Println("left full")
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
