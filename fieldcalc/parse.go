package fieldcalc

import (
	"fmt"

	"github.com/spf13/pflag"
)

// RegisterFlags adds the --calc flag to the provided flag set.
func RegisterFlags(fs *pflag.FlagSet) {
	fs.String("calc", "", `compute a new field via arithmetic, e.g. --calc "rate=bytes/secs"`)
}

// FromFlags constructs a Calc (or None) from parsed flags.
// Returns an error if the expression is provided but malformed.
func FromFlags(fs *pflag.FlagSet) (*Calc, error) {
	expr, err := fs.GetString("calc")
	if err != nil {
		return nil, fmt.Errorf("fieldcalc: flag error: %w", err)
	}
	if expr == "" {
		return None(), nil
	}
	return New(expr)
}
