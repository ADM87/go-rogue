package core

import "fmt"

type IQuadObject interface {
	IPoint
}

type IQuadNode interface {
	fmt.Stringer
	IRectangle
	Find(IQuadObject) bool
	Insert(IQuadObject) bool
	Remove(IQuadObject) bool
	Query(IRectangle, bool) []IQuadObject
}

type QuadNode struct {
	*Rectangle
	objects         []IQuadObject
	branches        []*QuadNode
	parent          IQuadNode
	depth, capacity int
}

func NewQuadBranch(parent IQuadNode, x, y, width, height, depth, capacity int) *QuadNode {
	return &QuadNode{NewRectangle(x, y, width, height), nil, nil, parent, depth, capacity}
}

func NewQuadTree(x, y, width, height, depth, capacity int) *QuadNode {
	return NewQuadBranch(nil, x, y, width, height, depth, capacity)
}

func (n *QuadNode) String() string {
	rect := n.Rectangle.String()
	objs := len(n.objects)
	brch := len(n.branches)
	return fmt.Sprintf("{Rectangle: %s, Objects: %v, Branches: %v}", rect, objs, brch)
}

// Find Processes a full search for the object in the tree.
func (n *QuadNode) Find(obj IQuadObject) bool {
	if !n.Contains(obj.GetX(), obj.GetY()) {
		return false
	}
	if len(n.branches) == 0 {
		for _, o := range n.objects {
			if o.Equals(obj) {
				return true
			}
		}
		return false
	}
	if n.branches[0].Find(obj) || n.branches[1].Find(obj) ||
		n.branches[2].Find(obj) || n.branches[3].Find(obj) {
		return true
	}
	return false
}

func (n *QuadNode) Insert(obj IQuadObject) bool {
	if n.Find(obj) {
		return false
	}
	return n.internalInsert(obj)
}

func (n *QuadNode) Remove(obj IQuadObject) bool {
	if n.Find(obj) {
		return false
	}
	return false
}

func (n *QuadNode) Query(rect IRectangle, cull bool) []IQuadObject {
	return n.internalQuery(rect, cull, []IQuadObject{})
}

// Private /////////////////////////////////////////////////////////////////////

func (n *QuadNode) internalInsert(obj IQuadObject) bool {
	if !n.Contains(obj.GetX(), obj.GetY()) {
		return false
	}
	if len(n.branches) == 0 {
		if n.depth == 0 || len(n.objects) < n.capacity || !n.subdivide() {
			n.objects = append(n.objects, obj)
			return true
		}
	}
	if n.branches[0].internalInsert(obj) || n.branches[1].internalInsert(obj) ||
		n.branches[2].internalInsert(obj) || n.branches[3].internalInsert(obj) {
		return true
	}
	return false
}

func (n *QuadNode) internalQuery(region IRectangle, cull bool, objects []IQuadObject) []IQuadObject {
	if !n.Overlaps(region) {
		return objects
	}
	if len(n.branches) == 0 {
		if !cull {
			return append(objects, n.objects...)
		}
		for _, o := range n.objects {
			if region.Contains(o.GetX(), o.GetY()) {
				objects = append(objects, o)
			}
		}
		return objects
	}
	objects = n.branches[0].internalQuery(region, cull, objects)
	objects = n.branches[1].internalQuery(region, cull, objects)
	objects = n.branches[2].internalQuery(region, cull, objects)
	objects = n.branches[3].internalQuery(region, cull, objects)
	return objects
}

func (n *QuadNode) subdivide() bool {
	hw, hh := n.GetWidth()>>1, n.GetHeight()>>1
	if hw < 2 || hh < 2 {
		n.depth = 0
		return false
	}
	ow, oh := n.GetWidth()&1, n.GetHeight()&1
	n.branches = []*QuadNode{
		NewQuadBranch(n, n.Left(), n.Top(), hw, hh, n.depth-1, n.capacity),
		NewQuadBranch(n, n.Left()+hw, n.Top(), hw+ow, hh, n.depth-1, n.capacity),
		NewQuadBranch(n, n.Left(), n.Top()+hh, hw, hh+oh, n.depth-1, n.capacity),
		NewQuadBranch(n, n.Left()+hw, n.Top()+hh, hw+ow, hh+oh, n.depth-1, n.capacity),
	}
	n.redistribute()
	return true
}

func (n *QuadNode) redistribute() {
	for _, o := range n.objects {
		if n.branches[0].internalInsert(o) || n.branches[1].internalInsert(o) ||
			n.branches[2].internalInsert(o) || n.branches[3].internalInsert(o) {
			continue
		}
	}
	n.objects = nil
}
