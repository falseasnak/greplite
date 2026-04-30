package truncate

import (
	"fmt"
	"strconv"
)

// FromFlags constructs a Truncator from CLI flag values.
// maxLen is the raw string value of the --max-line-length flag (empty = disabled).
// A value of "0" is treated as disabled (no truncation).
// Negative values are rejected with an error.
func FromFlags(maxLen string) (*Truncator, error) {
	if maxLen == "" || maxLen == "0" {
		return None(), nil
	}
	n, err := strconv.Atoi(maxLen)
	if err != nil {
		return nil, fmt.Errorf("truncate: invalid --max-line-length %q: %w", maxLen, err)
	}
	if n < 0 {
		return nil, fmt.Errorf("truncate: --max-line-length must be >= 0, got %d", n)
	}
	return New(n), nil
}

// MustFromFlags is like FromFlags but panics on error.
// Intended for use in tests or cases where the flag value has already been validated.
func MustFromFlags(maxLen string) *Truncator {
	t, err := FromFlags(maxLen)
	if err != nil {
		panic(err)
	}
	return t
}
