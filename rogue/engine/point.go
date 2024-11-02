package engine

type IPoint interface {
	GetX() int
	GetY() int
	GetXY() (int, int)

	SetX(int)
	SetY(int)
	SetXY(int, int)
}

type Point struct {
	_x, _y int
}

func NewPoint(x, y int) *Point {
	return &Point{x, y}
}

func (p *Point) GetX() int {
	return p._x
}

func (p *Point) GetY() int {
	return p._y
}

func (p *Point) GetXY() (int, int) {
	return p._x, p._y
}

func (p *Point) SetX(x int) {
	p._x = x
}

func (p *Point) SetY(y int) {
	p._y = y
}

func (p *Point) SetXY(x, y int) {
	p._x = x
	p._y = y
}
