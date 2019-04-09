package main

import "go.uber.org/zap"

// NewErrorWFields carries fields that can be logged
func NewErrorWFields(err error, fields ...zap.Field) error {
	return errorWFields{err, fields}
}

//func (e errorWFields) Cause() error {
//	type causer interface {
//		Cause() error
//	}
//
//	if cause, ok := e.error.(causer); ok {
//		return cause.Cause()
//	}
//	return e.error
//}

// GetFields retrieves context fields from an error and its causes if any
func GetFields(err error) (fields []zap.Field) {
	for err != nil {
		if e, hasFields := err.(errorWFields); hasFields {
			fields = append(fields, e.fields...)
			err = e.error
		}

		causer, hasCause := err.(interface{ Cause() error })
		if hasCause {
			err = causer.Cause()
		} else {
			err = err
		}
	}

	return fields
}

type errorWFields struct {
	error
	fields []zap.Field
}
