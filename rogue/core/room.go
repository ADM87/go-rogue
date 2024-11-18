package core

import "math"

type IRoom interface {
	IRectangle
	Visit()
	IsWall(int, int) bool
	HasBeenVisited() bool
	GetNeighbor(int) IRoom
	SetNeighbor(int, IRoom)
}

type Room struct {
	*Rectangle
	neighbors  []IRoom
	wasVisited bool
}

func NewRoom(x, y, width, height int) *Room {
	return &Room{
		Rectangle: NewRectangle(x, y, width, height),
		neighbors: make([]IRoom, 4),
	}
}

func (r *Room) IsWall(x, y int) bool {
	if !r.Contains(x, y) {
		return false
	}
	isWall := x == r.Left() || x == r.Right()-1 || y == r.Top() || y == r.Bottom()-1
	return isWall && !r.isDoor(x, y)
}

func (r *Room) GetNeighbor(index int) IRoom {
	return r.neighbors[index]
}

func (r *Room) SetNeighbor(index int, neighbor IRoom) {
	r.neighbors[index] = neighbor
}

func (r *Room) HasBeenVisited() bool {
	return r.wasVisited
}

func (r *Room) Visit() {
	r.wasVisited = true
}

// Private /////////////////////////////////////////////////////////////////////

func (r *Room) isDoor(x, y int) bool {
	if !r.Contains(x, y) {
		return false
	}
	return y == r.Top() && r.checkConnection(North, x, y) ||
		x == r.Right()-1 && r.checkConnection(East, x, y) ||
		y == r.Bottom()-1 && r.checkConnection(South, x, y) ||
		x == r.Left() && r.checkConnection(West, x, y)
}

func (r *Room) checkConnection(direction, x, y int) bool {
	var neighbor IRoom
	if neighbor = r.GetNeighbor(direction); neighbor == nil {
		return false
	}

	cx, cy := r.Center()

	switch direction {
	case North, South:
		w := int(math.Min(float64(r.GetWidth()), float64(neighbor.GetWidth()))) >> 1
		return x >= (cx-w)+2 && x <= (cx+w)-2

	case East, West:
		h := int(math.Min(float64(r.GetHeight()), float64(neighbor.GetHeight()))) >> 1
		return y >= (cy-h)+2 && y <= (cy+h)-2

	default:
		return false
	}
}
