package core

type EntityMovementHandler func(IEntity, int, int)

type IEntity interface {
	IRectangle
	MoveBy(int, int)     // MoveBy moves the camera by the given x and y distances.
	MoveTo(int, int)     // MoveTo moves the camera to the given x and y coordinates.
	OnCollision(IEntity) // OnCollision is called when the entity collides with another entity.
}

type Entity struct {
	*Rectangle
	moveHandler EntityMovementHandler
}

func NewEntity(x, y int, movementHandler EntityMovementHandler) *Entity {
	return &Entity{NewRectangle(x, y, 1, 1), movementHandler}
}

// MoveBy moves the camera by the given x and y distances.
func (e *Entity) MoveBy(x, y int) {
	e.moveHandler(e, e.x+x, e.y+y)
}

// MoveTo moves the camera to the given x and y coordinates.
func (e *Entity) MoveTo(x, y int) {
	e.moveHandler(e, x, y)
}

// OnCollision is called when the entity collides with another entity.
func (e *Entity) OnCollision(other IEntity) {
}
