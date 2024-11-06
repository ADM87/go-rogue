package core

import "fmt"

// IQuadObject is an interface representing an object in a quadtree.
type IQuadObject interface {
	IPoint // IPoint the position of the object.
}

// IQuadNode is an interface representing a node in a quadtree.
type IQuadNode interface {
	fmt.Stringer                          // String returns a string representation of the node.
	IRectangle                            // IRectangle the bounds of the node.
	Find(IQuadObject) bool                // Find processes a full search for the object in the tree.
	Insert(IQuadObject) bool              // Insert inserts the object into the tree.
	Remove(IQuadObject) bool              // Remove removes the object from the tree.
	Query(IRectangle, bool) []IQuadObject // Query returns a list of objects in the region.
}

// QuadNode is a struct representing a node in a quadtree.
type QuadNode struct {
	*Rectangle                    // Rectangle the bounds of the node.
	objects         []IQuadObject // objects in the node.
	branches        []*QuadNode   // branches of the node.
	parent          IQuadNode     // parent of the node.
	depth, capacity int           // depth and capacity of the node.
}

// NewQuadBranch returns a new quadtree branch.
func NewQuadBranch(parent IQuadNode, x, y, width, height, depth, capacity int) *QuadNode {
	return &QuadNode{NewRectangle(x, y, width, height), nil, nil, parent, depth, capacity}
}

// NewQuadTree returns a new quadtree.
func NewQuadTree(x, y, width, height, depth, capacity int) *QuadNode {
	return NewQuadBranch(nil, x, y, width, height, depth, capacity)
}

// String returns a string representation of the node.
func (n *QuadNode) String() string {
	rect := n.Rectangle.String()
	objs := len(n.objects)
	brch := len(n.branches)
	return fmt.Sprintf("{Rectangle: %s, Objects: %v, Branches: %v}", rect, objs, brch)
}

// Find processes a full search for the object in the tree.
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
	return n.branches[0].Find(obj) || n.branches[1].Find(obj) ||
		n.branches[2].Find(obj) || n.branches[3].Find(obj)
}

// Insert inserts the object into the tree.
func (n *QuadNode) Insert(obj IQuadObject) bool {
	if n.Find(obj) {
		return false
	}
	return n.internalInsert(obj)
}

// Remove removes the object from the tree.
func (n *QuadNode) Remove(obj IQuadObject) bool {
	if n.Find(obj) {
		return false
	}
	return false
}

// Query returns a list of objects in the region.
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
	return n.branches[0].internalInsert(obj) || n.branches[1].internalInsert(obj) ||
		n.branches[2].internalInsert(obj) || n.branches[3].internalInsert(obj)
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
