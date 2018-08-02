package main

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestBuildRPN(t *testing.T) {
    tests := []struct {
        input  []string
        expect []string
    }{
        { []string{"1"}, []string{"1"} },
        { []string{"2", "+", "3"}, []string{"2", "3", "+"} },
        { []string{"2", "+", "3", "*", "4", "/", "2"}, []string{"2", "3", "4", "*", "2", "/", "+"} },
        { []string{"2", "-", "(", "3", "+", "1", ")"}, []string{"2", "3", "1", "+", "-"} },
        { []string{"2", "-", "(", "3", ")"}, []string{"2", "3", "-"} },
    }

    assert := assert.New(t)
    for _, test := range tests {
        rpn, err := BuildRPN(test.input)
        if err != nil {
            t.Errorf("Failed to build RPN for %v", test.input)
        }

        assert.Equal(test.expect, rpn, test.input)
    }
}

func TestEvalRPN(t *testing.T) {
    tests := []struct {
        input  []string
        expect float64
    }{
        {[]string{"1"}, 1},
        {[]string{"2", "3", "+"}, 5},
        {[]string{"2", "3", "1", "+", "-"}, -2},
        {[]string{"2", "3", "4", "*", "2", "/", "+"}, 8},
        {[]string{"1", "2", "/"}, 0.5},
    }
    assert := assert.New(t)

    for _, test := range tests {
        assert.InDelta(test.expect, EvalRPN(test.input), 0.001, test.input)
    }
}
