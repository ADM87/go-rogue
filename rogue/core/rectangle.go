package core

import "fmt"

// IRectangle is an interface representing a rectangle in a map.
type IRectangle interface {
	ICopy                         // Copy returns a copy of the rectangle.
	IPoint                        // IPoint the position of the rectangle.
	Left() int                    // Left returns the left side of the rectangle.
	Right() int                   // Right returns the right side of the rectangle.
	Top() int                     // Top returns the top side of the rectangle.
	Bottom() int                  // Bottom returns the bottom side of the rectangle.
	Center() (int, int)           // Center returns the center of the rectangle.
	Min() (int, int)              // Min returns the left and top sides of the rectangle.
	Max() (int, int)              // Max returns the right and bottom sides of the rectangle.
	GetWidth() int                // GetWidth returns the width of the rectangle.
	GetHeight() int               // GetHeight returns the height of the rectangle.
	GetSize() (int, int)          // GetSize returns the width and height of the rectangle.
	SetWidth(int)                 // SetWidth sets the width of the rectangle.
	SetHeight(int)                // SetHeight sets the height of the rectangle.
	SetSize(int, int)             // SetSize sets the width and height of the rectangle.
	Contains(int, int) bool       // Contains returns true if the rectangle contains the point.
	CollidesWith(IRectangle) bool // Collides returns true if the rectangle intersects with the other rectangle.
	Overlaps(IRectangle) bool     // Overlaps returns true if the rectangle overlaps with the other rectangle.
}

// Rectangle is a struct representing a rectangle in a map.
type Rectangle struct {
	*Point            // Point the position of the rectangle.
	width, height int // width and height of the rectangle.
}

// NewRectangle returns a new rectangle.
func NewRectangle(x, y, width, height int) *Rectangle {
	return &Rectangle{NewPoint(x, y), width, height}
}

// String returns a string representation of the rectangle.
func (r *Rectangle) String() string {
	return fmt.Sprintf("{x: %d, y: %d, width: %d, height: %d}", r.x, r.y, r.width, r.height)
}

// Copy returns a copy of the rectangle.
func (r *Rectangle) Copy() interface{} {
	return &Rectangle{r.Point.Copy().(*Point), r.width, r.height}
}

// Left returns the left side of the rectangle.
func (r *Rectangle) Left() int {
	return r.x
}

// Right returns the right side of the rectangle.
func (r *Rectangle) Right() int {
	return r.x + r.width
}

// Top returns the top side of the rectangle.
func (r *Rectangle) Top() int {
	return r.y
}

// Bottom returns the bottom side of the rectangle.
func (r *Rectangle) Bottom() int {
	return r.y + r.height
}

// Center returns the center of the rectangle.
func (r *Rectangle) Center() (int, int) {
	return r.x + r.width>>1, r.y + r.height>>1
}

// Min returns the left and top sides of the rectangle.
func (r *Rectangle) Min() (int, int) {
	return r.Left(), r.Top()
}

// Max returns the right and bottom sides of the rectangle.
func (r *Rectangle) Max() (int, int) {
	return r.Right(), r.Bottom()
}

// GetWidth returns the width of the rectangle.
func (r *Rectangle) GetWidth() int {
	return r.width
}

// GetHeight returns the height of the rectangle.
func (r *Rectangle) GetHeight() int {
	return r.height
}

// GetSize returns the width and height of the rectangle.
func (r *Rectangle) GetSize() (int, int) {
	return r.width, r.height
}

// SetWidth sets the width of the rectangle.
func (r *Rectangle) SetWidth(width int) {
	r.width = width
}

// SetHeight sets the height of the rectangle.
func (r *Rectangle) SetHeight(height int) {
	r.height = height
}

// SetSize sets the width and height of the rectangle.
func (r *Rectangle) SetSize(width, height int) {
	r.width, r.height = width, height
}

// Contains returns true if the rectangle contains the point.
func (r *Rectangle) Contains(x, y int) bool {
	return r.Left() <= x && x < r.Right() && r.Top() <= y && y < r.Bottom()
}

// CollidesWith returns true if the rectangle intersects with the other rectangle.
func (r *Rectangle) CollidesWith(other IRectangle) bool {
	return r.Left() <= other.Right() && other.Left() <= r.Right() && r.Top() <= other.Bottom() && other.Top() <= r.Bottom()
}

// Overlaps returns true if the rectangle overlaps with the other rectangle.
func (r *Rectangle) Overlaps(other IRectangle) bool {
	return r.Left() < other.Right() && other.Left() < r.Right() && r.Top() < other.Bottom() && other.Top() < r.Bottom()
}
