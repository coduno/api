package model

//go:generate generator

type RobotTestResults struct {
	Failed        bool
	ReachedFinish bool
	Moves         TypedPositions
}

const (
	Up = iota
	Right
	Down
	Left
)

type Position struct {
	X,
	Y int
}

func (p Position) Equals(pos Position) bool {
	return p.X == pos.X && p.Y == pos.Y
}

type TypedPosition struct {
	Position
	Type string
}

func (p TypedPosition) Move(d int) TypedPosition {
	tp := p
	tp.Type = "MOVE"
	switch d {
	case Up:
		tp.X--
	case Down:
		tp.X++
	case Left:
		tp.Y--
	case Right:
		tp.Y++
	}
	return tp
}

func (p *TypedPosition) Validate(m Map) bool {
	if p.X < 0 || p.Y < 0 || p.X > m.Max.X || p.X > m.Max.Y {
		p.Type = "OUT_OF_BOUNDS"
		return false
	}
	for _, val := range m.Obstacles {
		if p.X == val.X && p.Y == val.Y {
			p.Type = "HIT_OBSTACLE"
			return false
		}
	}
	return true
}

type TypedPositions []TypedPosition

func (m TypedPositions) Last() TypedPosition {
	return m[len(m)-1]
}

func (m TypedPositions) Validate(pos int, tmap Map) bool {
	return (&m[pos]).Validate(tmap)
}

// Map represents a list of obstacles and coins
type Map struct {
	Start     Position
	Finish    Position
	Coins     []Position
	Obstacles []Position
	Max       Position
}
