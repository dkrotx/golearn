package main

import (
	"fmt"
	"sync"
	"time"
)

/* based on https://kaviraj.me/understanding-condition-variable-in-go/
 */

type Record struct {
	sync.Mutex
	data string

	cond *sync.Cond
}

func NewRecord() *Record {
	r := Record{}
	r.cond = sync.NewCond(&r)
	return &r
}

func RunExternalTaskForever(rec *Record) {
	for {
		rec.Lock()
		rec.cond.Wait() // will unlock mutex
		rec.Unlock()
		fmt.Println("Data: ", rec.data)
	}
}

func main() {
	rec := NewRecord()

	rec.data = fmt.Sprintf("gopher: %s", time.Now())
	go RunExternalTaskForever(rec)

	// periodically send signal to goroutine
	for {
		time.Sleep(1 * time.Second)
		/*rec.Lock()
		rec.data = fmt.Sprintf("gopher: %s", time.Now())
		rec.Unlock()*/

		rec.cond.Signal()
	}
}