package data

import "math"

type IMapConfig interface {
	GetMinWidth() int
	GetMaxWidth() int
	GetMinHeight() int
	GetMaxHeight() int
	GetMinRooms() int
	GetMaxRooms() int
}

type MapConfig struct {
	minWidth, maxWidth   int
	minHeight, maxHeight int
	minRooms, maxRooms   int
}

func NewMapConfig(minWidth, maxWidth, minHeight, maxHeight, minRooms, maxRooms int) *MapConfig {
	return &MapConfig{minWidth, maxWidth, minHeight, maxHeight, minRooms, maxRooms}
}

func (c *MapConfig) GetMinWidth() int {
	return int(math.Min(float64(c.minWidth), float64(c.maxWidth)))
}

func (c *MapConfig) GetMaxWidth() int {
	return int(math.Max(float64(c.minWidth), float64(c.maxWidth)))
}

func (c *MapConfig) GetMinHeight() int {
	return int(math.Min(float64(c.minHeight), float64(c.maxHeight)))
}

func (c *MapConfig) GetMaxHeight() int {
	return int(math.Max(float64(c.minHeight), float64(c.maxHeight)))
}

func (c *MapConfig) GetMinRooms() int {
	return int(math.Min(float64(c.minRooms), float64(c.maxRooms)))
}

func (c *MapConfig) GetMaxRooms() int {
	return int(math.Max(float64(c.minRooms), float64(c.maxRooms)))
}
