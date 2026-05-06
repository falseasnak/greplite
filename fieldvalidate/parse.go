package fieldvalidate

import (
	"flag"
	"fmt"
	"strings"
)

// RegisterFlags adds fieldvalidate flags to the given FlagSet.
func RegisterFlags(fs *flag.FlagSet) {
	fs.String("validate", "", "Comma-separated field validation rules (field:kind or field:regex:pattern).\nKinds: nonempty, number, bool, regex")
}

// FromFlags builds a Validator from the parsed flag set.
// Returns None() when no --validate flag is provided.
func FromFlags(fs *flag.FlagSet) (*Validator, error) {
	f := fs.Lookup("validate")
	if f == nil || f.Value.String() == "" {
		return None(), nil
	}
	return ParseCSV(f.Value.String())
}

// ParseCSV parses a comma-separated list of rule strings into a Validator.
// Rules using the regex kind may embed commas inside the pattern only if the
// whole rule is quoted; for typical use cases plain comma splitting is sufficient.
func ParseCSV(csv string) (*Validator, error) {
	if strings.TrimSpace(csv) == "" {
		return None(), nil
	}
	raw := strings.Split(csv, ",")
	rules := make([]Rule, 0, len(raw))
	for _, s := range raw {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		r, err := ParseRule(s)
		if err != nil {
			return nil, fmt.Errorf("fieldvalidate: ParseCSV: %w", err)
		}
		rules = append(rules, r)
	}
	return New(rules), nil
}
