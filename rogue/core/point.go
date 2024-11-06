package core

import "fmt"

// IPoint is an interface representing a point in a map.
type IPoint interface {
	fmt.Stringer       // String returns a string representation of the point.
	ICopy              // Copy returns a copy of the point.
	IEquals            // Equals returns true if the point is equal to the other point.
	GetX() int         // GetX returns the x coordinate of the point.
	GetY() int         // GetY returns the y coordinate of the point.
	GetXY() (int, int) // GetXY returns the x and y coordinates of the point.
	SetX(int)          // SetX sets the x coordinate of the point.
	SetY(int)          // SetY sets the y coordinate of the point.
	SetXY(int, int)    // SetXY sets the x and y coordinates of the point.
}

// Point is a struct representing a point in a map.
type Point struct {
	x, y int
}

// NewPoint returns a new point.
func NewPoint(x, y int) *Point {
	return &Point{x, y}
}

// String returns a string representation of the point.
func (p *Point) String() string {
	return fmt.Sprintf("{x: %d, y: %d}", p.x, p.y)
}

// Copy returns a copy of the point.
func (p *Point) Copy() interface{} {
	return &Point{p.x, p.y}
}

// Equals returns true if the point is equal to the other point.
func (p *Point) Equals(other interface{}) bool {
	if other == nil {
		return false
	}
	otherPoint, ok := other.(*Point)
	if !ok {
		return false
	}
	return p.x == otherPoint.x && p.y == otherPoint.y
}

// GetX returns the x coordinate of the point.
func (p *Point) GetX() int {
	return p.x
}

// GetY returns the y coordinate of the point.
func (p *Point) GetY() int {
	return p.y
}

// GetXY returns the x and y coordinates of the point.
func (p *Point) GetXY() (int, int) {
	return p.x, p.y
}

// SetX sets the x coordinate of the point.
func (p *Point) SetX(x int) {
	p.x = x
}

// SetY sets the y coordinate of the point.
func (p *Point) SetY(y int) {
	p.y = y
}

// SetXY sets the x and y coordinates of the point.
func (p *Point) SetXY(x, y int) {
	p.x, p.y = x, y
}
