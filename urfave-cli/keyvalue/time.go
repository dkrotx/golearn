package main

import (
	"errors"
	"time"
)

type Timestamp struct {
	time.Time
}

const timestampLayout = "\"2006-01-02T15:04:05.999999999\""


// UnmarshalJSON decode JSON strings with Timestamp
func (t *Timestamp) UnmarshalJSON(b []byte) (err error) {
	if string(b) == "null" {
		return errors.New("unable to unmarshal null value into time")
	}

	t.Time, err = time.Parse(timestampLayout, string(b))
	return
}

func (t Timestamp) MarshalJSON() (bytes []byte, err error) {
	return []byte(t.Format(timestampLayout)), nil
}
