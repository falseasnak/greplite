package fieldgroup

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// RegisterFlags adds fieldgroup CLI flags to fs.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.String("group-dest", "", "destination field name for grouped value")
	fs.StringSlice("group-src", nil, "comma-separated source fields to group (ordered)")
	fs.String("group-sep", " ", "separator used when joining grouped fields")
}

// FromFlags constructs a Grouper from parsed flags.
// Returns None() when --group-dest is not set.
func FromFlags(fs *pflag.FlagSet) (*Grouper, error) {
	dest, err := fs.GetString("group-dest")
	if err != nil {
		return nil, err
	}
	if dest == "" {
		return None(), nil
	}

	srcs, err := fs.GetStringSlice("group-src")
	if err != nil {
		return nil, err
	}
	// flatten in case user passed comma-separated values as a single string
	var sources []string
	for _, s := range srcs {
		for _, part := range strings.Split(s, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				sources = append(sources, part)
			}
		}
	}
	if len(sources) == 0 {
		return nil, fmt.Errorf("fieldgroup: --group-src must list at least one source field")
	}

	sep, err := fs.GetString("group-sep")
	if err != nil {
		return nil, err
	}
	return New(dest, sources, sep)
}
