// Package fieldvalidate provides field value validation for structured log records.
// It allows asserting that specific fields match expected types or value constraints,
// and can be used to drop or flag records that fail validation.
package fieldvalidate

import (
	"fmt"
	"regexp"
	"strconv"
)

// Rule describes a single field validation constraint.
type Rule struct {
	Field   string
	Kind    string // "nonempty", "number", "bool", "regex"
	Pattern *regexp.Regexp // only set when Kind == "regex"
}

// Validator holds a set of rules and applies them to log records.
type Validator struct {
	rules []Rule
}

// None returns a Validator that always passes every record.
func None() *Validator { return &Validator{} }

// New builds a Validator from a slice of Rules.
func New(rules []Rule) *Validator {
	return &Validator{rules: rules}
}

// Valid returns true when the record satisfies all rules.
// Records missing a validated field are considered invalid.
func (v *Validator) Valid(fields map[string]string) bool {
	for _, r := range v.rules {
		val, ok := fields[r.Field]
		if !ok {
			return false
		}
		switch r.Kind {
		case "nonempty":
			if val == "" {
				return false
			}
		case "number":
			if _, err := strconv.ParseFloat(val, 64); err != nil {
				return false
			}
		case "bool":
			if _, err := strconv.ParseBool(val); err != nil {
				return false
			}
		case "regex":
			if r.Pattern == nil || !r.Pattern.MatchString(val) {
				return false
			}
		}
	}
	return true
}

// ParseRule parses a rule string of the form "field:kind" or "field:regex:pattern".
func ParseRule(s string) (Rule, error) {
	parts := splitN(s, ':', 3)
	if len(parts) < 2 {
		return Rule{}, fmt.Errorf("fieldvalidate: invalid rule %q: expected field:kind", s)
	}
	r := Rule{Field: parts[0], Kind: parts[1]}
	switch r.Kind {
	case "nonempty", "number", "bool":
		// valid
	case "regex":
		if len(parts) < 3 || parts[2] == "" {
			return Rule{}, fmt.Errorf("fieldvalidate: regex rule requires a pattern")
		}
		re, err := regexp.Compile(parts[2])
		if err != nil {
			return Rule{}, fmt.Errorf("fieldvalidate: bad regex %q: %w", parts[2], err)
		}
		r.Pattern = re
	default:
		return Rule{}, fmt.Errorf("fieldvalidate: unknown kind %q", r.Kind)
	}
	return r, nil
}

func splitN(s string, sep byte, n int) []string {
	var out []string
	for len(out) < n-1 {
		i := indexByte(s, sep)
		if i < 0 {
			break
		}
		out = append(out, s[:i])
		s = s[i+1:]
	}
	return append(out, s)
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
