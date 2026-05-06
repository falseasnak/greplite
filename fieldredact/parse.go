package fieldredact

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// RegisterFlags adds the --redact-fields and --redact-placeholder flags to fs.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.String("redact-fields", "", "comma-separated list of field names whose values will be redacted")
	fs.String("redact-placeholder", defaultPlaceholder, "replacement text used for redacted values")
}

// FromFlags builds a Redactor from the flags registered by RegisterFlags.
// It returns None() when no fields are specified.
func FromFlags(fs *pflag.FlagSet) (*Redactor, error) {
	raw, err := fs.GetString("redact-fields")
	if err != nil {
		return nil, fmt.Errorf("fieldredact: %w", err)
	}
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return None(), nil
	}

	placeholder, err := fs.GetString("redact-placeholder")
	if err != nil {
		return nil, fmt.Errorf("fieldredact: %w", err)
	}

	parts := strings.Split(raw, ",")
	fields := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			return nil, fmt.Errorf("fieldredact: empty field name in --redact-fields %q", raw)
		}
		fields = append(fields, p)
	}
	return New(fields, placeholder), nil
}
