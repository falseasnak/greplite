// Package filter provides field-level filtering for structured log entries.
// It supports matching log fields by key=value pairs, regex patterns, and
// comparison operators for use with parsed JSON or logfmt log lines.
package filter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Op represents a comparison operator used in a filter expression.
type Op int

const (
	OpEqual Op = iota
	OpNotEqual
	OpContains
	OpRegex
)

// Filter represents a single field filter expression (e.g. level=error).
type Filter struct {
	Field   string
	Op      Op
	Value   string
	compiled *regexp.Regexp
}

// Parse parses a filter expression string into a Filter.
// Supported formats:
//   field=value     (exact match)
//   field!=value    (not equal)
//   field~=value    (contains)
//   field/regex/    (regex match)
func Parse(expr string) (*Filter, error) {
	if idx := strings.Index(expr, "!="); idx != -1 {
		return &Filter{Field: expr[:idx], Op: OpNotEqual, Value: expr[idx+2:]}, nil
	}
	if idx := strings.Index(expr, "~="); idx != -1 {
		return &Filter{Field: expr[:idx], Op: OpContains, Value: expr[idx+2:]}, nil
	}
	if idx := strings.Index(expr, "="); idx != -1 {
		return &Filter{Field: expr[:idx], Op: OpEqual, Value: expr[idx+1:]}, nil
	}
	// regex: field/pattern/
	parts := strings.SplitN(expr, "/", 3)
	if len(parts) == 3 && parts[2] == "" {
		re, err := regexp.Compile(parts[1])
		if err != nil {
			return nil, fmt.Errorf("invalid regex in filter %q: %w", expr, err)
		}
		return &Filter{Field: parts[0], Op: OpRegex, Value: parts[1], compiled: re}, nil
	}
	return nil, fmt.Errorf("unrecognized filter expression: %q", expr)
}

// Match reports whether the given fields map satisfies the filter.
func (f *Filter) Match(fields map[string]string) bool {
	val, ok := fields[f.Field]
	if !ok {
		// treat missing field as empty string for equality checks
		val = ""
	}
	switch f.Op {
	case OpEqual:
		return val == f.Value
	case OpNotEqual:
		return val != f.Value
	case OpContains:
		return strings.Contains(val, f.Value)
	case OpRegex:
		if f.compiled != nil {
			return f.compiled.MatchString(val)
		}
		return false
	}
	return false
}

// MatchAll reports whether all filters in the slice match the given fields.
func MatchAll(filters []*Filter, fields map[string]string) bool {
	for _, f := range filters {
		if !f.Match(fields) {
			return false
		}
	}
	return true
}

// numericVal attempts to parse a string as float64 for future numeric ops.
func numericVal(s string) (float64, bool) {
	v, err := strconv.ParseFloat(s, 64)
	return v, err == nil
}

// ensure numericVal is used to avoid lint errors during early development.
var _ = numericVal
