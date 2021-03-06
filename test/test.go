package test

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"

	"github.com/coduno/api/model"
	"github.com/coduno/api/ws"

	"golang.org/x/net/context"
)

type Tester int

type TesterFunc func(ctx context.Context, t model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) error

const (
	Simple Tester = 1 + iota
	Junit
	Diff
	IO
	Robot
	CoderJunit
	SpringInt
	maxTester
)

var testers = make([]TesterFunc, maxTester)

type ErrMissingParam string

func (e ErrMissingParam) Error() string {
	return fmt.Sprintf("missing parameter %q", e)
}

// RegisterTester registers a function to be called lated via Call. It is
// usually called from the init function of a package that contains a Tester.
func RegisterTester(t Tester, f TesterFunc) {
	if t >= maxTester {
		panic("test: RegisterTester of unknown tester function")
	}
	testers[t] = f
}

// Call looks up a registered Resulter and calls it.
func (t Tester) Call(ctx context.Context, test model.KeyedTest, sub model.KeyedSubmission, ball io.Reader) error {
	if t > 0 && t < maxTester {
		f := testers[t]
		if f != nil {
			return f(ctx, test, sub, ball)
		}
	}
	panic("test: requested tester function #" + strconv.Itoa(int(t)) + " is unavailable")
}

// marshalJSON is a helper that will write to the WebSocket identified by key.
func marshalJSON(sub *model.KeyedSubmission, v interface{}) error {
	ww, err := ws.NewWriter(sub.Key.Parent())
	if err != nil {
		return err
	}
	err = json.NewEncoder(ww).Encode(v)
	ww.Close()
	return err
}
