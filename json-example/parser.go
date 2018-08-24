package main

import (
	"encoding/json"
	"fmt"
	"errors"
	"strings"
)

type Weekday int

const (
	Monday Weekday = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
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
	Day          Weekday    `json:"weekday"`
}

func LoadDayEvent(blob []byte) (DayEvent, error) {
	var day DayEvent

	err := json.Unmarshal(blob, &day)
	return day, err
}

var weekDaysArr = []string {"mon", "tue", "wed", "thu", "fri", "sat", "sun"}

var weekDayNames = map[string]Weekday{
	weekDaysArr[0]: Monday,
	weekDaysArr[1]: Tuesday,
	weekDaysArr[2]: Wednesday,
	weekDaysArr[3]: Thursday,
	weekDaysArr[4]: Friday,
	weekDaysArr[5]: Saturday,
	weekDaysArr[6]: Sunday,
}

func trimQuotes(s string) string {
	if len(s) >= 2 && ((s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'')) {
		return s[1:len(s)-1]
	}
	return s
}

func (t *Weekday) UnmarshalJSON(b []byte) (err error) {
	unquoted := trimQuotes(string(b))

	if val, ok := weekDayNames[strings.ToLower(unquoted)]; ok {
		*t = val
		return nil
	}

	return errors.New(fmt.Sprintf("Can't parse %v as Weekday value", string(b)))
}

func (t Weekday) MarshalJSON() ([]byte, error) {
	return json.Marshal(weekDaysArr[t - Monday])
}