package aggregate

import (
	"fmt"
	"strings"
)

// Config holds the parsed aggregation options supplied via CLI flags.
type Config struct {
	// Field is the structured-log field name to group by.
	Field string
	// TopN, when > 0, limits output to the N most frequent values.
	TopN int
}

// FromFlags validates and constructs a Config from raw flag values.
// field must be non-empty. topN must be >= 0 (0 means unlimited).
func FromFlags(field string, topN int) (*Config, error) {
	field = strings.TrimSpace(field)
	if field == "" {
		return nil, fmt.Errorf("aggregate: --agg-field must not be empty")
	}
	if topN < 0 {
		return nil, fmt.Errorf("aggregate: --agg-top must be >= 0, got %d", topN)
	}
	return &Config{Field: field, TopN: topN}, nil
}

// Apply trims the results slice to at most TopN entries when TopN > 0.
func (cfg *Config) Apply(results []Result) []Result {
	if cfg.TopN > 0 && len(results) > cfg.TopN {
		return results[:cfg.TopN]
	}
	return results
}
