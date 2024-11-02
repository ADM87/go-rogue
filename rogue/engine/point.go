package engine

type IPoint interface {
	GetXY() (int, int)
	GetX() int
	GetY() int

	SetXY(int, int)
	SetX(int)
	SetY(int)

	Copy() IPoint
}

type Point struct {
	_x, _y int
}

func NewPoint(x, y int) *Point {
	return &Point{x, y}
}

func (p *Point) GetXY() (int, int) {
	return p._x, p._y
}

func (p *Point) GetX() int {
	return p._x
}

func (p *Point) GetY() int {
	return p._y
}

func (p *Point) SetXY(x, y int) {
	p._x = x
	p._y = y
}

func (p *Point) SetX(x int) {
	p._x = x
}

func (p *Point) SetY(y int) {
	p._y = y
}

func (p *Point) Copy() IPoint {
	return NewPoint(p._x, p._y)
}
