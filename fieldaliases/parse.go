package fieldaliases

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// RegisterFlags adds the --field-alias flag to the supplied flag set.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.StringArray(
		"field-alias",
		nil,
		"rename a field key before output: oldkey=newkey (repeatable)",
	)
}

// FromFlags reads the --field-alias values from fs and returns a Mapper.
// Returns None() when no aliases are configured.
func FromFlags(fs *pflag.FlagSet) (*Mapper, error) {
	specs, err := fs.GetStringArray("field-alias")
	if err != nil {
		return nil, fmt.Errorf("fieldaliases: %w", err)
	}
	var expanded []string
	for _, s := range specs {
		for _, part := range strings.Split(s, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				expanded = append(expanded, part)
			}
		}
	}
	if len(expanded) == 0 {
		return None(), nil
	}
	return New(expanded)
}
