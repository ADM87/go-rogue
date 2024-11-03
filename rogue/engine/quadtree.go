package engine

type IQuadObject interface {
	IPoint
}

type IQuadNode interface {
	IRectangle

	CountNodes() int
	CountObjects() int

	Find(IQuadObject) (IQuadNode, bool)
	Insert(IQuadObject) bool

	Query(IRectangle, bool) []IQuadObject
}

type QuadNode struct {
	*Rectangle

	objects []IQuadObject
	nodes   []*QuadNode
	parent  *QuadNode

	depth, capacity int
}

func NewQuadNode(x, y, width, height, depth, capacity int, parent *QuadNode) *QuadNode {
	return &QuadNode{
		Rectangle: NewRectangle(x, y, width, height),

		parent:   parent,
		depth:    depth,
		capacity: capacity,
	}
}

func NewQuadTree(width, height, depth, capacity int) *QuadNode {
	return NewQuadNode(0, 0, width, height, depth, capacity, nil)
}

func (q *QuadNode) CountNodes() int {
	count := 1
	for _, n := range q.nodes {
		count += n.CountNodes()
	}
	return count
}

func (q *QuadNode) CountObjects() int {
	count := len(q.objects)
	for _, n := range q.nodes {
		count += n.CountObjects()
	}
	return count
}

// Searches entire tree for the object; does not check for containment
func (q *QuadNode) Find(obj IQuadObject) (IQuadNode, bool) {
	for _, o := range q.objects {
		if o == obj {
			return q, true
		}
	}
	for _, n := range q.nodes {
		if node, found := n.Find(obj); found {
			return node, true
		}
	}
	return nil, false
}

func (q *QuadNode) Insert(obj IQuadObject) bool {
	if _, found := q.Find(obj); found {
		return false
	}
	return q.internalInsert(obj)
}

func (q *QuadNode) Query(rect IRectangle, cull bool) []IQuadObject {
	if !q.Intersects(rect) {
		return nil
	}
	objects := make([]IQuadObject, 0)
	for _, n := range q.nodes {
		objects = append(objects, n.Query(rect, cull)...)
	}
	for _, obj := range q.objects {
		if !cull || rect.ContainsPoint(obj) {
			objects = append(objects, obj)
		}
	}
	return objects
}

// Private methods

func (q *QuadNode) internalInsert(obj IQuadObject) bool {
	if !q.ContainsPoint(obj) {
		return false
	}
	if q.tryInsertOnNodes(obj) {
		return true
	}
	if len(q.objects) < q.capacity || q.depth == 0 {
		q.objects = append(q.objects, obj)
		return true
	}
	if q.subdivide() {
		q.distribute()
		return q.tryInsertOnNodes(obj)
	}
	q.objects = append(q.objects, obj)
	return true
}

func (q *QuadNode) tryInsertOnNodes(obj IQuadObject) bool {
	for _, n := range q.nodes {
		if n.internalInsert(obj) {
			return true
		}
	}
	return false
}

func (q *QuadNode) subdivide() bool {
	ex := CalculateExtents(q.GetWidth(), q.GetHeight())
	hw, hh := ex[0].GetX(), ex[0].GetY()
	if hw == 0 || hh == 0 {
		q.depth = 0
		return false
	}
	rw, rh := ex[1].GetX(), ex[1].GetY()
	center := q.GetCenter()
	q.nodes = []*QuadNode{
		NewQuadNode(center.GetX()-hw, center.GetY()-hh, hw+rw, hh+rh, q.depth-1, q.capacity, q),
		NewQuadNode(center.GetX()+hw, center.GetY()-hh, hw+rw, hh+rh, q.depth-1, q.capacity, q),
		NewQuadNode(center.GetX()-hw, center.GetY()+hh, hw+rw, hh+rh, q.depth-1, q.capacity, q),
		NewQuadNode(center.GetX()+hw, center.GetY()+hh, hw+rw, hh+rh, q.depth-1, q.capacity, q),
	}
	return true
}

func (q *QuadNode) distribute() {
	for _, obj := range q.objects {
		q.tryInsertOnNodes(obj)
	}
	q.objects = nil
}
