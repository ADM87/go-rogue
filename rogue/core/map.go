package core

import (
	"math/rand"
	"rogue/data"
)

type IMap interface {
	IRectangle
	GetStart() (int, int)
	GetEnd() (int, int)
	GetRooms(region IRectangle) []IRoom
}

type Map struct {
	*Rectangle
	start IPoint
	end   IPoint
	rooms []IRoom
}

func NewMap(config data.IMapConfig) *Map {
	m := &Map{}

	m.generateLayout(config)
	m.CalculateArea()

	sX, sY := m.rooms[0].Center()
	m.start = NewPoint(sX, sY)

	eX, eY := m.rooms[len(m.rooms)-1].Center()
	m.end = NewPoint(eX, eY)

	return m
}

func (m *Map) GetStart() (int, int) {
	return m.start.GetX(), m.start.GetY()
}

func (m *Map) GetEnd() (int, int) {
	return m.end.GetX(), m.end.GetY()
}

func (m *Map) CalculateArea() {
	left, top, right, bottom := 0, 0, 0, 0
	for _, room := range m.rooms {
		if room.Left() < left {
			left = room.Left()
		}
		if room.Top() < top {
			top = room.Top()
		}
		if room.Right() > right {
			right = room.Right()
		}
		if room.Bottom() > bottom {
			bottom = room.Bottom()
		}
	}
	m.Rectangle = NewRectangle(left, top, right-left, bottom-top)
}

func (m *Map) GetRooms(region IRectangle) []IRoom {
	rooms := make([]IRoom, 0)
	for _, room := range m.rooms {
		if room.Overlaps(region) {
			rooms = append(rooms, room)
		}
	}
	return rooms
}

// Private //////////////////////////////////////////////////////////////////////

func (m *Map) generateLayout(config data.IMapConfig) {
	minRooms, maxRooms := config.GetMinRooms(), config.GetMaxRooms()
	if minRooms == 0 || maxRooms == 0 {
		panic("MapConfig must have MinRooms and MaxRooms set")
	}
	m.createRooms(config, nil, minRooms+rand.Intn(maxRooms-minRooms+1))
}

func (m *Map) createRooms(config data.IMapConfig, previousRoom IRoom, total int) {
	if len(m.rooms) == total {
		return
	}
	width, height := m.newRoomSize(config)
	if previousRoom == nil {
		m.rooms = append(m.rooms, NewRoom(0, 0, width, height))
		m.createRooms(config, m.rooms[0], total)
		return
	}
	directions := []int{North, East, South, West}
	for {
		if len(directions) == 0 || len(m.rooms) == total {
			return
		}

		i := rand.Intn(len(directions))
		direction := directions[i]
		directions = append(directions[:i], directions[i+1:]...)

		if neighbor := previousRoom.GetNeighbor(direction); neighbor != nil {
			m.createRooms(config, neighbor, total)
			continue
		}

		x, y := m.newRoomPosition(previousRoom, direction, width, height)

		room := NewRoom(x, y, width, height)
		if m.overlapsOthers(room) {
			continue
		}

		previousRoom.SetNeighbor(direction, room)
		room.SetNeighbor((direction+2)%4, previousRoom)

		m.rooms = append(m.rooms, room)
		m.createRooms(config, room, total)
	}
}

func (m *Map) newRoomSize(config data.IMapConfig) (int, int) {
	minWidth, maxWidth := config.GetMinWidth(), config.GetMaxWidth()
	if minWidth == 0 || maxWidth == 0 {
		panic("MapConfig must have MinRoomWidth and MaxRoomWidth set")
	}
	minHeight, maxHeight := config.GetMinHeight(), config.GetMaxHeight()
	if minHeight == 0 || maxHeight == 0 {
		panic("MapConfig must have MinRoomHeight and MaxRoomHeight set")
	}
	width := minWidth + rand.Intn(maxWidth-minWidth+1)
	if width&1 == 0 {
		width++
	}
	height := minHeight + rand.Intn(maxHeight-minHeight+1)
	if height&1 == 0 {
		height++
	}
	return width, height
}

func (m *Map) newRoomPosition(neighbor IRoom, direction, width, height int) (int, int) {
	x, y := neighbor.Center()
	switch direction {
	case North:
		x -= width >> 1
		y = neighbor.Top() - height
	case East:
		x = neighbor.Right()
		y -= height >> 1
	case South:
		x -= width >> 1
		y = neighbor.Bottom()
	case West:
		x = neighbor.Left() - width
		y -= height >> 1
	}
	return x, y
}

func (m *Map) overlapsOthers(room IRoom) bool {
	for _, other := range m.rooms {
		if room.Overlaps(other) {
			return true
		}
	}
	return false
}
