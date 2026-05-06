package fieldsplit

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
)

// RegisterFlags adds fieldsplit flags to the provided FlagSet.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.String("split-field", "", "source field to split")
	fs.StringSlice("split-into", nil, "comma-separated destination field names")
	fs.String("split-sep", ",", "separator used to split the source field")
}

// FromFlags builds a Splitter from the registered flags.
// Returns None if --split-field is not set.
func FromFlags(fs *pflag.FlagSet) (*Splitter, error) {
	src, err := fs.GetString("split-field")
	if err != nil {
		return nil, err
	}
	if src == "" {
		return None, nil
	}

	dests, err := fs.GetStringSlice("split-into")
	if err != nil {
		return nil, err
	}
	// GetStringSlice may return a single element that is itself comma-separated.
	var flat []string
	for _, d := range dests {
		for _, part := range strings.Split(d, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				flat = append(flat, part)
			}
		}
	}
	if len(flat) == 0 {
		return nil, fmt.Errorf("fieldsplit: --split-into must specify at least one destination field")
	}

	sep, err := fs.GetString("split-sep")
	if err != nil {
		return nil, err
	}

	return New(src, flat, sep)
}
