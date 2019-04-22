package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ticker := time.NewTicker(time.Second)
	var prev time.Time
	var done bool

	for !done {
		fmt.Println("For exit press ^C twice shortly")
		select {
		case <-ticker.C:
		case <-c:
			now := time.Now()
			if now.Sub(prev) < time.Second {
				done = true
			}
			prev = now
		}
	}
}
