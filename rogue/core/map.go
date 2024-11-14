package core

import "rogue/data"

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
	return &Map{Rectangle: NewRectangle(0, 0, 2000, 1000)}
}

func (m *Map) GetStart() (int, int) {
	return m.start.GetX(), m.start.GetY()
}

func (m *Map) GetEnd() (int, int) {
	return m.end.GetX(), m.end.GetY()
}

// Private //////////////////////////////////////////////////////////////////////

func (m *Map) newRoom() IRectangle {
	return NewRectangle(0, 0, 0, 0)
}
