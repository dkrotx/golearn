package main

import (
    "fmt"
    "sync"
    "time"
    "runtime"
)

func timedPrintf(format string, a ...interface{}) {
    t := time.Now()
    newFmt := "[" + t.Format(time.StampMilli) + "] " + format
    fmt.Printf(newFmt, a...)
}

func timedPrintln(s string) {
    timedPrintf(s + "\n")
}

func longCalculation(wid int, wg *sync.WaitGroup) {
    defer wg.Done()

    var res int64
    for i := int64(0); i < 1<<32; i++ {
        res += i % 177
    }

    timedPrintf("Worker #%d is finished\n", wid)
}

const ngoroutines = 4

func main() {
    var wg sync.WaitGroup

    runtime.GOMAXPROCS(1) // play with this
    timedPrintln("Launch goroutines")
    for i := 0; i < ngoroutines; i++ {
        wg.Add(1)
        go longCalculation(i, &wg)
    }

    timedPrintln("Waiting")
    wg.Wait()
    timedPrintln("All goroutines are finished")
}
