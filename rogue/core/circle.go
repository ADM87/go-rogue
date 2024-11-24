package core

type ICircle interface {
	IPoint
	GetRadius() int
	SetRadius(int)
	Contains(int, int) bool
	OverlapsCircle(ICircle) bool
	OverlapsRectangle(IRectangle) bool
	SetYScale(float64)
	GetYScale() float64
}

type Circle struct {
	*Point
	radius int
	yScale float64
}

// NewCircle creates a new Circle instance
func NewCircle(x, y, radius int, yScale float64) *Circle {
	return &Circle{
		Point:  NewPoint(x, y),
		radius: radius,
		yScale: yScale,
	}
}

func (c *Circle) GetX() int {
	return c.Point.GetX()
}

func (c *Circle) GetY() int {
	return c.Point.GetY() * int(c.yScale)
}

// GetRadius returns the radius of the circle
func (c *Circle) GetRadius() int {
	return c.radius
}

// SetRadius sets the radius of the circle
func (c *Circle) SetRadius(radius int) {
	c.radius = radius
}

// SetYScale sets the y-axis scaling factor
func (c *Circle) SetYScale(scale float64) {
	c.yScale = scale
}

// GetYScale gets the y-axis scaling factor
func (c *Circle) GetYScale() float64 {
	return c.yScale
}

// Contains checks if a point (x, y) is within the circle
func (c *Circle) Contains(x, y int) bool {
	dx := c.GetX() - x
	dy := c.GetY() - int(float64(y)*c.yScale)
	return dx*dx+dy*dy <= c.radius*c.radius
}

// OverlapsCircle checks if another circle overlaps with this circle
func (c *Circle) OverlapsCircle(circle ICircle) bool {
	dx := c.GetX() - circle.GetX()
	dy := c.GetY() - int(float64(circle.GetY())*c.yScale)
	distance := dx*dx + dy*dy
	radiusSum := c.radius + circle.GetRadius()
	return distance <= radiusSum*radiusSum
}

// OverlapsRectangle checks if the circle overlaps with a rectangle
func (c *Circle) OverlapsRectangle(rect IRectangle) bool {
	closestX := c.GetX()
	closestY := c.GetY()

	// Find the closest x coordinate
	if c.GetX() < rect.Left() {
		closestX = rect.Left()
	} else if c.GetX() > rect.Right() {
		closestX = rect.Right()
	}

	// Find the closest y coordinate
	if c.GetY() < rect.Top() {
		closestY = rect.Top()
	} else if c.GetY() > rect.Bottom() {
		closestY = rect.Bottom()
	}

	// Apply y-axis scaling
	closestY = int(float64(closestY) * c.yScale)

	// Calculate distance between circle center and closest point
	dx := c.GetX() - closestX
	dy := c.GetY() - closestY
	return dx*dx+dy*dy <= c.radius*c.radius
}
