package engine

import "fmt"

// ICamera defines the interface for a camera with movement and viewport functionalities.
type ICamera interface {
	Move(dx, dy int)      // Move the camera by dx and dy units.
	Goto(x, y int)        // Move the camera to the specified x, y coordinates.
	Viewport() IRectangle // Get the current viewport of the camera.
}

// Camera struct embeds a Rectangle to represent the camera's viewable area.
type Camera struct {
	*Rectangle
}

// NewCamera creates a new Camera instance with specified position and size.
func NewCamera(x, y, w, h int) *Camera {
	return &Camera{
		Rectangle: NewRectangle(x, y, w, h),
	}
}

// Move shifts the camera's center by the specified delta values.
func (c *Camera) Move(dx, dy int) {
	center := c.Rectangle.GetCenter()
	c.Rectangle.SetCenter(center.GetX()+dx, center.GetY()+dy)
}

// Goto sets the camera's center to the specified coordinates.
func (c *Camera) Goto(x, y int) {
	c.Rectangle.SetCenter(x, y)
}

// Viewport returns a copy of the camera's current viewable area.
func (c *Camera) Viewport() IRectangle {
	return c.Copy()
}

// String provides a string representation of the Rectangle.
func (r *Rectangle) String() string {
	center := r.GetCenter()
	size := r.GetSize()
	return fmt.Sprintf("Rectangle{Center: %s, Size: %s}", center, size)
}
