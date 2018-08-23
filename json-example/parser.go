package main

import (
	"encoding/json"
)

type StepInfo struct {
	Name    string  `json:"name"`
	Success bool    `json:"success"`
	Amount  int     `json:"amount"`
}

type DayEvent struct {
	StartTime    Timestamp  `json:"start_time"`
	EndTime      Timestamp  `json:"end_time"`
	Emails       []string   `json:"inform"`
	Steps        []StepInfo `json:"steps"`
}

func LoadDayEvent(blob []byte) DayEvent {
	var day DayEvent

	json.Unmarshal(blob, &day)
	return day
}

