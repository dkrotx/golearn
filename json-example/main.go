package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"encoding/json"
	"time"
)

func main() {
	text, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		panic(fmt.Sprintf("Failed to parse input json: %s", err))
	}

	day := LoadDayEvent(text)
	fmt.Printf("day structure:\n%v\n\n", day)

	// use obtained values
	fmt.Printf("Duration is %v\n", time.Duration(day.EndTime.UnixNano() - day.StartTime.UnixNano()))

	js, err := json.Marshal(day)
	if err != nil {
		panic(fmt.Sprintf("Failed to make output json: %s", err))
	}

	fmt.Printf("Marshalled:\n%s\n\n", js)
}