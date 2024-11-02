package engine

func CalculateExtents(width, height int) [2]IPoint {
	hw, hh := width>>1, height>>1
	rw, rh := width&1, height&1
	return [2]IPoint{
		NewPoint(-hw, -hh),
		NewPoint(hw+rw, hh+rh),
	}
}
