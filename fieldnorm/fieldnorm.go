// Package fieldnorm normalises field values in parsed log records.
// It supports trimming whitespace, lowercasing, and uppercasing string
// values for named fields, making downstream filtering and aggregation
// more robust against inconsistent log formatting.
package fieldnorm

import "strings"

// Op is a normalisation operation applied to a field value.
type Op int

const (
	OpTrim      Op = iota // strip leading/trailing whitespace
	OpLower               // convert to lowercase
	OpUpper               // convert to uppercase
	OpTrimLower           // trim then lowercase
	OpTrimUpper           // trim then uppercase
)

// Rule pairs a field name with the operation to apply.
type Rule struct {
	Field string
	op    Op
}

// Normaliser applies a set of rules to a record.
type Normaliser struct {
	rules []Rule
}

// None returns a Normaliser that makes no changes.
func None() *Normaliser { return &Normaliser{} }

// New returns a Normaliser that applies the given rules.
func New(rules []Rule) *Normaliser {
	r := make([]Rule, len(rules))
	copy(r, rules)
	return &Normaliser{rules: r}
}

// Apply returns a shallow copy of fields with normalisation applied.
// Fields not covered by any rule are passed through unchanged.
func (n *Normaliser) Apply(fields map[string]interface{}) map[string]interface{} {
	if len(n.rules) == 0 {
		return fields
	}
	out := make(map[string]interface{}, len(fields))
	for k, v := range fields {
		out[k] = v
	}
	for _, r := range n.rules {
		v, ok := out[r.Field]
		if !ok {
			continue
		}
		s, ok := v.(string)
		if !ok {
			continue
		}
		out[r.Field] = applyOp(s, r.op)
	}
	return out
}

// Rules returns a copy of the configured rules.
func (n *Normaliser) Rules() []Rule {
	out := make([]Rule, len(n.rules))
	copy(out, n.rules)
	return out
}

func applyOp(s string, op Op) string {
	switch op {
	case OpTrim:
		return strings.TrimSpace(s)
	case OpLower:
		return strings.ToLower(s)
	case OpUpper:
		return strings.ToUpper(s)
	case OpTrimLower:
		return strings.ToLower(strings.TrimSpace(s))
	case OpTrimUpper:
		return strings.ToUpper(strings.TrimSpace(s))
	default:
		return s
	}
}
