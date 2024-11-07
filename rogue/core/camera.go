package core

import "fmt"

// ICamera is an interface representing a viewable area of a map.
type ICamera interface {
	fmt.Stringer              // String returns a string representation of the camera.
	IEntity                   // IEntity the position of the camera
	ClampToBounds(IRectangle) // ClampToBounds clamps the camera to the bounds of the map.
	Viewport() IRectangle     // Viewport returns the viewable area of the camera centered around the camera's position.
}

// Camera is a struct representing a viewable area of a map.
type Camera struct {
	*Entity        // Entity the position of the camera.
	size    IPoint // size of the camera.
}

// NewCamera returns a new camera.
func NewCamera(x, y, width, height int) *Camera {
	return &Camera{NewEntity(x, y), NewPoint(width, height)}
}

// String returns a string representation of the camera.
func (c *Camera) String() string {
	return fmt.Sprintf("{Position: %s, size: %s, Viewport: %s}", c.Entity.String(), c.size.String(), c.Viewport().String())
}

// ClampToBounds clamps the camera to the bounds of the map.
func (c *Camera) ClampToBounds(bounds IRectangle) {
	x, y := c.GetXY()
	width, height := c.size.GetXY()
	bx, by := bounds.GetXY()
	bw, bh := bounds.GetSize()
	if x-width/2 < bx {
		c.SetX(bx + width/2)
	}
	if x+width/2 > bx+bw {
		c.SetX(bx + bw - width/2)
	}
	if y-height/2 < by {
		c.SetY(by + height/2)
	}
	if y+height/2 > by+bh {
		c.SetY(by + bh - height/2)
	}
}

// Viewport returns the viewable area of the camera centered around the camera's position.
func (c *Camera) Viewport() IRectangle {
	x, y := c.GetXY()
	width, height := c.size.GetXY()
	return NewRectangle(x-width/2, y-height/2, width, height)
}
