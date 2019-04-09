package main

import (
	"fmt"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	time.Sleep(2 * time.Second)
	fmt.Fprintf(w, "Welcome to my website!\n")
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":12345", nil)
}