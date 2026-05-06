// Package fieldclip truncates individual field values to a maximum byte or
// rune length, leaving other fields untouched.
package fieldclip

import "unicode/utf8"

// Clipper truncates named field values that exceed a maximum length.
type Clipper struct {
	fields map[string]int // field name → max rune count
	suffix string
}

// None returns a no-op Clipper that passes every record through unchanged.
func None() *Clipper { return &Clipper{} }

// New creates a Clipper that truncates each named field to maxRunes runes,
// appending suffix (e.g. "…") when a value is actually clipped.
// Returns an error when fields is empty or maxRunes < 1.
func New(fields map[string]int, suffix string) (*Clipper, error) {
	if len(fields) == 0 {
		return nil, fmt.Errorf("fieldclip: at least one field required")
	}
	for f, n := range fields {
		if n < 1 {
			return nil, fmt.Errorf("fieldclip: max runes for %q must be >= 1", f)
		}
	}
	copy := make(map[string]int, len(fields))
	for k, v := range fields {
		copy[k] = v
	}
	return &Clipper{fields: copy, suffix: suffix}, nil
}

// Apply returns a shallow copy of rec with clipped field values.
// If the Clipper is a no-op (None), the original map is returned as-is.
func (c *Clipper) Apply(rec map[string]interface{}) map[string]interface{} {
	if len(c.fields) == 0 {
		return rec
	}
	out := make(map[string]interface{}, len(rec))
	for k, v := range rec {
		max, ok := c.fields[k]
		if !ok {
			out[k] = v
			continue
		}
		s, ok := v.(string)
		if !ok {
			out[k] = v
			continue
		}
		out[k] = clipString(s, max, c.suffix)
	}
	return out
}

// clipString truncates s to at most maxRunes runes, appending suffix if clipped.
func clipString(s string, maxRunes int, suffix string) string {
	if utf8.RuneCountInString(s) <= maxRunes {
		return s
	}
	// Walk rune by rune to find the byte boundary.
	i := 0
	for n := 0; n < maxRunes; n++ {
		_, size := utf8.DecodeRuneInString(s[i:])
		i += size
	}
	return s[:i] + suffix
}
