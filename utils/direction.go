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

func (d Direction) NewPos(x, y int) (int, int) {
	switch d {
	case Up:
		return x, y - 1
	case Right:
		return x + 1, y
	case Down:
		return x, y + 1
	case Left:
		return x - 1, y
	}
	return 0, 0
}
