package utils

type Direction int

const (
	Up Direction = iota
	Right
	Down
	Left
	Wtf
)

func (d Direction) String() string {
	switch d {
	case Up:
		return "up"
	case Right:
		return "right"
	case Down:
		return "down"
	case Left:
		return "left"
	}
	return "wtf"
}

func (d Direction) Reverse() Direction {
	switch d {
	case Up:
		return Down
	case Right:
		return Left
	case Down:
		return Up
	case Left:
		return Right
	}
	return Wtf
}
