package ratelimit

import (
	"fmt"
	"strconv"
	"strings"
)

// FromFlags constructs a Limiter from the CLI flag values.
// rateStr is the raw --rate-limit flag value (e.g. "100" or "50/drop").
// An empty string returns a no-op limiter.
func FromFlags(rateStr string) (*Limiter, error) {
	if rateStr == "" {
		return None(), nil
	}
	drop := false
	parts := strings.SplitN(rateStr, "/", 2)
	if len(parts) == 2 {
		switch strings.ToLower(strings.TrimSpace(parts[1])) {
		case "drop":
			drop = true
		case "block", "":
			drop = false
		default:
			return nil, fmt.Errorf("ratelimit: unknown mode %q (want 'drop' or 'block')", parts[1])
		}
	}
	n, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil || n <= 0 {
		return nil, fmt.Errorf("ratelimit: invalid rate %q: must be a positive integer", parts[0])
	}
	return New(n, drop), nil
}
