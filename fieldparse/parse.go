package fieldparse

import (
	"flag"
	"fmt"
)

// RegisterFlags adds --parse-field and --parse-format flags to fs.
func RegisterFlags(fs *flag.FlagSet, field, format *string) {
	fs.StringVar(field, "parse-field", "",
		"parse this field's string value as a nested structured record and merge its keys")
	fs.StringVar(format, "parse-format", "auto",
		"format of the nested record: auto, json, or logfmt (default: auto)")
}

// FromFlags constructs a Parser from already-parsed flag values.
// Returns None when field is empty so callers can always call Apply safely.
func FromFlags(field, format string) (*Parser, error) {
	if field == "" {
		return None, nil
	}
	p, err := New(field, format)
	if err != nil {
		return nil, fmt.Errorf("fieldparse: %w", err)
	}
	return p, nil
}
