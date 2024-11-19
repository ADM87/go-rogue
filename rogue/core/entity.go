package core

type EntityMovementHandler func(IEntity, int, int)

type IEntity interface {
	IRectangle
	MoveBy(int, int)                  // MoveBy moves the camera by the given x and y distances.
	MoveTo(int, int)                  // MoveTo moves the camera to the given x and y coordinates.
	OnCollisionStart(IEntity)         // OnCollision is called when the entity collides with another entity.
	OnCollisionEnd()                  // OnCollision is called when the entity stops colliding with another entity.
	IsColliding() bool                // IsColliding returns true if the entity is colliding with another entity.
	GetComponent(string) interface{}  // GetComponent returns the component with the given name.
	SetComponent(string, interface{}) // SetComponent sets the component with the given name.
}

type Entity struct {
	*Rectangle
	moveHandler EntityMovementHandler
	isColliding bool
}

func NewEntity(x, y int, movementHandler EntityMovementHandler) *Entity {
	return &Entity{NewRectangle(x, y, 1, 1), movementHandler, false}
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
func (e *Entity) OnCollisionStart(other IEntity) {
	e.isColliding = true
}

// OnCollision is called when the entity stops colliding with another entity.
func (e *Entity) OnCollisionEnd() {
	e.isColliding = false
}

// IsColliding returns true if the entity is colliding with another entity.
func (e *Entity) IsColliding() bool {
	return e.isColliding
}

// GetComponent returns the component with the given name.
func (e *Entity) GetComponent(name string) interface{} {
	return nil
}

// SetComponent sets the component with the given name.
func (e *Entity) SetComponent(name string, component interface{}) {
}
