// Package fieldredact replaces sensitive field values with a redaction
// placeholder before lines are emitted. It is useful for masking tokens,
// passwords, or PII fields in log output.
package fieldredact

import "strings"

const defaultPlaceholder = "[REDACTED]"

// Redactor replaces the values of named fields with a fixed placeholder.
type Redactor struct {
	fields      map[string]struct{}
	placeholder string
}

// None returns a no-op Redactor that passes all fields through unchanged.
func None() *Redactor {
	return &Redactor{}
}

// New creates a Redactor that masks the given field names.
// An optional placeholder may be supplied; if empty the default is used.
func New(fields []string, placeholder string) *Redactor {
	if placeholder == "" {
		placeholder = defaultPlaceholder
	}
	set := make(map[string]struct{}, len(fields))
	for _, f := range fields {
		set[strings.TrimSpace(f)] = struct{}{}
	}
	return &Redactor{fields: set, placeholder: placeholder}
}

// Apply returns a shallow copy of fields with sensitive values replaced.
// If the Redactor has no configured fields the original map is returned as-is.
func (r *Redactor) Apply(record map[string]any) map[string]any {
	if len(r.fields) == 0 {
		return record
	}
	out := make(map[string]any, len(record))
	for k, v := range record {
		if _, masked := r.fields[k]; masked {
			out[k] = r.placeholder
		} else {
			out[k] = v
		}
	}
	return out
}

// Fields returns the set of field names that will be redacted.
func (r *Redactor) Fields() []string {
	result := make([]string, 0, len(r.fields))
	for f := range r.fields {
		result = append(result, f)
	}
	return result
}
