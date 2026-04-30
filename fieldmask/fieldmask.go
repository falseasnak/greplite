// Package fieldmask provides field inclusion/exclusion filtering for
// structured log records. It allows callers to keep only a named set of
// fields (allowlist) or to drop a named set of fields (denylist).
package fieldmask

import "strings"

// Mode controls whether the mask is an allowlist or a denylist.
type Mode int

const (
	// ModeNone disables field masking; all fields pass through unchanged.
	ModeNone Mode = iota
	// ModeAllow keeps only the listed fields.
	ModeAllow
	// ModeDeny drops the listed fields.
	ModeDeny
)

// Mask filters fields in a log record according to its Mode and field list.
type Mask struct {
	mode   Mode
	fields map[string]struct{}
}

// None returns a Mask that never modifies a record.
func None() *Mask {
	return &Mask{mode: ModeNone}
}

// NewAllow returns a Mask that retains only the given fields.
// Field names are case-sensitive.
func NewAllow(fields []string) *Mask {
	return &Mask{mode: ModeAllow, fields: toSet(fields)}
}

// NewDeny returns a Mask that removes the given fields from every record.
func NewDeny(fields []string) *Mask {
	return &Mask{mode: ModeDeny, fields: toSet(fields)}
}

// Apply returns a new map with the mask applied to fields.
// The original map is never modified.
func (m *Mask) Apply(record map[string]string) map[string]string {
	if m.mode == ModeNone || len(m.fields) == 0 {
		return record
	}
	out := make(map[string]string, len(record))
	for k, v := range record {
		switch m.mode {
		case ModeAllow:
			if _, ok := m.fields[k]; ok {
				out[k] = v
			}
		case ModeDeny:
			if _, ok := m.fields[k]; !ok {
				out[k] = v
			}
		}
	}
	return out
}

// ParseCSV parses a comma-separated list of field names, trimming whitespace.
func ParseCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

func toSet(fields []string) map[string]struct{} {
	s := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		s[f] = struct{}{}
	}
	return s
}
