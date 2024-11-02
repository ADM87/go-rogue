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
	center := c.Rectangle.GetCenter()
	c.Rectangle.SetCenter(center.GetX()+dx, center.GetY()+dy)
}

func (c *Camera) Goto(x, y int) {
	c.Rectangle.SetCenter(x, y)
}

func (c *Camera) Viewport() IRectangle {
	return c.Copy()
}

func (r *Rectangle) String() string {
	center := r.GetCenter()
	size := r.GetSize()
	return fmt.Sprintf("Rectangle{Center: %s, Size: %s}", center, size)
}
