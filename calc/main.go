package main

import (
    "os"
    "fmt"
    "flag"
    "strings"
)

func main() {
    verbose := flag.Bool("v", false, "be vebose")
    flag.Parse()

    raw_expr := strings.Join(flag.Args(), "")
    if raw_expr == "" {
        flag.Usage()
        os.Exit(64)
    }

    tokens, err := ParseString(raw_expr)
    if err != nil {
        panic(fmt.Sprintf(`failed to parse "%v"`, raw_expr))
    }

    if *verbose {
        fmt.Printf("Tokens: %q\n", tokens)
    }

    rpn, err := BuildRPN(tokens)
    if err != nil {
        panic(fmt.Sprintf(`failed to build RPN from "%v" (%v)`, raw_expr, tokens))
    }

    if *verbose {
        fmt.Printf("RPN:    %q\n", rpn)
    }

    fmt.Println(EvalRPN(rpn))
}
