package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type exchange struct {
	wg  sync.WaitGroup
	ctx context.Context
	signal chan struct{}
}

func printForever(ex *exchange) {
	defer ex.wg.Done()

	for {
		select {
		case <-ex.ctx.Done():
			return
		case <-ex.signal:
			fmt.Println("Time to work")
			time.Sleep(5 * time.Second)
		}
	}
}

func main() {
	ex := exchange{wg: sync.WaitGroup{}}

	ctx, cancelFunc := context.WithCancel(context.Background())
	ex.ctx = ctx
	ex.signal = make(chan struct{})

	ex.wg.Add(1)
	go printForever(&ex)

	for i := 0; i < 20; i++ {
		fmt.Printf("Iteration: %d\n", i)
		select {
		case ex.signal <- struct{}{}:
			fmt.Println("message sent")
		default:
			fmt.Println("skip this time")
		}

		time.Sleep(time.Second)
	}

	cancelFunc()
	ex.wg.Wait()
}
