// Package fieldaliases provides a mapper that rewrites field keys in a parsed
// record according to a set of alias rules before the record reaches downstream
// pipeline stages.
package fieldaliases

import "fmt"

// Mapper rewrites field keys according to a fixed alias table.
type Mapper struct {
	aliases map[string]string // old key → new key
}

// None returns a no-op Mapper that leaves every record untouched.
func None() *Mapper { return &Mapper{} }

// New constructs a Mapper from the supplied alias pairs.
// Each entry in aliases must be of the form "oldkey=newkey".
func New(aliases []string) (*Mapper, error) {
	m := &Mapper{aliases: make(map[string]string, len(aliases))}
	for _, a := range aliases {
		old, neu, err := splitAlias(a)
		if err != nil {
			return nil, err
		}
		if _, dup := m.aliases[old]; dup {
			return nil, fmt.Errorf("fieldaliases: duplicate source field %q", old)
		}
		m.aliases[old] = neu
	}
	return m, nil
}

// Apply returns a new map with aliased keys renamed.
// Keys that have no alias rule are passed through unchanged.
// If two source keys would map to the same destination key the last one wins.
func (m *Mapper) Apply(record map[string]any) map[string]any {
	if len(m.aliases) == 0 {
		return record
	}
	out := make(map[string]any, len(record))
	for k, v := range record {
		if alias, ok := m.aliases[k]; ok {
			out[alias] = v
		} else {
			out[k] = v
		}
	}
	return out
}

// Aliases returns a copy of the current alias table (old → new).
func (m *Mapper) Aliases() map[string]string {
	cp := make(map[string]string, len(m.aliases))
	for k, v := range m.aliases {
		cp[k] = v
	}
	return cp
}

func splitAlias(s string) (old, neu string, err error) {
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			old, neu = s[:i], s[i+1:]
			if old == "" {
				return "", "", fmt.Errorf("fieldaliases: empty source field in %q", s)
			}
			if neu == "" {
				return "", "", fmt.Errorf("fieldaliases: empty destination field in %q", s)
			}
			return old, neu, nil
		}
	}
	return "", "", fmt.Errorf("fieldaliases: missing '=' in alias spec %q", s)
}
