package engine

type IRectangle interface {
	IPoint

	GetWidth() int
	GetHeight() int
	GetSize() (int, int)

	SetWidth(int)
	SetHeight(int)
	SetSize(int, int)

	MinX() int
	MinY() int
	MinXY() (int, int)

	MaxX() int
	MaxY() int
	MaxXY() (int, int)
}

type Rectangle struct {
	Point

	_width, _height int
}

func NewRectangle(x, y, width, height int) *Rectangle {
	return &Rectangle{
		Point: Point{x, y},

		_width:  width,
		_height: height,
	}
}

func (r *Rectangle) GetWidth() int {
	return r._width
}

func (r *Rectangle) GetHeight() int {
	return r._height
}

func (r *Rectangle) GetSize() (int, int) {
	return r._width, r._height
}

func (r *Rectangle) SetWidth(width int) {
	r._width = width
}

func (r *Rectangle) SetHeight(height int) {
	r._height = height
}

func (r *Rectangle) SetSize(width, height int) {
	r._width = width
	r._height = height
}

func (r *Rectangle) MinX() int {
	return r.GetX()
}

func (r *Rectangle) MinY() int {
	return r.GetY()
}

func (r *Rectangle) MinXY() (int, int) {
	return r.GetXY()
}

func (r *Rectangle) MaxX() int {
	return r.GetX() + r.GetWidth()
}

func (r *Rectangle) MaxY() int {
	return r.GetY() + r.GetHeight()
}

func (r *Rectangle) MaxXY() (int, int) {
	return r.GetX() + r.GetWidth(), r.GetY() + r.GetHeight()
}
