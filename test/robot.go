package test

import (
	"archive/tar"
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/coduno/api/model"
	"github.com/coduno/api/util"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
)

func turnLeft(d int) int {
	d--
	if d < 0 {
		d = model.Left
	}
	return d
}

func turnRight(d int) int {
	d++
	if d > 3 {
		d = model.Up
	}
	return d
}

func init() {
	RegisterTester(Robot, robot)
}

func robot(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) (err error) {
	log.Debugf(ctx, "Executing robot tester")
	var testrc io.ReadCloser
	if testrc, err = util.Load(util.CloudContext(ctx), util.TemplateBucket, t.Params["tests"]); err != nil {
		return
	}
	defer testrc.Close()

	in := tar.NewReader(ball)
	if _, err = in.Next(); err != nil {
		return err
	}

	var tmap []byte
	if tmap, err = ioutil.ReadAll(testrc); err != nil {
		return
	}

	var m model.Map
	if err = json.Unmarshal(tmap, &m); err != nil {
		return
	}
	tr, err := testRobot(m, in)
	if err != nil {
		// TODO(victorbalan): Pass the error to the ws so the client knows what he's doing wrong
		return
	}
	if _, err = tr.PutWithParent(ctx, sub.Key); err != nil {
		return
	}
	return marshalJSON(&sub, moves)
}

func testRobot(m model.Map, in io.Reader) (tr model.RobotTestResults, err error) {
	tr.Moves = append(tr.Moves, model.TypedPosition{Position: m.Start, Type: "MOVE"})
	direction := model.Right
	r := bufio.NewReader(in)

	current, err := r.ReadBytes('\n')
	if err != nil {
		return tr, err
	}
	current = bytes.TrimSuffix(current, []byte("\n"))
	lastMove := false
	for {
		next, err := r.ReadBytes('\n')
		if err == io.EOF {
			lastMove = true
		} else if err != nil {
			return tr, err
		}

		next = bytes.TrimSuffix(next, []byte("\n"))
		current = bytes.ToLower(current)

		c := bytes.Split(current, []byte(" "))
		switch string(c[0]) {
		case "move":
			n, err := strconv.Atoi(string(c[1]))
			if err != nil {
				return tr, err
			}
			var nextIsPick = false
			if !lastMove {
				nextIsPick = bytes.Compare(next, []byte("pick")) == 0
			}
			if ok := move(&tr, n, direction, m, nextIsPick); !ok {
				return tr, nil
			}
		case "right":
			direction = turnRight(direction)
		case "left":
			direction = turnLeft(direction)
		case "pick":
			pick(&tr, m.Coins)
		case "":
			if lastMove {
				return tr, nil
			}
			fallthrough
		default:
			return tr, errors.New("bad command")
		}
		current = next
	}
	if m.Finish.Equals(tr.Moves.Last().Position) {
		tr.ReachedFinish = true
	}
	return tr, nil
}

func pick(tr *model.RobotTestResults, coins []model.Position) {
	pos := tr.Moves.Last()
	if !isCoin(pos.Position, coins) {
		pos.Type = "WRONG_PICK"
		tr.Failed = true
	} else {
		pos.Type = "PICKED"
	}
	tr.Moves[len(tr.Moves)-1] = pos
}

func move(tr *model.RobotTestResults, n, direction int, m model.Map, nextIsPick bool) bool {
	for i := 0; i < n; i++ {
		tr.Moves = append(tr.Moves, tr.Moves.Last().Move(direction))
		if !tr.Moves.Validate(len(tr.Moves)-1, m) {
			tr.Failed = true
			tr.ReachedFinish = false
			return false
		}

		if isCoin(tr.Moves.Last().Position, m.Coins) && (!nextIsPick || i < n-1) {
			move := tr.Moves.Last()
			move.Type = "MISSED_COIN"
			tr.Moves[len(tr.Moves)-1] = move
			tr.Failed = true
		}
	}
	return true
}

func isCoin(pos model.Position, coins []model.Position) bool {
	for _, val := range coins {
		if pos.Equals(val) {
			return true
		}
	}
	return false
}
