package main

import (
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

func requestBackend(url string) error {
	// let's say we're making request to remote via HTTP
	// It has failed, and we return error to indicate this
	//
	// Try to switch between 503 and 404 HTTP codes
	return NewHTTPError(errors.New("failed to connect to backend"), http.StatusBadGateway)
}

func middleendCallSomething(domain string) error {
	if err := requestBackend("http://" + domain + "/favicon.ico"); err != nil {
		// usually we need to describe error here
		return errors.Wrap(err, "can't handle domain")
	}

	return nil
}

func main() {
	if err := middleendCallSomething("google.com"); err != nil {
		// useful things from error:
		// - message (with all descriptions)
		// - stacktrace
		// - check whatever original error is temporary or not
		fmt.Printf("Temporary: %v\n", IsTemporaryCause(err))
		fmt.Printf("Message: %s\n", err.Error())
		fmt.Println("Stack: ", SPrintStackTrace(err, -1))
	}
}
