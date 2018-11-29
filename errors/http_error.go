package main

import "github.com/pkg/errors"

type HTTPError struct {
	error
	code int
}

func NewHTTPError(err error, code int) HTTPError {
	return HTTPError{err, code}
}

func (e HTTPError) StackTrace() errors.StackTrace {
	return StackTrace(e.error)
}

func (e HTTPError) IsTemporary() bool {
	category := e.code / 100
	// 4xx errors are request errors, no sense to retry (w/o modifying request)
	return category != 4
}
