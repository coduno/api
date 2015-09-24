package test

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"path"
	"strconv"
	"strings"

	"github.com/coduno/api/cache"
	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"github.com/coduno/api/ws"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func init() {
	RegisterTester(Robot, robot)
}

func robot(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission) (err error) {
	log.Debugf(ctx, "Executing robot tester")
	cctx := util.CloudContext(ctx)
	var testMap, stdin io.Reader
	if testMap, err = cache.PutGCS(cctx, util.TemplateBucket, t.Params["tests"]); err != nil {
		return
	}

	if stdin, err = cache.PutGCS(cctx, sub.Code.Bucket, path.Dir(sub.Code.Name)+"/"+util.FileNames["robot"]); err != nil {
		return
	}
	var testMapBytes, stdinBytes []byte
	if stdinBytes, err = ioutil.ReadAll(stdin); err != nil {
		return
	}

	if testMapBytes, err = ioutil.ReadAll(testMap); err != nil {
		return
	}

	var m Map
	if err = json.Unmarshal(testMapBytes, &m); err != nil {
		return
	}
	moves, err := testRobot(m, string(stdinBytes))
	if err != nil {
		// TODO(victorbalan): Pass the error to the ws so the client knows what he`s doing wrong
		return
	}
	var body []byte
	if body, err = json.Marshal(moves); err != nil {
		return
	}
	return ws.Write(sub.Key.Parent(), body)
}

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

func (p TypedPosition) move(d int) TypedPosition {
	tp := p
	tp.Type = "MOVE"
	switch d {
	case up:
		tp.X--
	case down:
		tp.X++
	case left:
		tp.Y--
	case right:
		tp.Y++
	}
	return tp
}

func (move *TypedPosition) validate(m Map) bool {
	if move.X < 0 || move.Y < 0 || move.X > m.Max.X || move.X > m.Max.Y {
		move.Type = "OUT_OF_BOUNDS"
		return false
	}
	for _, val := range m.Obstacles {
		if move.X == val.X && move.Y == val.Y {
			move.Type = "HIT_OBSTACLE"
			return false
		}
	}
	return true
}

type Moves []TypedPosition

func (m Moves) last() TypedPosition {
	return m[len(m)-1]
}

func (mv Moves) validate(pos int, m Map) bool {
	return (&mv[pos]).validate(m)
}

type Map struct {
	Start     Position
	Finish    Position
	Coins     []Position
	Obstacles []Position
	Max       Position
}

const (
	up = iota
	right
	down
	left
)

func turnLeft(d int) int {
	d--
	if d < 0 {
		d = left
	}
	return d
}

func turnRight(d int) int {
	d++
	if d > 3 {
		d = up
	}
	return d
}

func testRobot(m Map, in string) (moves []TypedPosition, err error) {
	moves = append(moves, TypedPosition{Position: m.Start, Type: "MOVE"})

	commands := strings.Split(in, "\n")

	direction := right

	for counter, val := range commands {
		c := strings.Split(val, " ")
		switch c[0] {
		case "MOVE":
			n, err := strconv.Atoi(c[1])
			if err != nil {
				return moves, errors.New("bad move argument")
			}
			nextIsPick := !((counter < len(commands)-1 && commands[counter+1] != "PICK") || counter > len(commands))
			var ok bool
			moves, ok = move(moves, n, direction, m, nextIsPick)
			if !ok {
				return moves, nil
			}
		case "RIGHT":
			direction = turnRight(direction)
		case "LEFT":
			direction = turnLeft(direction)
		case "PICK":
			moves = pick(moves, m.Coins)
		default:
			return moves, errors.New("bad command")
		}
	}
	return
}

func pick(moves Moves, coins []Position) []TypedPosition {
	pos := moves.last()
	if !isCoin(pos.Position, coins) {
		pos.Type = "WRONG_PICK"
	} else {
		pos.Type = "PICKED"
	}
	moves[len(moves)-1] = pos
	return moves
}

func move(moves Moves, n, direction int, m Map, nextIsPick bool) ([]TypedPosition, bool) {
	for i := 0; i < n; i++ {
		moves = append(moves, moves.last().move(direction))
		if !moves.validate(len(moves)-1, m) {
			return moves, false
		}

		if isCoin(moves.last().Position, m.Coins) && (!nextIsPick || i < n-1) {
			move := moves.last()
			move.Type = "MISSED_COIN"
			moves[len(moves)-1] = move
		}
	}
	return moves, true
}

func isCoin(pos Position, coins []Position) bool {
	for _, val := range coins {
		if pos.Equals(val) {
			return true
		}
	}
	return false
}
