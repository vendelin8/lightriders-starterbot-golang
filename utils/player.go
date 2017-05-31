package utils

type Player struct {
	X, Y, Id int
	IdStr    byte
	LastMove Direction
}

func (p *Player) SetMove(d Direction) {
	p.LastMove = d
}

func (p *Player) Move() {
	switch p.LastMove {
	case Up:
		p.Y--
	case Right:
		p.X++
	case Down:
		p.Y++
	case Left:
		p.X--
	}
}

func (p *Player) MoveField(field *Field) {
	field.Rows[p.Y][p.X] = byte('x')
	p.Move()
	field.Rows[p.Y][p.X] = p.IdStr
}

func (p *Player) WhereIs(in string, field *Field) Direction {
	//searching for next pos of the player in the string field
	if p.X > 0 && in[(p.Y*field.Width+p.X-1)*2] == p.IdStr {
		return Left
	}
	if p.X+1 < field.Width && in[(p.Y*field.Width+p.X+1)*2] == p.IdStr {
		return Right
	}
	if p.Y > 0 && in[((p.Y-1)*field.Width+p.X)*2] == p.IdStr {
		return Up
	}
	return Down
}

func (p *Player) SetAndMoveField(fieldStr string, field *Field) {
	p.LastMove = p.WhereIs(fieldStr, field)
	p.MoveField(field)
}
