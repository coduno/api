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

	var m model.Map
	if err = json.Unmarshal(tmap, &m); err != nil {
		return
	}

	tr, err := testRobot(&m, in)
	if err != nil {
		return marshalJSON(&sub, struct {
			Error string
		}{
			Error: err.Error(),
		})
	}

	if _, err = tr.PutWithParent(ctx, sub.Key); err != nil {
		return
	}
	return marshalJSON(&sub, tr.Moves)
}

func testRobot(m *model.Map, in io.Reader) (tr model.RobotTestResults, err error) {
	tr.Moves = append(tr.Moves, model.RobotLogEntry{Position: m.Start, Event: model.Move})
	direction := model.Right
	r := bufio.NewReader(in)
	var picks []int
	for {
		current, err := r.ReadBytes('\n')
		if len(current) == 0 {
			break
		} else if err != nil && err != io.EOF {
			return tr, err
		}

		current = bytes.ToLower(current)
		current = bytes.TrimSuffix(current, []byte("\n"))

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
			tr.Moves = append(tr.Moves, tr.Moves[len(tr.Moves)-1])
			picks = append(picks, len(tr.Moves)-1)
		default:
			return tr, errors.New("bad command")
		}
	}
	if m.Finish == tr.Moves[len(tr.Moves)-1].Position {
		tr.ReachedFinish = true
	}

	picked := 0
	for _, pick := range picks {
		if m.CoinAt(tr.Moves[pick].Position) {
			picked++
			tr.Moves[pick].Event = model.Picked
		} else {
			tr.Failed = true
			tr.Moves[pick].Event = model.WrongPick
		}
	}
	if picked != len(m.Coins) {
		tr.Failed = true
	}
	return tr, nil
}

func move(tr *model.RobotTestResults, n int, direction model.RobotDirection, m *model.Map) bool {
	for n > 0 {
		n--

		i := len(tr.Moves) - 1
		tr.Moves = append(tr.Moves, tr.Moves[i].Move(direction))
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
