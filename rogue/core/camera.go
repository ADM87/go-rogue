package core

import (
	"fmt"
)

// ICamera is an interface representing a viewable area of a map.
type ICamera interface {
	fmt.Stringer              // String returns a string representation of the camera.
	IPoint                    // IEntity the position of the camera
	MoveBy(int, int)          // MoveBy moves the camera by the given x and y distances.
	MoveTo(int, int)          // MoveTo moves the camera to the given x and y coordinates.
	ClampToBounds(IRectangle) // ClampToBounds clamps the camera to the bounds of the map.
	Viewport() IRectangle     // Viewport returns the viewable area of the camera centered around the camera's position.
}

// Camera is a struct representing a viewable area of a map.
type Camera struct {
	*Point        // Entity the position of the camera.
	size   IPoint // size of the camera.
}

// NewCamera returns a new camera.
func NewCamera(x, y, width, height int) *Camera {
	return &Camera{NewPoint(x, y), NewPoint(width, height)}
}

// String returns a string representation of the camera.
func (c *Camera) String() string {
	return fmt.Sprintf("{Position: %s, size: %s, Viewport: %s}", c.Point.String(), c.size.String(), c.Viewport().String())
}

// MoveBy moves the camera by the given x and y distances.
func (c *Camera) MoveBy(x, y int) {
	c.x += x
	c.y += y
}

// MoveTo moves the camera to the given x and y coordinates.
func (c *Camera) MoveTo(x, y int) {
	c.x = x
	c.y = y
}

// ClampToBounds clamps the camera to the bounds of the map.
func (c *Camera) ClampToBounds(bounds IRectangle) {
	viewport := c.Viewport()
	if viewport.Left() < bounds.Left() {
		c.x += bounds.Left() - viewport.Left()
	} else if viewport.Right() > bounds.Right() {
		c.x -= viewport.Right() - bounds.Right()
	}
	if viewport.Top() < bounds.Top() {
		c.y += bounds.Top() - viewport.Top()
	} else if viewport.Bottom() > bounds.Bottom() {
		c.y -= viewport.Bottom() - bounds.Bottom()
	}
}

// Viewport returns the viewable area of the camera centered around the camera's position.
func (c *Camera) Viewport() IRectangle {
	x, y := c.GetXY()
	width, height := c.size.GetXY()
	return NewRectangle(x-width/2, y-height/2, width, height)
}
