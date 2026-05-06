// Package fieldgroup merges multiple fields into a single composite field.
package fieldgroup

import (
	"fmt"
	"strings"
)

// Grouper merges a set of source fields into a new destination field.
type Grouper struct {
	dest      string
	sources   []string
	separator string
}

// None returns a no-op Grouper that passes records through unchanged.
func None() *Grouper { return &Grouper{} }

// New creates a Grouper that concatenates sources into dest using sep.
// dest and every element of sources must be non-empty.
func New(dest string, sources []string, sep string) (*Grouper, error) {
	if dest == "" {
		return nil, fmt.Errorf("fieldgroup: dest field name must not be empty")
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("fieldgroup: at least one source field is required")
	}
	for _, s := range sources {
		if s == "" {
			return nil, fmt.Errorf("fieldgroup: source field names must not be empty")
		}
	}
	return &Grouper{dest: dest, sources: sources, separator: sep}, nil
}

// Apply returns a copy of rec with the grouped field added.
// Fields missing from rec are treated as empty strings.
// If the Grouper is a no-op (None), rec is returned unchanged.
func (g *Grouper) Apply(rec map[string]any) map[string]any {
	if g.dest == "" {
		return rec
	}
	parts := make([]string, 0, len(g.sources))
	for _, s := range g.sources {
		if v, ok := rec[s]; ok {
			parts = append(parts, fmt.Sprintf("%v", v))
		} else {
			parts = append(parts, "")
		}
	}
	out := make(map[string]any, len(rec)+1)
	for k, v := range rec {
		out[k] = v
	}
	out[g.dest] = strings.Join(parts, g.separator)
	return out
}
