package engine

type IRectangle interface {
	GetExtents() (IPoint, IPoint)

	GetCenterX() int
	SetCenterX(int)

	GetCenterY() int
	SetCenterY(int)

	GetCenter() IPoint
	SetCenter(int, int)

	GetSize() IPoint
	SetSize(int, int)

	GetWidth() int
	SetWidth(int)

	GetHeight() int
	SetHeight(int)

	Min() IPoint
	Max() IPoint

	Left() int
	Right() int
	Top() int
	Bottom() int

	ContainsXY(int, int) bool
	ContainsPoint(IPoint) bool

	Intersects(IRectangle) bool

	Copy() IRectangle
}

type Rectangle struct {
	_center  IPoint
	_size    IPoint
	_extents [2]IPoint
}

func NewRectangle(x, y, width, height int) *Rectangle {
	return &Rectangle{
		_center:  NewPoint(x, y),
		_size:    NewPoint(width, height),
		_extents: CalculateExtends(width, height),
	}
}

func (r *Rectangle) ContainsXY(x, y int) bool {
	return x >= r.Left() && x <= r.Right() && y >= r.Top() && y <= r.Bottom()
}

func (r *Rectangle) ContainsPoint(p IPoint) bool {
	return r.ContainsXY(p.GetX(), p.GetY())
}

func (r *Rectangle) Intersects(other IRectangle) bool {
	return r.Left() <= other.Right() && r.Right() >= other.Left() &&
		r.Top() <= other.Bottom() && r.Bottom() >= other.Top()
}

func (r *Rectangle) GetExtents() (IPoint, IPoint) {
	return r._extents[0], r._extents[1]
}

func (r *Rectangle) GetCenterX() int {
	return r._center.GetX()
}

func (r *Rectangle) GetCenterY() int {
	return r._center.GetY()
}

func (r *Rectangle) GetCenter() IPoint {
	return r._center.Copy()
}

func (r *Rectangle) GetWidth() int {
	return r._size.GetX()
}

func (r *Rectangle) GetHeight() int {
	return r._size.GetY()
}

func (r *Rectangle) GetSize() IPoint {
	return r._size.Copy()
}

func (r *Rectangle) SetCenterX(x int) {
	r._center.SetX(x)
}

func (r *Rectangle) SetCenterY(y int) {
	r._center.SetY(y)
}

func (r *Rectangle) SetCenter(x, y int) {
	r._center.SetXY(x, y)
}

func (r *Rectangle) SetWidth(width int) {
	r._size.SetX(width)
	r._extents = CalculateExtends(width, r.GetHeight())
}

func (r *Rectangle) SetHeight(height int) {
	r._size.SetY(height)
	r._extents = CalculateExtends(r.GetWidth(), height)
}

func (r *Rectangle) SetSize(width, height int) {
	r._size.SetXY(width, height)
	r._extents = CalculateExtends(width, height)
}

func (r *Rectangle) Min() IPoint {
	return r._extents[0]
}

func (r *Rectangle) Max() IPoint {
	return r._extents[1]
}

func (r *Rectangle) Left() int {
	return r._center.GetX() + r.Min().GetX()
}

func (r *Rectangle) Right() int {
	return r._center.GetX() + r.Max().GetX()
}

func (r *Rectangle) Top() int {
	return r._center.GetY() + r.Min().GetY()
}

func (r *Rectangle) Bottom() int {
	return r._center.GetY() + r.Max().GetY()
}

func (r *Rectangle) Copy() IRectangle {
	return &Rectangle{
		_center:  r._center.Copy(),
		_extents: r._extents,
		_size:    r._size.Copy(),
	}
}
