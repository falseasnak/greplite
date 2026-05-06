package fieldexist

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// RegisterFlags adds --require-fields and --exclude-fields flags to the
// provided flag set.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.String("require-fields", "", "comma-separated list of fields that must be present in each record")
	fs.String("exclude-fields", "", "comma-separated list of fields whose presence causes a record to be dropped")
}

// FromFlags constructs a Filter from the registered flag values.
// Returns None() when neither flag is set.
func FromFlags(fs *pflag.FlagSet) (*Filter, error) {
	requireRaw, err := fs.GetString("require-fields")
	if err != nil {
		return nil, fmt.Errorf("fieldexist: %w", err)
	}
	excludeRaw, err := fs.GetString("exclude-fields")
	if err != nil {
		return nil, fmt.Errorf("fieldexist: %w", err)
	}

	if strings.TrimSpace(requireRaw) == "" && strings.TrimSpace(excludeRaw) == "" {
		return None(), nil
	}

	return New(splitCSV(requireRaw), splitCSV(excludeRaw))
}

func splitCSV(s string) []string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	out := parts[:0]
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
