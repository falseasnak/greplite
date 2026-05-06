// Package fieldparse provides a pipeline stage that parses a raw string field
// into a structured sub-record and merges the resulting key-value pairs into
// the parent record. Supported inner formats are "json" and "logfmt".
package fieldparse

import (
	"fmt"

	"github.com/yourorg/greplite/parser"
)

// Parser extracts a string field from a record, parses it as a nested
// structured format, and promotes the inner fields into the record.
type Parser struct {
	field  string
	format string // "auto", "json", "logfmt"
}

// None is a no-op Parser that passes every record through unchanged.
var None = &Parser{}

// New creates a Parser that reads field and parses its value using format.
// format must be one of "auto", "json", or "logfmt".
func New(field, format string) (*Parser, error) {
	if field == "" {
		return nil, fmt.Errorf("fieldparse: field name must not be empty")
	}
	switch format {
	case "auto", "json", "logfmt":
	default:
		return nil, fmt.Errorf("fieldparse: unknown format %q (want auto, json, or logfmt)", format)
	}
	return &Parser{field: field, format: format}, nil
}

// Apply parses the target field in rec and merges inner fields into a copy of
// rec. The original target field is preserved. If the field is absent or
// cannot be parsed the original record is returned unchanged.
func (p *Parser) Apply(rec map[string]any) map[string]any {
	if p.field == "" {
		return rec
	}
	raw, ok := rec[p.field]
	if !ok {
		return rec
	}
	s, ok := raw.(string)
	if !ok || s == "" {
		return rec
	}

	inner, err := parser.Auto([]byte(s))
	if err != nil || len(inner) == 0 {
		return rec
	}

	out := make(map[string]any, len(rec)+len(inner))
	for k, v := range rec {
		out[k] = v
	}
	for k, v := range inner {
		if _, exists := out[k]; !exists {
			out[k] = v
		}
	}
	return out
}
