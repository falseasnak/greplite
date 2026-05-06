// Package fieldcount provides a filter that accepts or rejects log records
// based on the number of parsed fields they contain. This is useful for
// skipping sparse or malformed records that lack sufficient structure.
package fieldcount

import "fmt"

// Filter holds the min/max field-count bounds.
type Filter struct {
	min int
	max int // -1 means no upper bound
}

// None returns a Filter that accepts every record regardless of field count.
func None() *Filter {
	return &Filter{min: 0, max: -1}
}

// New creates a Filter that accepts records whose field count is in [min, max].
// Pass max == -1 to impose no upper bound.
// Returns an error if min < 0 or (max != -1 && max < min).
func New(min, max int) (*Filter, error) {
	if min < 0 {
		return nil, fmt.Errorf("fieldcount: min must be >= 0, got %d", min)
	}
	if max != -1 && max < min {
		return nil, fmt.Errorf("fieldcount: max (%d) must be >= min (%d)", max, min)
	}
	return &Filter{min: min, max: max}, nil
}

// Accept returns true when the number of fields in fields satisfies the
// configured bounds.
func (f *Filter) Accept(fields map[string]interface{}) bool {
	n := len(fields)
	if n < f.min {
		return false
	}
	if f.max != -1 && n > f.max {
		return false
	}
	return true
}

// String returns a human-readable description of the filter bounds.
func (f *Filter) String() string {
	if f.min == 0 && f.max == -1 {
		return "fieldcount:none"
	}
	if f.max == -1 {
		return fmt.Sprintf("fieldcount:min=%d", f.min)
	}
	return fmt.Sprintf("fieldcount:min=%d,max=%d", f.min, f.max)
}
