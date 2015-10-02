package model

//go:generate generator

type RobotTestResults struct {
	Failed        bool
	ReachedFinish bool
	Moves         []RobotLogEntry
}

type RobotDirection int

const (
	Up RobotDirection = iota
	Right
	Down
	Left
)

type Position struct {
	X,
	Y int
}

type RobotEvent string

const (
	Move        RobotEvent = RobotEvent("MOVE")
	OutOfBounds            = RobotEvent("OUT_OF_BOUNDS")
	Obstacle               = RobotEvent("HIT_OBSTACLE")
	MissedCoin             = RobotEvent("MISSED_COIN")
	WrongPick              = RobotEvent("WRONG_PICK")
	Picked                 = RobotEvent("PICKED")
)

type RobotLogEntry struct {
	Position
	Event RobotEvent
}

func (p *RobotLogEntry) Validate(m *Map) bool {
	if p.X < 0 || p.Y < 0 || p.X > m.Max.X || p.X > m.Max.Y {
		p.Event = OutOfBounds
		return false
	}
	for _, val := range m.Obstacles {
		if p.X == val.X && p.Y == val.Y {
			p.Event = Obstacle
			return false
		}
	}
	return true
}

type Map struct {
	Start, Finish, Max Position
	Coins, Obstacles   []Position
}

func (m Map) CoinAt(p Position) bool {
	for _, coin := range m.Coins {
		if p == coin {
			return true
		}
	}
	return false
}
