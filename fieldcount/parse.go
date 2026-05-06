package fieldcount

import (
	"fmt"
	"strconv"
	"strings"
)

// RegisterFlags adds --min-fields and --max-fields to any flag set that
// implements a StringVar-style interface. Callers typically pass *flag.FlagSet.
type flagSet interface {
	StringVar(p *string, name, value, usage string)
}

// Flags holds the raw string values captured from CLI flags.
type Flags struct {
	MinFields string
	MaxFields string
}

// RegisterFlags registers --min-fields and --max-fields on fs and stores the
// parsed strings in f.
func RegisterFlags(fs flagSet, f *Flags) {
	fs.StringVar(&f.MinFields, "min-fields", "", "only emit records with at least N parsed fields")
	fs.StringVar(&f.MaxFields, "max-fields", "", "only emit records with at most N parsed fields")
}

// FromFlags constructs a Filter from the values captured by RegisterFlags.
// Returns None() when neither flag is set.
func FromFlags(f Flags) (*Filter, error) {
	if f.MinFields == "" && f.MaxFields == "" {
		return None(), nil
	}

	min := 0
	max := -1

	if f.MinFields != "" {
		v, err := parseNonNeg(f.MinFields, "--min-fields")
		if err != nil {
			return nil, err
		}
		min = v
	}

	if f.MaxFields != "" {
		v, err := parseNonNeg(f.MaxFields, "--max-fields")
		if err != nil {
			return nil, err
		}
		max = v
	}

	return New(min, max)
}

func parseNonNeg(s, flag string) (int, error) {
	s = strings.TrimSpace(s)
	v, err := strconv.Atoi(s)
	if err != nil || v < 0 {
		return 0, fmt.Errorf("%s: expected non-negative integer, got %q", flag, s)
	}
	return v, nil
}
