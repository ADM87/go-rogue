package core

import "math"

type ICopy interface {
	Copy() interface{}
}

func Raycast(x0, y0, x1, y1 int, yScale float64, callback func(x, y int) bool) {
	dx := int(math.Abs(float64(x1 - x0)))
	dy := int(math.Abs(float64(y1-y0)) * yScale)
	sx := 1
	sy := 1

	if x0 > x1 {
		sx = -1
	}
	if y0 > y1 {
		sy = -1
	}

	err := dx - dy
	currentX, currentY := x0, y0

	for {
		if callback(currentX, currentY) {
			break
		}
		if currentX == x1 && currentY == y1 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			currentX += sx
		}
		if e2 < dx {
			err += dx
			currentY += sy
		}
	}
}
