package engine

// IRectangle defines the interface for a rectangle with various properties and methods.
type IRectangle interface {
	GetExtents() (IPoint, IPoint) // Get the extents of the rectangle.
	GetCenterX() int              // Get the X coordinate of the center.
	SetCenterX(int)               // Set the X coordinate of the center.
	GetCenterY() int              // Get the Y coordinate of the center.
	SetCenterY(int)               // Set the Y coordinate of the center.
	GetCenter() IPoint            // Get the center point.
	SetCenter(int, int)           // Set the center point.
	GetSize() IPoint              // Get the size of the rectangle.
	SetSize(int, int)             // Set the size of the rectangle.
	GetWidth() int                // Get the width of the rectangle.
	SetWidth(int)                 // Set the width of the rectangle.
	GetHeight() int               // Get the height of the rectangle.
	SetHeight(int)                // Set the height of the rectangle.
	Min() IPoint                  // Get the minimum point of the rectangle.
	Max() IPoint                  // Get the maximum point of the rectangle.
	Left() int                    // Get the left boundary of the rectangle.
	Right() int                   // Get the right boundary of the rectangle.
	Top() int                     // Get the top boundary of the rectangle.
	Bottom() int                  // Get the bottom boundary of the rectangle.
	ContainsXY(int, int) bool     // Check if the rectangle contains a point with X and Y coordinates.
	ContainsPoint(IPoint) bool    // Check if the rectangle contains a point.
	Intersects(IRectangle) bool   // Check if the rectangle intersects with another rectangle.
	Copy() IRectangle             // Create a copy of the rectangle.
}

// Rectangle represents a rectangle with center, size, and extents.
type Rectangle struct {
	center  IPoint
	size    IPoint
	extents [2]IPoint
}

// NewRectangle creates a new Rectangle with the specified position and size.
func NewRectangle(x, y, width, height int) *Rectangle {
	return &Rectangle{
		center:  NewPoint(x, y),
		size:    NewPoint(width, height),
		extents: CalculateExtents(width, height),
	}
}

// ContainsXY checks if the rectangle contains a point with X and Y coordinates.
func (r *Rectangle) ContainsXY(x, y int) bool {
	return x >= r.Left() && x <= r.Right() && y >= r.Top() && y <= r.Bottom()
}

// ContainsPoint checks if the rectangle contains a point.
func (r *Rectangle) ContainsPoint(p IPoint) bool {
	return r.ContainsXY(p.GetX(), p.GetY())
}

// Intersects checks if the rectangle intersects with another rectangle.
func (r *Rectangle) Intersects(other IRectangle) bool {
	return r.Left() <= other.Right() && r.Right() >= other.Left() &&
		r.Top() <= other.Bottom() && r.Bottom() >= other.Top()
}

// GetExtents returns the extents of the rectangle.
func (r *Rectangle) GetExtents() (IPoint, IPoint) {
	return r.extents[0], r.extents[1]
}

// GetCenterX returns the X coordinate of the center.
func (r *Rectangle) GetCenterX() int {
	return r.center.GetX()
}

// GetCenterY returns the Y coordinate of the center.
func (r *Rectangle) GetCenterY() int {
	return r.center.GetY()
}

// GetCenter returns the center point of the rectangle.
func (r *Rectangle) GetCenter() IPoint {
	return r.center.Copy()
}

// GetWidth returns the width of the rectangle.
func (r *Rectangle) GetWidth() int {
	return r.size.GetX()
}

// GetHeight returns the height of the rectangle.
func (r *Rectangle) GetHeight() int {
	return r.size.GetY()
}

// GetSize returns the size of the rectangle.
func (r *Rectangle) GetSize() IPoint {
	return r.size.Copy()
}

// SetCenterX sets the X coordinate of the center.
func (r *Rectangle) SetCenterX(x int) {
	r.center.SetX(x)
}

// SetCenterY sets the Y coordinate of the center.
func (r *Rectangle) SetCenterY(y int) {
	r.center.SetY(y)
}

// SetCenter sets the center point of the rectangle.
func (r *Rectangle) SetCenter(x, y int) {
	r.center.SetXY(x, y)
}

// SetWidth sets the width of the rectangle.
func (r *Rectangle) SetWidth(width int) {
	r.size.SetX(width)
	r.extents = CalculateExtents(width, r.GetHeight())
}

// SetHeight sets the height of the rectangle.
func (r *Rectangle) SetHeight(height int) {
	r.size.SetY(height)
	r.extents = CalculateExtents(r.GetWidth(), height)
}

// SetSize sets the size of the rectangle.
func (r *Rectangle) SetSize(width, height int) {
	r.size.SetXY(width, height)
	r.extents = CalculateExtents(width, height)
}

// Min returns the minimum point of the rectangle.
func (r *Rectangle) Min() IPoint {
	return r.extents[0]
}

// Max returns the maximum point of the rectangle.
func (r *Rectangle) Max() IPoint {
	return r.extents[1]
}

// Left returns the left boundary of the rectangle.
func (r *Rectangle) Left() int {
	return r.center.GetX() + r.Min().GetX()
}

// Right returns the right boundary of the rectangle.
func (r *Rectangle) Right() int {
	return r.center.GetX() + r.Max().GetX()
}

// Top returns the top boundary of the rectangle.
func (r *Rectangle) Top() int {
	return r.center.GetY() + r.Min().GetY()
}

// Bottom returns the bottom boundary of the rectangle.
func (r *Rectangle) Bottom() int {
	return r.center.GetY() + r.Max().GetY()
}

// Copy creates a copy of the rectangle.
func (r *Rectangle) Copy() IRectangle {
	return &Rectangle{
		center:  r.center.Copy(),
		extents: r.extents,
		size:    r.size.Copy(),
	}
}
