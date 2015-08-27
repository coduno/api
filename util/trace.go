package util

import (
	"fmt"
	"runtime"
)

type trace struct {
	e error
	t []byte
}

func (t trace) Error() string {
	return fmt.Sprintf("%s\n%s", t.e, t.t)
}

// Trace wraps the passed error and generates a new error that will
// expand into a full stack trace. Be aware that this is expensive,
// as it will stop the world to collect the trace, and should be
// used with caution!
func Trace(err error) error {
	r := trace{
		e: err,
		t: make([]byte, 65536),
	}
	runtime.Stack(r.t, false)
	return r
}
