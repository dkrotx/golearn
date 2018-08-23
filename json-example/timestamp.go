package main

import (
	"time"
	"errors"
)

const timestampLayout = "\"2006-01-02T15:04:05.999999999\""

// Timestamp is a wrapper around time.Time with support for JSON unmarshal of ISO8601 formatted timestamps.
type Timestamp struct {
	time.Time
}

// UnmarshalJSON unmarshals JSON strings with timestamps as they are formatted by uOrchestrate.
func (t *Timestamp) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return errors.New("unable to unmarshal null value into time")
	}

	t.Time, err = time.Parse(timestampLayout, string(b))
	return
}
