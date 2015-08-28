package test

import (
	"strconv"

	"golang.org/x/net/context"
	"google.golang.org/cloud/datastore"
)

// Tester is a reference to a unique implementation of a TesterFunc.
type Tester int

// TesterFunc is a function that will compute skills for the referenced
// Result. It will call Taskers of encapsulated tasks on demand.
type TesterFunc func(ctx context.Context, resultKey *datastore.Key) error

const (
	Simple Tester = 1 + iota
	maxTester
)

var testers = make([]TesterFunc, maxTester)

// RegisterTester registers a function to be called lated via Call. It is
// usually called from the init function of a package that contains a Tester.
func RegisterTester(t Tester, f TesterFunc) {
	if t >= maxTester {
		panic("test: RegisterTester of unknown tester function")
	}
	testers[t] = f
}

// Call looks up a registered Resulter and calls it.
func (t Tester) Call(ctx context.Context, resultKey *datastore.Key) error {
	if t > 0 && t < maxTester {
		f := testers[t]
		if f != nil {
			return f(ctx, resultKey)
		}
	}
	panic("logic: requested resulter function #" + strconv.Itoa(int(t)) + " is unavailable")
}
