package core

type IEntity interface {
	IPoint
	MoveBy(int, int) // MoveBy moves the camera by the given x and y distances.
	MoveTo(int, int) // MoveTo moves the camera to the given x and y coordinates.
}

type Entity struct {
	*Point
}

func NewEntity(x, y int) *Entity {
	return &Entity{NewPoint(x, y)}
}

func (e *Entity) Equals(other interface{}) bool {
	if other == nil {
		return false
	}
	otherEntity, ok := other.(*Entity)
	if !ok {
		return false
	}
	return e.x == otherEntity.x && e.y == otherEntity.y
}

// MoveBy moves the camera by the given x and y distances.
func (e *Entity) MoveBy(x, y int) {
	e.x += x
	e.y += y
}

// MoveTo moves the camera to the given x and y coordinates.
func (e *Entity) MoveTo(x, y int) {
	e.x = x
	e.y = y
}
