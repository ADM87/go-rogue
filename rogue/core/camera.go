package core

import "fmt"

// ICamera is an interface representing a viewable area of a map.
type ICamera interface {
	fmt.Stringer          // String returns a string representation of the camera.
	IPoint                // IPoint the position of the camera.
	MoveBy(int, int)      // MoveBy moves the camera by the given x and y distances.
	MoveTo(int, int)      // MoveTo moves the camera to the given x and y coordinates.
	Viewport() IRectangle // Viewport returns the viewable area of the camera centered around the camera's position.
}

// Camera is a struct representing a viewable area of a map.
type Camera struct {
	*Point        // Point the position of the camera.
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

// Viewport returns the viewable area of the camera centered around the camera's position.
func (c *Camera) Viewport() IRectangle {
	x, y := c.GetXY()
	width, height := c.size.GetXY()
	return NewRectangle(x-width/2, y-height/2, width, height)
}
