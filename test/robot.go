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
)

func turnLeft(d model.RobotDirection) model.RobotDirection {
	d--
	if d < 0 {
		d = model.Left
	}
	return d
}

func turnRight(d model.RobotDirection) model.RobotDirection {
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
	testrc, err := util.Load(util.CloudContext(ctx), util.TemplateBucket, t.Params["tests"])
	if err != nil {
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

	var m *model.Map
	if err = json.Unmarshal(tmap, m); err != nil {
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

func testRobot(m *model.Map, in io.Reader) (tr model.RobotTestResults, err error) {
	tr.Moves = append(tr.Moves, model.RobotLogEntry{Position: m.Start, Event: model.Move})
	direction := model.Right
	r := bufio.NewReader(in)
	picked := 0

	for {
		current, err := r.ReadBytes('\n')
		if err == io.EOF || len(current) == 0 {
			break
		} else if err != nil {
			return tr, err
		}

		current = bytes.ToLower(current)

		c := bytes.Split(current, []byte(" "))
		switch string(c[0]) {
		case "move":
			if len(c) < 2 {
				return tr, errors.New("move missing parameter")
			}

			n, err := strconv.Atoi(string(c[1]))
			if err != nil {
				return tr, err
			}

			if ok := move(&tr, n, direction, m); !ok {
				return tr, nil
			}
		case "right":
			direction = turnRight(direction)
		case "left":
			direction = turnLeft(direction)
		case "pick":
			if pick(&tr, m) {
				picked++
			}
		default:
			return tr, errors.New("bad command")
		}
	}
	if m.Finish == tr.Moves[len(tr.Moves)-1].Position {
		tr.ReachedFinish = true
	}
	if picked != len(m.Coins) {
		tr.Failed = true
	}
	return tr, nil
}

func pick(tr *model.RobotTestResults, m *model.Map) bool {
	i := len(tr.Moves) - 1

	// Check whether there is a coin.
	if !m.CoinAt(tr.Moves[i].Position) {
		tr.Moves[i].Event = model.WrongPick
		tr.Failed = true
		return false
	}

	// Check whether the coin has been picked up already.
	for _, m := range tr.Moves {
		if m == (model.RobotLogEntry{tr.Moves[i].Position, model.Picked}) {
			tr.Moves[i].Event = model.WrongPick
			tr.Failed = true
			return false
		}
	}

	tr.Moves[i].Event = model.Picked
	return true
}

func move(tr *model.RobotTestResults, n int, direction model.RobotDirection, m *model.Map) bool {
	for n > 0 {
		n--

		i := len(tr.Moves) - 1
		tr.Moves = append(tr.Moves, tr.Moves[i])
		i++

		if tr.Moves[i].Validate(m) {
			continue
		}

		tr.Failed = true
		tr.ReachedFinish = false
		return false
	}
	return true
}
