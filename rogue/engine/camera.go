package engine

import "fmt"

type ICamera interface {
	Move(dx, dy int)
	Goto(x, y int)
	Viewport() IRectangle
}

type Camera struct {
	*Rectangle
}

func NewCamera(x, y, w, h int) *Camera {
	return &Camera{
		Rectangle: NewRectangle(x, y, w, h),
	}
}

func (c *Camera) Move(dx, dy int) {
	x, y := c.Rectangle.GetCenter()
	c.Rectangle.SetCenter(x+dx, y+dy)
}

func (c *Camera) Goto(x, y int) {
	c.Rectangle.SetCenter(x, y)
}

func (c *Camera) Viewport() IRectangle {
	return c.Copy()
}

func (r *Rectangle) String() string {
	return fmt.Sprintf("Rectangle[Center: (%d, %d), Width: %d, Height: %d]", r.GetCenterX(), r.GetCenterY(), r.GetWidth(), r.GetHeight())
}
