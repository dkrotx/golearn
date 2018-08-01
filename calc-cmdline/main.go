package main

/* 
 * Evaluate command-line args as mathematical expression which contains +,-,*,/  and numbers
 */

import (
    "os"
    "fmt"
    "strconv"
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

func getOperatorPrecedence(op string) int {
    return operators[op]
}

func parseExpression(tokens []string) []string {
    var x string
    out := make([]string, 0)
    ops := make([]string, 0)

    for _, tok := range tokens {
        if isOperator(tok) {
            prec := getOperatorPrecedence(tok)
            for ; len(ops) != 0 && prec <= getOperatorPrecedence(ops[len(ops)-1]); {
                x, ops = ops[len(ops)-1], ops[:len(ops)-1]
                out = append(out, x)
            }
            ops = append(ops, tok)
        } else {
            out = append(out, tok)
        }
    }

    for ; len(ops) != 0; {
        x, ops = ops[len(ops)-1], ops[:len(ops)-1]
        out = append(out, x)
    }

    return out
}

func evalSingleOp(a, b float64, op string) float64 {
    switch op {
        case "+": return a + b
        case "-": return a - b
        case "*": return a * b
        case "/": return a / b
        default: panic(fmt.Sprintf("Unknown operator '%v'", op))
    }
}

func evalExpression(expr []string) float64 {
    stack := make([]float64, 0)

    for _, val := range expr {
        if isOperator(val) {
            a, b := stack[len(stack)-2], stack[len(stack)-1]
            stack = stack[:len(stack)-2]
            stack = append(stack, evalSingleOp(a, b, val))
        } else {
            num, _ := strconv.ParseFloat(val, 64)
            stack = append(stack, num)
        }
    }

    return stack[0]
}

func main() {
    expr := parseExpression(os.Args[1:])
    fmt.Printf("RPN: %v\n", expr)
    fmt.Printf("RES: %.2f\n", evalExpression(expr))
}
