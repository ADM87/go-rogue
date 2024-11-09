package core

import "fmt"

// IQuadNode is an interface representing a node in a quadtree.
type IQuadNode interface {
	fmt.Stringer                      // String returns a string representation of the node.
	IRectangle                        // IRectangle the bounds of the node.
	Find(IEntity) bool                // Find processes a full search for the object in the tree.
	Insert(IEntity) bool              // Insert inserts the object into the tree.
	Remove(IEntity) bool              // Remove removes the object from the tree.
	Query(IRectangle, bool) []IEntity // Query returns a list of objects in the region.
	TotalNodes() int                  // TotalNodes returns the total number of nodes in the tree.
	TotalObjects() int                // TotalObjects returns the total number of objects in the tree.
	IsBorder(int, int) bool           // IsBorder returns true if the entity is on the border of the node.
}

// QuadNode is a struct representing a node in a quadtree.
type QuadNode struct {
	*Rectangle                  // Rectangle the bounds of the node.
	objects         []IEntity   // objects in the node.
	branches        []*QuadNode // branches of the node.
	parent          IQuadNode   // parent of the node.
	depth, capacity int         // depth and capacity of the node.
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
func (n *QuadNode) Find(obj IEntity) bool {
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
func (n *QuadNode) Insert(obj IEntity) bool {
	if n.Find(obj) {
		return false
	}
	return n.internalInsert(obj)
}

// Remove removes the object from the tree.
func (n *QuadNode) Remove(obj IEntity) bool {
	if n.Find(obj) {
		return n.internalRemove(obj)
	}
	return false
}

// Query returns a list of objects in the region.
func (n *QuadNode) Query(rect IRectangle, cull bool) []IEntity {
	return n.internalQuery(rect, cull, []IEntity{})
}

// TotalNodes returns the total number of nodes in the tree.
func (n *QuadNode) TotalNodes() int {
	if len(n.branches) == 0 {
		return 1
	}
	return 1 + n.branches[0].TotalNodes() + n.branches[1].TotalNodes() +
		n.branches[2].TotalNodes() + n.branches[3].TotalNodes()
}

func (n *QuadNode) TotalObjects() int {
	if len(n.branches) == 0 {
		return len(n.objects)
	}
	return n.branches[0].TotalObjects() + n.branches[1].TotalObjects() +
		n.branches[2].TotalObjects() + n.branches[3].TotalObjects()
}

// IsBorder returns true if the entity is on the border of the node.
func (n *QuadNode) IsBorder(x, y int) bool {
	minX, minY := n.Rectangle.Min()
	maxX, maxY := n.Rectangle.Max()
	if x == minX || x == maxX {
		return y >= minY && y <= maxY
	}

	if y == minY || y == maxY {
		return x >= minX && x <= maxX
	}
	if len(n.branches) == 4 {
		return n.branches[0].IsBorder(x, y) || n.branches[1].IsBorder(x, y) ||
			n.branches[2].IsBorder(x, y) || n.branches[3].IsBorder(x, y)
	}
	return false
}

// Private /////////////////////////////////////////////////////////////////////

func (n *QuadNode) internalInsert(obj IEntity) bool {
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

func (n *QuadNode) internalRemove(obj IEntity) bool {
	if !n.Contains(obj.GetX(), obj.GetY()) {
		return false
	}
	if len(n.branches) == 0 {
		for i, object := range n.objects {
			if object.Equals(obj) {
				n.objects = append(n.objects[:i], n.objects[i+1:]...)
				return true
			}
		}
		return false
	}
	for _, branch := range n.branches {
		if branch.internalRemove(obj) {
			branch.tryMerge()
			return true
		}
	}
	return false
}

func (n *QuadNode) tryMerge() {
	if len(n.branches) == 0 || n.TotalObjects() > n.capacity {
		return
	}
	n.objects = append(n.objects, n.Query(n, false)...)
	n.branches = nil
	fmt.Println("Merged")
}

func (n *QuadNode) internalQuery(region IRectangle, cull bool, objects []IEntity) []IEntity {
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
