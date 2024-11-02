package engine

type IQuadObject interface {
	IPoint
}

type IQuadNode interface {
	IRectangle
}

type QuadNode struct {
	*Rectangle

	objects []IQuadObject
	nodes   []IQuadNode
	parent  IQuadNode

	depth, capacity int
}

func NewQuadNode(x, y, width, height, depth, capacity int, parent IQuadNode) *QuadNode {
	return &QuadNode{
		Rectangle: NewRectangle(x, y, width, height),

		parent:   parent,
		depth:    depth,
		capacity: capacity,
	}
}

func NewQuadTree(x, y, width, height, depth, capacity int) *QuadNode {
	return NewQuadNode(x, y, width, height, depth, capacity, nil)
}
