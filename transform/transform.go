// Package transform provides field extraction and reshaping of structured
// log records before they are passed to the output formatter.
//
// A Transform selects a subset of fields from a parsed record map and,
// optionally, renames them.  When no field selectors are configured the
// original record is returned unchanged.
package transform

import (
	"fmt"
	"strings"
)

// Selector describes a single field projection: the source key in the parsed
// record and the (possibly different) key that should appear in the output.
type Selector struct {
	From string // original field name
	To   string // output field name (equals From when no alias is given)
}

// ParseSelector parses a selector expression of the form "field" or
// "field:alias".  An empty string returns an error.
func ParseSelector(s string) (Selector, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return Selector{}, fmt.Errorf("transform: empty selector")
	}
	parts := strings.SplitN(s, ":", 2)
	sel := Selector{From: parts[0]}
	if len(parts) == 2 && parts[1] != "" {
		sel.To = parts[1]
	} else {
		sel.To = parts[0]
	}
	return sel, nil
}

// Transform applies a set of field selectors to a parsed record.  When the
// selector list is empty the original map is returned as-is (zero allocation
// path).  Missing fields are silently skipped unless RequireAll is true.
type Transform struct {
	selectors  []Selector
	requireAll bool
}

// New creates a Transform from a slice of raw selector strings (e.g. from CLI
// flags).  Pass requireAll=true to treat missing fields as an error.
func New(raw []string, requireAll bool) (*Transform, error) {
	sels := make([]Selector, 0, len(raw))
	for _, r := range raw {
		sel, err := ParseSelector(r)
		if err != nil {
			return nil, err
		}
		sels = append(sels, sel)
	}
	return &Transform{selectors: sels, requireAll: requireAll}, nil
}

// Apply projects the given record according to the configured selectors.  When
// no selectors are configured the input map is returned unchanged.
func (t *Transform) Apply(record map[string]any) (map[string]any, error) {
	if len(t.selectors) == 0 {
		return record, nil
	}
	out := make(map[string]any, len(t.selectors))
	for _, sel := range t.selectors {
		v, ok := record[sel.From]
		if !ok {
			if t.requireAll {
				return nil, fmt.Errorf("transform: field %q not found in record", sel.From)
			}
			continue
		}
		out[sel.To] = v
	}
	return out, nil
}

// Fields returns the list of source field names this transform will extract.
// An empty slice means all fields are passed through.
func (t *Transform) Fields() []string {
	if len(t.selectors) == 0 {
		return nil
	}
	names := make([]string, len(t.selectors))
	for i, s := range t.selectors {
		names[i] = s.From
	}
	return names
}
