package levelfilter

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// RegisterFlags adds the --level flag to the provided flag set.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.String("level", "", `minimum log level to display (trace|debug|info|warn|error|fatal)`)
}

// FromFlags constructs a Filter from the parsed flag set.
// Returns None() when the flag is empty.
func FromFlags(fs *pflag.FlagSet) (*Filter, error) {
	val, err := fs.GetString("level")
	if err != nil {
		return nil, fmt.Errorf("levelfilter: %w", err)
	}
	val = strings.TrimSpace(val)
	if val == "" {
		return None(), nil
	}
	f, err := New(val)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// KnownLevels returns the list of recognised level names in severity order.
func KnownLevels() []string {
	return []string{"trace", "debug", "info", "warn", "error", "fatal"}
}
