// Package fieldcalc provides a pipeline stage that computes a new field
// from a simple arithmetic expression over existing numeric fields.
package fieldcalc

import (
	"fmt"
	"strconv"
	"strings"
)

// Op is a supported binary arithmetic operator.
type Op int

const (
	OpAdd Op = iota
	OpSub
	OpMul
	OpDiv
)

// Calc computes a new field value from two source fields and an operator.
type Calc struct {
	dest  string
	left  string
	op    Op
	right string
}

// None returns a no-op Calc that passes records through unchanged.
func None() *Calc { return nil }

// New creates a Calc that writes dest = left op right.
// expr must be in the form "dest=left+right" (operators: + - * /).
func New(expr string) (*Calc, error) {
	for _, sep := range []string{"+", "-", "*", "/"} {
		parts := strings.SplitN(expr, "=", 2)
		if len(parts) != 2 {
			continue
		}
		dest := strings.TrimSpace(parts[0])
		rhs := parts[1]
		var op Op
		var operands []string
		switch {
		case strings.Contains(rhs, "+"):
			operands = strings.SplitN(rhs, "+", 2)
			op = OpAdd
		case strings.Contains(rhs, "-"):
			operands = strings.SplitN(rhs, "-", 2)
			op = OpSub
		case strings.Contains(rhs, "*"):
			operands = strings.SplitN(rhs, "*", 2)
			op = OpMul
		case strings.Contains(rhs, "/"):
			operands = strings.SplitN(rhs, "/", 2)
			op = OpDiv
		default:
			return nil, fmt.Errorf("fieldcalc: no operator found in %q", expr)
		}
		if dest == "" {
			return nil, fmt.Errorf("fieldcalc: empty destination field in %q", expr)
		}
		left := strings.TrimSpace(operands[0])
		right := strings.TrimSpace(operands[1])
		if left == "" || right == "" {
			return nil, fmt.Errorf("fieldcalc: empty operand in %q", expr)
		}
		return &Calc{dest: dest, left: left, op: op, right: right}, nil
		_ = sep
	}
	return nil, fmt.Errorf("fieldcalc: invalid expression %q", expr)
}

// Apply computes the expression and adds the result to a copy of fields.
// If either operand is missing or non-numeric the record is returned unchanged.
func (c *Calc) Apply(fields map[string]string) map[string]string {
	if c == nil {
		return fields
	}
	lv, ok1 := numVal(fields, c.left)
	rv, ok2 := numVal(fields, c.right)
	if !ok1 || !ok2 {
		return fields
	}
	var result float64
	switch c.op {
	case OpAdd:
		result = lv + rv
	case OpSub:
		result = lv - rv
	case OpMul:
		result = lv * rv
	case OpDiv:
		if rv == 0 {
			return fields
		}
		result = lv / rv
	}
	out := make(map[string]string, len(fields)+1)
	for k, v := range fields {
		out[k] = v
	}
	out[c.dest] = strconv.FormatFloat(result, 'f', -1, 64)
	return out
}

func numVal(fields map[string]string, key string) (float64, bool) {
	v, ok := fields[key]
	if !ok {
		return 0, false
	}
	f, err := strconv.ParseFloat(v, 64)
	return f, err == nil
}
