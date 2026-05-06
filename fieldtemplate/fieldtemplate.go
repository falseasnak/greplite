// Package fieldtemplate provides a computed field that is derived from
// other fields using a Go text/template expression.
package fieldtemplate

import (
	"bytes"
	"fmt"
	"text/template"
)

// Applier adds or overwrites a single field in a record using a
// rendered template.
type Applier struct {
	destField string
	tmpl      *template.Template
}

// None is a no-op Applier that leaves records unchanged.
func None() *Applier { return nil }

// New creates an Applier that writes the result of tmplSrc into
// destField for every record it processes.
//
// The template receives the record map as its dot value, so fields are
// accessible via {{.fieldName}}.
func New(destField, tmplSrc string) (*Applier, error) {
	if destField == "" {
		return nil, fmt.Errorf("fieldtemplate: destination field must not be empty")
	}
	t, err := template.New("fieldtemplate").Option("missingkey=zero").Parse(tmplSrc)
	if err != nil {
		return nil, fmt.Errorf("fieldtemplate: parse template: %w", err)
	}
	return &Applier{destField: destField, tmpl: t}, nil
}

// Apply renders the template against rec and stores the result in the
// destination field. The original map is not mutated; a shallow copy
// is returned.
//
// If a is nil the record is returned unchanged.
func (a *Applier) Apply(rec map[string]string) (map[string]string, error) {
	if a == nil {
		return rec, nil
	}
	var buf bytes.Buffer
	if err := a.tmpl.Execute(&buf, rec); err != nil {
		return rec, fmt.Errorf("fieldtemplate: execute: %w", err)
	}
	out := make(map[string]string, len(rec)+1)
	for k, v := range rec {
		out[k] = v
	}
	out[a.destField] = buf.String()
	return out, nil
}

// DestField returns the name of the field that will be written.
func (a *Applier) DestField() string {
	if a == nil {
		return ""
	}
	return a.destField
}
