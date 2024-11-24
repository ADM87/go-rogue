package core

import (
	"math/rand"
	"rogue/data"
)

type IMap interface {
	IRectangle
	GetStart() (int, int)
	GetEnd() (int, int)
	GetRooms(IRectangle) []IRoom
	SetActiveRegion(IRectangle)
	Render(int, int, int, int, bool) int
}

type Map struct {
	*Rectangle
	start       IPoint
	end         IPoint
	rooms       []IRoom
	activeRooms []IRoom
}

func NewMap(config data.IMapConfig) *Map {
	m := &Map{}

	m.generateLayout(config)
	m.calculateDimensions()

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

func (m *Map) GetRooms(region IRectangle) []IRoom {
	rooms := make([]IRoom, 0)
	for _, room := range m.rooms {
		if room.Overlaps(region) {
			rooms = append(rooms, room)
		}
	}
	return rooms
}

func (m *Map) SetActiveRegion(region IRectangle) {
	m.activeRooms = m.activeRooms[:0]
	for _, room := range m.rooms {
		if room.Overlaps(region) {
			m.activeRooms = append(m.activeRooms, room)
		}
	}
}

func (m *Map) Render(x, y, focalX, focalY int, los bool) int {
	if !m.Contains(x, y) {
		return data.OutOfBounds
	}
	visible := true
	if los {
		Raycast(x, y, focalX, focalY, 1, func(x, y int) bool {
			for _, room := range m.activeRooms {
				if !room.Contains(x, y) {
					continue
				}
				if room.IsWall(x, y) {
					visible = false
					return true
				}
			}
			visible = true
			return false
		})
	}
	result := data.OutOfBounds
	for _, room := range m.activeRooms {
		if room.Contains(x, y) {
			if room.IsWall(x, y) {
				result = data.Wall
				break
			}
			if !visible {
				result = data.NotVisible
				break
			}
			result = data.Floor
			break
		}
	}
	return result
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
	directions := []int{data.North, data.East, data.South, data.West}
	for {
		if len(m.rooms) == total {
			return
		}
		dir := directions[rand.Intn(len(directions))]
		if neighbor := previousRoom.GetNeighbor(dir); neighbor != nil {
			m.createRooms(config, neighbor, total)
			continue
		}
		x, y := m.newRoomPosition(previousRoom, dir, width, height)
		room := NewRoom(x, y, width, height)
		if m.overlapsOthers(room) {
			continue
		}
		if others, collides := m.collidesWithOthers(room); collides {
			for _, other := range others {
				if other == previousRoom {
					continue
				}
				m.tryConnectRooms(room, other)
				break
			}
		}
		previousRoom.SetNeighbor(dir, room)
		room.SetNeighbor((dir+2)%4, previousRoom)
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
	case data.North:
		x -= width >> 1
		y = neighbor.Top() - height
	case data.East:
		x = neighbor.Right()
		y -= height >> 1
	case data.South:
		x -= width >> 1
		y = neighbor.Bottom()
	case data.West:
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

func (m *Map) collidesWithOthers(room IRoom) ([]IRoom, bool) {
	rooms := make([]IRoom, 0)
	for _, other := range m.rooms {
		if room.CollidesWith(other) {
			rooms = append(rooms, other)
		}
	}
	return rooms, len(rooms) > 0
}

func (m *Map) tryConnectRooms(r1, r2 IRoom) bool {
	if r1.GetX() != r2.GetX() && r1.GetY() != r2.GetY() {
		return false
	}
	dir := r1.GetNeighborDirection(r2)
	opp := (dir + 2) % 4
	if r1.GetNeighbor(dir) != nil && r2.GetNeighbor(opp) != nil {
		return false
	}
	r1.SetNeighbor(dir, r2)
	r2.SetNeighbor(opp, r1)
	return true
}

func (m *Map) calculateDimensions() {
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
