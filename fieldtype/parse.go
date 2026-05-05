package fieldtype

import (
	"fmt"
	"strings"
)

// FieldSpec pairs a field name with an explicit Kind override.
// Example CLI representation: "latency:number" or just "latency" (auto).
type FieldSpec struct {
	Field string
	Kind  Kind
}

// ParseFieldSpec parses a single field-type specification of the form
// "fieldname" or "fieldname:kind".
func ParseFieldSpec(s string) (FieldSpec, error) {
	parts := strings.SplitN(s, ":", 2)
	field := strings.TrimSpace(parts[0])
	if field == "" {
		return FieldSpec{}, fmt.Errorf("fieldtype: empty field name in spec %q", s)
	}
	if len(parts) == 1 {
		return FieldSpec{Field: field, Kind: KindAuto}, nil
	}
	kind, err := ParseKind(parts[1])
	if err != nil {
		return FieldSpec{}, err
	}
	return FieldSpec{Field: field, Kind: kind}, nil
}

// ParseFieldSpecs parses a comma-separated list of field specs.
func ParseFieldSpecs(csv string) ([]FieldSpec, error) {
	if strings.TrimSpace(csv) == "" {
		return nil, nil
	}
	parts := strings.Split(csv, ",")
	specs := make([]FieldSpec, 0, len(parts))
	for _, p := range parts {
		spec, err := ParseFieldSpec(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		specs = append(specs, spec)
	}
	return specs, nil
}

// Registry maps field names to their configured Kind for fast lookup.
type Registry map[string]Kind

// NewRegistry builds a Registry from a slice of FieldSpecs.
func NewRegistry(specs []FieldSpec) Registry {
	r := make(Registry, len(specs))
	for _, s := range specs {
		r[s.Field] = s.Kind
	}
	return r
}

// Resolve returns a typed Value for the given field name and raw string,
// using the configured Kind if present, otherwise KindAuto.
func (r Registry) Resolve(field, raw string) Value {
	kind, ok := r[field]
	if !ok {
		kind = KindAuto
	}
	return Coerce(raw, kind)
}
