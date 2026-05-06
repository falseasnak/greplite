// Package fieldmatch filters records by matching one or more field values
// against patterns (exact, contains, or regex).
package fieldmatch

import (
	"fmt"
	"regexp"
	"strings"
)

// Mode describes how a field value is matched.
type Mode int

const (
	ModeExact    Mode = iota // field == value
	ModeContains             // strings.Contains
	ModeRegex                // regexp match
)

// Rule is a single field-match constraint.
type Rule struct {
	Field   string
	Mode    Mode
	Pattern string
	re      *regexp.Regexp
}

// Matcher holds a set of rules and applies them to records.
type Matcher struct {
	rules []Rule
}

// None returns a Matcher that accepts every record.
func None() *Matcher { return &Matcher{} }

// New builds a Matcher from the provided rules.
func New(rules []Rule) (*Matcher, error) {
	for i, r := range rules {
		if r.Field == "" {
			return nil, fmt.Errorf("fieldmatch: rule %d has empty field name", i)
		}
		if r.Mode == ModeRegex {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				return nil, fmt.Errorf("fieldmatch: rule %d bad regex: %w", i, err)
			}
			rules[i].re = re
		}
	}
	return &Matcher{rules: rules}, nil
}

// Accept returns true when the record satisfies every rule.
func (m *Matcher) Accept(fields map[string]interface{}) bool {
	for _, r := range m.rules {
		v, ok := fields[r.Field]
		if !ok {
			return false
		}
		s := fmt.Sprintf("%v", v)
		switch r.Mode {
		case ModeExact:
			if s != r.Pattern {
				return false
			}
		case ModeContains:
			if !strings.Contains(s, r.Pattern) {
				return false
			}
		case ModeRegex:
			if !r.re.MatchString(s) {
				return false
			}
		}
	}
	return true
}
