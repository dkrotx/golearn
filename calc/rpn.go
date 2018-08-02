package main

import (
    "fmt"
    "errors"
    "strconv"
)

func getOperatorPrecedence(op string) int {
    return operators[op]
}

type operator_info struct {
    value      string
    precedence int
    depth      int
}

func sameOrLowerPrecedence(depth, precedence int, other *operator_info) bool {
    return depth < other.depth || (depth == other.depth && precedence <= other.precedence)
}

func BuildRPN(tokens []string) ([]string, error) {
    var op operator_info
    var depth int
    out := make([]string, 0)
    ops := make([]operator_info, 0)

    for _, tok := range tokens {
        if tok == "(" {
            depth++
        } else if tok == ")" {
            depth--
            if depth < 0 {
                return nil, errors.New("Unexpected closing bracket")
            }
        } else if isOperator(tok) {
            prec := getOperatorPrecedence(tok)
            for ; len(ops) != 0 && sameOrLowerPrecedence(depth, prec, &ops[len(ops)-1]); {
                op, ops = ops[len(ops)-1], ops[:len(ops)-1]
                out = append(out, op.value)
            }
            ops = append(ops, operator_info{value: tok, precedence: prec, depth: depth})
        } else {
            out = append(out, tok)
        }
    }

    if depth != 0 {
        return nil, errors.New("Expecting more closing brackets")
    }

    for ; len(ops) != 0; {
        op, ops = ops[len(ops)-1], ops[:len(ops)-1]
        out = append(out, op.value)
    }

    return out, nil
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

func EvalRPN(expr []string) float64 {
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
