package main

import (
	"fmt"
	"github.com/pkg/errors"
	"math"
)

type Temporary interface {
	IsTemporary() bool
}

func IsTemporaryError(err error) bool {
	te, ok := err.(Temporary)
	return ok && te.IsTemporary()
}

func IsTemporaryCause(err error) bool {
	return IsTemporaryError(errors.Cause(err))
}

// SPrintStackTrace prints stack trace with desired depth (-1 for unlimited)
func SPrintStackTrace(err error, depth int) string {
	e, ok := errors.Cause(err).(stackTracer)
	if !ok {
		return "[NO STACK AVAILABLE]"
	}

	if depth == -1 {
		depth = math.MaxInt32
	}

	st := e.StackTrace()
	if len(st) < depth {
		depth = len(st)
	}

	return fmt.Sprintf("%+v", st[0:depth])
}

