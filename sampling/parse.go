package sampling

import (
	"fmt"
	"strconv"
	"strings"
)

// Config holds the raw CLI values for sampling flags before they are
// resolved into a Sampler.
type Config struct {
	// Rate, if > 0, enables rate-based sampling (keep every Nth line).
	Rate int
	// Prob, if > 0, enables random sampling with this probability.
	Prob float64
	// Seed is used for the random sampler; 0 means use default seed.
	Seed int64
}

// FromFlags builds a Sampler from a Config. Returns NewNone() when no
// sampling is configured, an error when both Rate and Prob are set, or
// when the individual values are invalid.
func FromFlags(cfg Config) (*Sampler, error) {
	if cfg.Rate > 0 && cfg.Prob > 0 {
		return nil, fmt.Errorf("sampling: --sample-rate and --sample-prob are mutually exclusive")
	}
	if cfg.Rate > 0 {
		return NewRate(cfg.Rate)
	}
	if cfg.Prob > 0 {
		return NewRandom(cfg.Prob, cfg.Seed)
	}
	return NewNone(), nil
}

// ParseRate parses a string such as "10" into an integer rate.
func ParseRate(s string) (int, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("sampling: invalid rate %q: %w", s, err)
	}
	return n, nil
}

// ParseProb parses a string such as "0.25" into a float64 probability.
func ParseProb(s string) (float64, error) {
	s = strings.TrimSpace(s)
	p, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("sampling: invalid probability %q: %w", s, err)
	}
	return p, nil
}
