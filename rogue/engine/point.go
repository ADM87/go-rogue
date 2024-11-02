package engine

import "fmt"

// IPoint defines the interface for a point with getter, setter, and copy methods.
type IPoint interface {
	GetXY() (int, int) // Get the X and Y coordinates.
	GetX() int         // Get the X coordinate.
	GetY() int         // Get the Y coordinate.
	SetXY(int, int)    // Set the X and Y coordinates.
	SetX(int)          // Set the X coordinate.
	SetY(int)          // Set the Y coordinate.
	Copy() IPoint      // Create a copy of the point.
}

// Point represents a point with X and Y coordinates.
type Point struct {
	x, y int
}

// NewPoint creates a new Point with the specified coordinates.
func NewPoint(x, y int) *Point {
	return &Point{x, y}
}

// GetXY returns the X and Y coordinates of the point.
func (p *Point) GetXY() (int, int) {
	return p.x, p.y
}

// GetX returns the X coordinate of the point.
func (p *Point) GetX() int {
	return p.x
}

// GetY returns the Y coordinate of the point.
func (p *Point) GetY() int {
	return p.y
}

// SetXY sets the X and Y coordinates of the point.
func (p *Point) SetXY(x, y int) {
	p.x = x
	p.y = y
}

// SetX sets the X coordinate of the point.
func (p *Point) SetX(x int) {
	p.x = x
}

// SetY sets the Y coordinate of the point.
func (p *Point) SetY(y int) {
	p.y = y
}

// Copy creates a new Point with the same coordinates.
func (p *Point) Copy() IPoint {
	return NewPoint(p.x, p.y)
}

// String returns a string representation of the point.
func (p *Point) String() string {
	return fmt.Sprintf("Point{%d, %d}", p.x, p.y)
}
