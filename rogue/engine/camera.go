package engine

import "fmt"

type ICamera interface {
	IPoint

	GetViewport() IRectangle

	MoveX(int)
	MoveY(int)
	MoveXY(int, int)
}

type Camera struct {
	Point

	_width, _height int
}

func NewCamera(x, y, width, height int) *Camera {
	return &Camera{
		Point: Point{x, y},

		_width:  width,
		_height: height,
	}
}

func (c *Camera) GetViewport() IRectangle {
	hw, hh := c._width>>1, c._height>>1
	rw, rh := c._width&1, c._height&1
	return NewRectangle(c.GetX()-hw+rw, c.GetY()-hh+rh, c._width, c._height)
}

func (c *Camera) MoveX(x int) {
	c.Point.SetX(c.GetX() + x)
}

func (c *Camera) MoveY(y int) {
	c.Point.SetY(c.GetY() + y)
}

func (c *Camera) MoveXY(x, y int) {
	c.MoveX(x)
	c.MoveY(y)
}

func (c *Camera) DrawDebug() string {
	vp := c.GetViewport()

	output := ""
	output += fmt.Sprintf("Camera: %v\n", c)
	output += fmt.Sprintf("Viewport: %v\n", vp)
	for y := vp.MinY(); y < vp.MaxY(); y++ {
		for x := vp.MinX(); x < vp.MaxX(); x++ {
			if x == 0 && y == 0 {
				output += "@"
				continue
			}
			if x == c.GetX() && y == c.GetY() {
				output += "X"
				continue
			}
			output += "."
		}
		output += "\n"
	}
	return output
}
