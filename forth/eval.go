//go:build !solution

package main

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

type Evaluator struct {
	words     map[string][]func() bool
	variables map[string][]int
	stack     []int
}

func (e *Evaluator) funcDefition(args []string) bool {
	InstructionName := args[1]
	InstructionSet := make([]func() bool, 0)
	for i := 2; args[i] != ";"; i++ {
		word, ok := e.words[args[i]]
		if !ok {
			return false
		}
		InstructionSet = append(InstructionSet, word...)
	}
	e.words[InstructionName] = InstructionSet
	return true
}
func (e *Evaluator) varDefition(args []string) bool {
	VarName := args[1]
	ind := -1
	for i := 2; i < len(args)-1; i++ {
		for word := range e.words {
			ok := strings.Contains(args[i], word)
			if ok {
				ind = i
			}
		}
	}
	for i := 2; args[i] != ";"; i++ {
		if i == ind {
			word, ok := e.words[args[ind]]
			for _, f := range word {
				ok := f()
				if !ok {
					return false
				}
			}
			if !ok {
				return false
			}
		} else {
			val, ok := e.variables[args[i]]
			if ok {
				e.stack = append(e.stack, val...)
			} else {
				n, _ := strconv.ParseInt(args[i], 10, 32)
				e.stack = append(e.stack, int(n))
			}
		}
	}
	e.variables[VarName] = e.stack
	e.stack = []int{}
	return true
}

func (e *Evaluator) isFuncDefinition(args []string) bool {
	cnt := 0
	for i := 2; i < len(args)-1; i++ {
		for word := range e.words {
			if strings.Contains(args[i], word) {
				cnt++
			}
		}
	}
	if cnt == len(args[2:len(args)-1]) {
		return true
	}
	return false
}

func (e *Evaluator) wordDefition(args []string) bool {
	if isNumber(args[1]) {
		return false
	}
	if !e.isFuncDefinition(args) {
		ok := e.varDefition(args)
		if !ok {
			return false
		}
	} else {
		ok := e.funcDefition(args)
		if !ok {
			return false
		}
	}
	return true
}
func (e *Evaluator) dup() bool {
	if len(e.stack) == 0 {
		return false
	}
	e.stack = append(e.stack, e.stack[len(e.stack)-1])
	return true
}
func (e *Evaluator) over() bool {
	if len(e.stack) < 2 {
		return false
	}
	e.stack = append(e.stack, e.stack[len(e.stack)-2])
	return true
}
func (e *Evaluator) drop() bool {
	if len(e.stack) == 0 {
		return false
	}
	e.stack = e.stack[:len(e.stack)-1]
	return true
}
func (e *Evaluator) swap() bool {
	if len(e.stack) < 2 {
		return false
	}
	e.stack[len(e.stack)-1], e.stack[len(e.stack)-2] = e.stack[len(e.stack)-2], e.stack[len(e.stack)-1]
	return true
}
func (e *Evaluator) GetOperands() (int, int, bool) {
	if len(e.stack) < 2 {
		return 0, 0, false
	}
	a := e.stack[len(e.stack)-1]
	b := e.stack[len(e.stack)-2]
	e.stack = e.stack[:len(e.stack)-2]
	return a, b, true
}
func (e *Evaluator) add() bool {
	a, b, enoughOperands := e.GetOperands()
	if !enoughOperands {
		return false
	}
	e.stack = append(e.stack, a+b)
	return true
}
func (e *Evaluator) minus() bool {
	a, b, enoughOperands := e.GetOperands()
	if !enoughOperands {
		return false
	}
	e.stack = append(e.stack, b-a)
	return true
}
func (e *Evaluator) mul() bool {
	a, b, enoughOperands := e.GetOperands()
	if !enoughOperands {
		return false
	}
	e.stack = append(e.stack, a*b)
	return true
}
func (e *Evaluator) div() bool {
	a, b, enoughOperands := e.GetOperands()
	if !enoughOperands || a == 0 {
		return false
	}
	e.stack = append(e.stack, b/a)
	return true
}

func NewEvaluator() *Evaluator {
	e := &Evaluator{
		stack: []int{},
	}
	e.variables = make(map[string][]int)
	e.words = map[string][]func() bool{
		"dup":  {e.dup},
		"over": {e.over},
		"drop": {e.drop},
		"swap": {e.swap},
		"+":    {e.add},
		"-":    {e.minus},
		"*":    {e.mul},
		"/":    {e.div},
	}
	return e
}

func isNumber(s string) bool {
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

// Process evaluates sequence of words or definition.
//
// Returns resulting stack state and an error.
func (e *Evaluator) Process(row string) ([]int, error) {
	args := strings.Fields(row)
	for i := 0; i < len(args); i++ {
		args[i] = strings.ToLower(args[i])
	}
	for _, arg := range args {
		_, isKnownCommand := e.words[arg]
		_, isVariable := e.variables[arg]
		switch {
		case arg == ":":
			ok := e.wordDefition(args)
			if !ok {
				return e.stack, errors.New(arg)
			}
			return e.stack, nil
		case isNumber(arg):
			number, _ := strconv.ParseInt(arg, 10, 32)
			e.stack = append(e.stack, int(number))
		case isVariable:
			e.stack = append(e.stack, (e.variables[arg])...)
		case isKnownCommand:
			for _, f := range e.words[arg] {
				ok := f()
				if !ok {
					return e.stack, errors.New(arg)
				}
			}
		default:
			return e.stack, errors.New(arg)
		}
	}
	return e.stack, nil
}
