package core

import (
	"math/rand"
	"rogue/data"
)

type IMap interface {
	IRectangle
	GetStart() (int, int)
	GetEnd() (int, int)
}

type Map struct {
	*Rectangle
	start IPoint
	end   IPoint
	rooms []IRectangle
}

func NewMap(config data.IMapConfig) *Map {
	m := &Map{}

	minRooms, maxRooms := config.GetMinRooms(), config.GetMaxRooms()
	if minRooms == 0 || maxRooms == 0 {
		panic("MapConfig must have MinRooms and MaxRooms set")
	}

	totalRooms := minRooms + rand.Intn(maxRooms-minRooms+1)
	for i := 0; i < totalRooms; i++ {
		room := m.newRoom(config)
		m.rooms = append(m.rooms, room)
	}

	sX, sY := m.rooms[0].Center()
	m.start = NewPoint(sX, sY)

	eX, eY := m.rooms[len(m.rooms)-1].Center()
	m.end = NewPoint(eX, eY)

	m.CalculateArea()
	return m
}

func (m *Map) GetStart() (int, int) {
	return m.start.GetX(), m.start.GetY()
}

func (m *Map) GetEnd() (int, int) {
	return m.end.GetX(), m.end.GetY()
}

// Private //////////////////////////////////////////////////////////////////////

func (m *Map) newRoom(config data.IMapConfig) IRectangle {
	minWidth, maxWidth := config.GetMinWidth(), config.GetMaxWidth()
	if minWidth == 0 || maxWidth == 0 {
		panic("MapConfig must have MinRoomWidth and MaxRoomWidth set")
	}
	minHeight, maxHeight := config.GetMinHeight(), config.GetMaxHeight()
	if minHeight == 0 || maxHeight == 0 {
		panic("MapConfig must have MinRoomHeight and MaxRoomHeight set")
	}

	width := minWidth + rand.Intn(maxWidth-minWidth+1)
	height := minHeight + rand.Intn(maxHeight-minHeight+1)

	return NewRectangle(0, 0, width, height)
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
