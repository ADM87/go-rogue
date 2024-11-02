package engine

type IQuadObject interface {
	IRectangle
}

type IQuadNode interface {
	IRectangle

	Insert(IQuadObject) bool
	Remove(IQuadObject) bool
	Contains(IQuadObject) bool

	IsLeaf() bool
	IsEmpty() bool

	Query(IRectangle) []IQuadObject
}

type QuadNode struct {
	Rectangle

	_parent   *QuadNode
	_children []*QuadNode
	_objects  []IQuadObject

	_depth, _capacity int
}

func NewQuadNode(x, y, width, height, depth, capacity int, parent *QuadNode) *QuadNode {
	return &QuadNode{
		Rectangle: *NewRectangle(x, y, width, height),

		_parent:   parent,
		_depth:    depth,
		_capacity: capacity,
	}
}

func NewQuadTree(x, y, width, height, depth, capacity int) *QuadNode {
	return NewQuadNode(x, y, width, height, depth, capacity, nil)
}

func (n *QuadNode) Insert(obj IQuadObject) bool {
	if n.Contains(obj) {
		return false
	}
	if n._depth == 0 || len(n._objects) < n._capacity || !n.subdivide() {
		n._objects = append(n._objects, obj)
		return true
	}
	for _, c := range n._children {
		if c.Insert(obj) {
			return true
		}
	}
	return false
}

func (n *QuadNode) Remove(obj IQuadObject) bool {
	if !n.Contains(obj) {
		return false
	}
	if n.IsLeaf() {
		for i, o := range n._objects {
			if o == obj {
				n._objects = append(n._objects[:i], n._objects[i+1:]...)
				if n.IsEmpty() {
					n._objects = nil
				}
				return true
			}
		}
	}
	removed := false
	for _, c := range n._children {
		if c.Remove(obj) {
			removed = true
			break
		}
	}
	if removed {
		n.mergeIfEmpty()
	}
	return removed
}

func (n *QuadNode) Query(region IRectangle) []IQuadObject {
	if !n.Intersects(region) {
		return nil
	}
	var objects []IQuadObject
	if n.IsLeaf() {
		for _, o := range n._objects {
			if region.Intersects(o) {
				objects = append(objects, o)
			}
		}
	} else {
		for _, c := range n._children {
			objects = append(objects, c.Query(region)...)
		}
	}
	return objects
}

func (n *QuadNode) IsLeaf() bool {
	return len(n._children) == 0
}

func (n *QuadNode) IsEmpty() bool {
	return len(n._objects) == 0
}

func (n *QuadNode) Contains(obj IQuadObject) bool {
	if !n.Intersects(obj) {
		return false
	}
	if n.IsLeaf() {
		for _, o := range n._objects {
			if o == obj {
				return true
			}
		}
	} else {
		for _, c := range n._children {
			if c.Contains(obj) {
				return true
			}
		}
	}
	return false
}

// Private

func (n *QuadNode) subdivide() bool {
	if n._depth == 0 {
		return false
	}
	hw, hh := n.GetWidth()>>1, n.GetHeight()>>1
	rw, rh := n.GetWidth()&1, n.GetHeight()&1
	n._children = []*QuadNode{
		NewQuadNode(n.GetX(), n.GetY(), hw, hh, n._depth-1, n._capacity, n),
		NewQuadNode(n.GetX()+hw, n.GetY(), hw+rw, hh, n._depth-1, n._capacity, n),
		NewQuadNode(n.GetX(), n.GetY()+hh, hw, hh+rh, n._depth-1, n._capacity, n),
		NewQuadNode(n.GetX()+hw, n.GetY()+hh, hw+rw, hh+rh, n._depth-1, n._capacity, n),
	}
	if n.IsEmpty() {
		return true
	}
	for _, o := range n._objects {
		for _, c := range n._children {
			if c.Insert(o) {
				break
			}
		}
	}
	n._objects = nil
	return true
}

func (n *QuadNode) mergeIfEmpty() {
	if n.IsEmpty() && n.allChildrenEmpty() {
		n._children = nil
		if n._parent != nil {
			n._parent.mergeIfEmpty()
		}
	}
}

func (n *QuadNode) allChildrenEmpty() bool {
	if n.IsLeaf() {
		return false
	}
	for _, c := range n._children {
		if !c.IsEmpty() {
			return false
		}
	}
	return true
}
