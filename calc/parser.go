package main

import (
    "regexp"
    "errors"
)

var operators = map[string]int {
    "+": 1,
    "-": 1,
    "*": 2,
    "/": 2,
}

func isOperator(s string) bool {
   _, found := operators[s]
   return found
}

func isNumber(s string) bool {
    return !(isOperator(s) || s == "(" || s == ")")
}

func ParseString(s string) ([]string, error) {
    re := regexp.MustCompile(`\d+|\(|\)|\+|\-|\*|\/`)
    tokens := re.FindAllString(s, -1)

    var unary bool
    res := make([]string, 0)

    for i, tok := range tokens {
        if isOperator(tok) {
            if i == 0 || !isNumber(tokens[i-1]) {
                if tok != "-" {
                    return nil, errors.New("Only '-' allowed as unary sing")
                }
                unary = true
            }
         } else {
            if unary {
                if !isNumber(tok) {
                    return nil, errors.New("Unary '-' i  wrong place")
                }
                tok = "-" + tok
                unary = false
            }
        }

        if !unary { res = append(res, tok) }
    }

    return res, nil
}
