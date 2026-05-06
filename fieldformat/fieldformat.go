// Package fieldformat applies printf-style format strings to named fields,
// writing the result into a new or existing destination field.
package fieldformat

import (
	"fmt"
	"strings"
)

// Formatter rewrites one field value using a format spec.
type Formatter struct {
	field  string
	dest   string
	fmt    string
}

// None is a no-op Formatter that passes records through unchanged.
func None() *Formatter { return nil }

// New creates a Formatter that reads field, applies fmtSpec (a fmt verb such as
// "%.3f" or "%08d"), and writes the result to dest. If dest is empty the
// source field is overwritten.
func New(field, dest, fmtSpec string) (*Formatter, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldformat: field name must not be empty")
	}
	if fmtSpec == "" {
		return nil, fmt.Errorf("fieldformat: format spec must not be empty")
	}
	if !strings.Contains(fmtSpec, "%") {
		return nil, fmt.Errorf("fieldformat: format spec %q contains no verb", fmtSpec)
	}
	if dest == "" {
		dest = field
	}
	return &Formatter{field: field, dest: dest, fmt: fmtSpec}, nil
}

// Apply formats the named field in rec and stores the result in dest.
// If the source field is absent the record is returned unchanged.
// A nil Formatter is a no-op.
func (f *Formatter) Apply(rec map[string]any) map[string]any {
	if f == nil {
		return rec
	}
	v, ok := rec[f.field]
	if !ok {
		return rec
	}
	out := make(map[string]any, len(rec))
	for k, val := range rec {
		out[k] = val
	}
	out[f.dest] = fmt.Sprintf(f.fmt, v)
	return out
}
