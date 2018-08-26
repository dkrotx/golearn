package main

import (
    "fmt"
    "time"
)

func readStdin(out chan<- string) {
    var name string

    fmt.Print("Enter you name: ")
    if n, err := fmt.Scan(&name); err == nil && n != 0 {
        out<-name
    }
    close(out)
}

func main() {
    ch := make(chan string)
    go readStdin(ch)

    select {
        case name := <-ch:
            fmt.Println("Hello, ", name)
        case <-time.After(5 * time.Second):
            fmt.Println("\nTimed out.")
    }
}
