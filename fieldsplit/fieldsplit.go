// Package fieldsplit provides a pipeline stage that splits a single string
// field into multiple fields using a configurable delimiter.
package fieldsplit

import (
	"fmt"
	"strings"
)

// Splitter splits a source field into one or more destination fields.
type Splitter struct {
	src   string
	dests []string
	sep   string
}

// None is a no-op Splitter that passes records through unchanged.
var None = &Splitter{}

// New creates a Splitter that splits src on sep and assigns the resulting
// parts to dests in order. Extra parts are discarded; missing parts are
// set to an empty string.
func New(src string, dests []string, sep string) (*Splitter, error) {
	if src == "" {
		return nil, fmt.Errorf("fieldsplit: source field must not be empty")
	}
	if len(dests) == 0 {
		return nil, fmt.Errorf("fieldsplit: at least one destination field is required")
	}
	for i, d := range dests {
		if d == "" {
			return nil, fmt.Errorf("fieldsplit: destination field at index %d must not be empty", i)
		}
	}
	if sep == "" {
		return nil, fmt.Errorf("fieldsplit: separator must not be empty")
	}
	return &Splitter{src: src, dests: dests, sep: sep}, nil
}

// Apply splits the source field in rec and writes destination fields.
// The original source field is preserved. If the source field is absent,
// the record is returned unchanged.
func (s *Splitter) Apply(rec map[string]any) map[string]any {
	if s.src == "" {
		return rec
	}
	v, ok := rec[s.src]
	if !ok {
		return rec
	}
	raw := fmt.Sprintf("%v", v)
	parts := strings.SplitN(raw, s.sep, len(s.dests))
	out := make(map[string]any, len(rec)+len(s.dests))
	for k, val := range rec {
		out[k] = val
	}
	for i, dest := range s.dests {
		if i < len(parts) {
			out[dest] = parts[i]
		} else {
			out[dest] = ""
		}
	}
	return out
}
