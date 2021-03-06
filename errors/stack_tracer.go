package main

import "github.com/pkg/errors"


// official approach recommended by https://godoc.org/github.com/pkg/errors

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func StackTrace(e error) errors.StackTrace {
	err, ok := e.(stackTracer)
	if !ok {
		return nil
	}
	return err.StackTrace()
}

func ErrorWithEarliestStackTrace(e error) error {
	if e == nil {
		return nil
	}

	cause := errors.Cause(e)
	if StackTrace(cause) != nil {
		return cause
	}
	return e
}