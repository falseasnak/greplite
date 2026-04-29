// Package fieldmap provides field renaming and aliasing for structured log output.
// It allows users to rename fields in parsed log records before they are
// formatted or filtered, enabling normalization across different log sources.
package fieldmap

import (
	"fmt"
	"strings"
)

// Mapper renames fields in a log record according to a configured mapping.
type Mapper struct {
	// mappings is a map from original field name to new field name.
	mappings map[string]string
}

// New creates a Mapper from a slice of "old=new" mapping strings.
// Returns an error if any mapping is malformed.
func New(specs []string) (*Mapper, error) {
	mappings := make(map[string]string, len(specs))
	for _, spec := range specs {
		parts := strings.SplitN(spec, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("fieldmap: invalid mapping %q: expected old=new", spec)
		}
		mappings[parts[0]] = parts[1]
	}
	return &Mapper{mappings: mappings}, nil
}

// None returns a no-op Mapper that does not rename any fields.
func None() *Mapper {
	return &Mapper{mappings: map[string]string{}}
}

// Apply returns a new record with fields renamed according to the mapping.
// Fields not present in the mapping are passed through unchanged.
// If the mapper has no mappings, the original record is returned as-is.
func (m *Mapper) Apply(record map[string]interface{}) map[string]interface{} {
	if len(m.mappings) == 0 {
		return record
	}
	out := make(map[string]interface{}, len(record))
	for k, v := range record {
		if newKey, ok := m.mappings[k]; ok {
			out[newKey] = v
		} else {
			out[k] = v
		}
	}
	return out
}

// Mappings returns a copy of the internal field rename map.
func (m *Mapper) Mappings() map[string]string {
	copy := make(map[string]string, len(m.mappings))
	for k, v := range m.mappings {
		copy[k] = v
	}
	return copy
}
