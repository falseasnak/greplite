package fieldclip

import (
	"fmt"
	"strconv"
	"strings"
)

// RegisterFlags adds --clip-field and --clip-suffix to the provided flag set.
// flags must expose a StringArrayVar / StringVar interface compatible with
// the stdlib flag package or pflag.
func RegisterFlags(fs interface {
	StringArrayVar(p *[]string, name, value, usage string)
	StringVar(p *string, name, value, usage string)
}, fields *[]string, suffix *string) {
	fs.StringArrayVar(fields, "clip-field", nil,
		"clip field value to N runes: FIELD=N (repeatable)")
	fs.StringVar(suffix, "clip-suffix", "…",
		"suffix appended to clipped values")
}

// FromFlags constructs a Clipper from parsed CLI flags.
// specs is a slice of "FIELD=N" strings produced by --clip-field.
func FromFlags(specs []string, suffix string) (*Clipper, error) {
	if len(specs) == 0 {
		return None(), nil
	}
	fields := make(map[string]int, len(specs))
	for _, spec := range specs {
		field, n, err := parseSpec(spec)
		if err != nil {
			return nil, err
		}
		fields[field] = n
	}
	return New(fields, suffix)
}

// parseSpec parses a single "FIELD=N" clip specification.
func parseSpec(spec string) (string, int, error) {
	parts := strings.SplitN(spec, "=", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", 0, fmt.Errorf("fieldclip: invalid spec %q, expected FIELD=N", spec)
	}
	n, err := strconv.Atoi(parts[1])
	if err != nil || n < 1 {
		return "", 0, fmt.Errorf("fieldclip: N in %q must be a positive integer", spec)
	}
	return parts[0], n, nil
}
