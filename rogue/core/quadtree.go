package core

import "fmt"

// IQuadNode is an interface representing a node in a quadtree.
type IQuadNode interface {
	fmt.Stringer                      // String returns a string representation of the node.
	IRectangle                        // IRectangle the bounds of the node.
	Find(IEntity) (IQuadNode, bool)   // Find processes a full search for the object in the tree.
	Insert(IEntity) bool              // Insert inserts the object into the tree.
	Remove(IEntity) bool              // Remove removes the object from the tree.
	Move(IEntity, int, int) bool      // Move moves the object in the tree.
	Query(IRectangle, bool) []IEntity // Query returns a list of objects in the region.
	TotalNodes() int                  // TotalNodes returns the total number of nodes in the tree.
	TotalObjects() int                // TotalObjects returns the total number of objects in the tree.

	TryToMerge()

	// Debug
	IsBorder(int, int) bool // IsBorder returns true if the entity is on the border of the node.
}

// QuadNode is a struct representing a node in a quadtree.
type QuadNode struct {
	*Rectangle                  // Rectangle the bounds of the node.
	objects         []IEntity   // objects in the node.
	branches        []*QuadNode // branches of the node.
	parent          *QuadNode   // parent of the node.
	depth, capacity int         // depth and capacity of the node.
}

// NewQuadBranch returns a new quadtree branch.
func NewQuadBranch(parent *QuadNode, x, y, width, height, depth, capacity int) *QuadNode {
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
func (n *QuadNode) Find(obj IEntity) (IQuadNode, bool) {
	if len(n.branches) == 4 {
		for _, branch := range n.branches {
			if node, ok := branch.Find(obj); ok {
				return node, true
			}
		}
	}
	if n.hasObject(obj) {
		return n, true
	}
	return nil, false
}

// Insert inserts the object into the tree.
func (n *QuadNode) Insert(obj IEntity) bool {
	if !n.Contains(obj.GetX(), obj.GetY()) || n.hasObject(obj) {
		return false
	}
	if len(n.branches) == 0 {
		if n.depth == 0 || len(n.objects) < n.capacity || !n.subdivide() {
			n.objects = append(n.objects, obj)
			return true
		}
	}
	for _, branch := range n.branches {
		if branch.Insert(obj) {
			return true
		}
	}
	return false
}

// Remove removes the object from the tree.
func (n *QuadNode) Remove(obj IEntity) bool {
	if !n.Contains(obj.GetX(), obj.GetY()) {
		return false
	}
	if len(n.branches) == 4 {
		for _, branch := range n.branches {
			if branch.Remove(obj) {
				branch.TryToMerge()
				return true
			}
		}
	}
	for i, o := range n.objects {
		if o == obj {
			n.objects = append(n.objects[:i], n.objects[i+1:]...)
			return true
		}
	}
	return false
}

// Move moves the object in the tree.
func (n *QuadNode) Move(obj IEntity, x, y int) bool {
	if !n.Remove(obj) {
		return false
	}
	obj.SetXY(x, y)
	return n.Insert(obj)
}

// Query returns a list of objects in the region.
func (n *QuadNode) Query(rect IRectangle, cull bool) []IEntity {
	if !n.Overlaps(rect) {
		return []IEntity{}
	}
	var objects []IEntity
	if len(n.branches) == 4 {
		objects = append(objects, n.branches[0].Query(rect, cull)...)
		objects = append(objects, n.branches[1].Query(rect, cull)...)
		objects = append(objects, n.branches[2].Query(rect, cull)...)
		objects = append(objects, n.branches[3].Query(rect, cull)...)
		return objects
	}
	for _, obj := range n.objects {
		if cull && !rect.Contains(obj.GetX(), obj.GetY()) {
			continue
		}
		objects = append(objects, obj)
	}
	return objects
}

// TotalNodes returns the total number of nodes in the tree.
func (n *QuadNode) TotalNodes() int {
	total := 1
	if len(n.branches) == 4 {
		total += n.branches[0].TotalNodes() + n.branches[1].TotalNodes() +
			n.branches[2].TotalNodes() + n.branches[3].TotalNodes()
	}
	return total
}

// TotalObjects returns the total number of objects in the tree.
func (n *QuadNode) TotalObjects() int {
	total := len(n.objects)
	if len(n.branches) == 4 {
		total += n.branches[0].TotalObjects() + n.branches[1].TotalObjects() +
			n.branches[2].TotalObjects() + n.branches[3].TotalObjects()
	}
	return total
}

// IsBorder returns true if the entity is on the border of the node.
func (n *QuadNode) IsBorder(x, y int) bool {
	if len(n.branches) == 4 {
		for _, branch := range n.branches {
			if branch.IsBorder(x, y) {
				return true
			}
		}
	}
	if x == n.GetX() || x == n.Right()-1 {
		return y >= n.GetY() && y < n.Bottom()
	}
	if y == n.GetY() || y == n.Bottom()-1 {
		return x >= n.GetX() && x < n.Right()
	}
	return false
}

func (n *QuadNode) TryToMerge() {
	if len(n.branches) == 0 || n.TotalObjects() > n.capacity {
		return
	}

	objects := n.Query(n.Rectangle, false)
	n.branches = nil

	for _, obj := range objects {
		n.Insert(obj)
	}
}

// Private /////////////////////////////////////////////////////////////////////

func (n *QuadNode) hasObject(obj IEntity) bool {
	for _, o := range n.objects {
		if o == obj {
			return true
		}
	}
	return false
}

func (n *QuadNode) subdivide() bool {
	hw, hh := n.GetWidth()>>1, n.GetHeight()>>1
	if hw < 4 || hh < 4 {
		n.depth = 0
		return false
	}
	ow, oh := n.GetWidth()&1, n.GetHeight()&1
	n.branches = []*QuadNode{
		NewQuadBranch(n, n.GetX(), n.GetY(), hw, hh, n.depth-1, n.capacity),
		NewQuadBranch(n, n.GetX()+hw, n.GetY(), hw+ow, hh, n.depth-1, n.capacity),
		NewQuadBranch(n, n.GetX(), n.GetY()+hh, hw, hh+oh, n.depth-1, n.capacity),
		NewQuadBranch(n, n.GetX()+hw, n.GetY()+hh, hw+ow, hh+oh, n.depth-1, n.capacity),
	}
	n.distributeToBranches()
	return true
}

func (n *QuadNode) distributeToBranches() {
	for _, obj := range n.objects {
		for _, branch := range n.branches {
			if branch.Insert(obj) {
				break
			}
		}
	}
	n.objects = nil
}
