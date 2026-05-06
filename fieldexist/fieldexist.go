// Package fieldexist filters log records based on the presence or absence
// of specific fields, regardless of their values.
package fieldexist

import "strings"

// Filter decides whether a record passes based on field existence rules.
type Filter struct {
	require []string
	exclude []string
}

// None returns a no-op Filter that accepts every record.
func None() *Filter { return &Filter{} }

// New creates a Filter that requires all fields in require to be present
// and rejects records that contain any field listed in exclude.
// Field names are compared case-insensitively.
func New(require, exclude []string) (*Filter, error) {
	norm := func(ss []string) []string {
		out := make([]string, len(ss))
		for i, s := range ss {
			out[i] = strings.ToLower(strings.TrimSpace(s))
		}
		return out
	}
	return &Filter{
		require: norm(require),
		exclude: norm(exclude),
	}, nil
}

// Apply returns true when the record satisfies all existence constraints.
func (f *Filter) Apply(fields map[string]interface{}) bool {
	if len(f.require) == 0 && len(f.exclude) == 0 {
		return true
	}

	// Build a lower-cased key set for O(1) lookup.
	keys := make(map[string]struct{}, len(fields))
	for k := range fields {
		keys[strings.ToLower(k)] = struct{}{}
	}

	for _, r := range f.require {
		if _, ok := keys[r]; !ok {
			return false
		}
	}
	for _, e := range f.exclude {
		if _, ok := keys[e]; ok {
			return false
		}
	}
	return true
}

// RequiredFields returns a copy of the required field list.
func (f *Filter) RequiredFields() []string {
	out := make([]string, len(f.require))
	copy(out, f.require)
	return out
}

// ExcludedFields returns a copy of the excluded field list.
func (f *Filter) ExcludedFields() []string {
	out := make([]string, len(f.exclude))
	copy(out, f.exclude)
	return out
}
