package timefilter

import (
	"fmt"
	"time"
)

// humanLayouts are the formats accepted on the CLI in addition to RFC 3339.
var humanLayouts = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

// ParseTime parses a user-supplied timestamp string using a set of
// human-friendly layouts. It returns an error describing the expected
// formats when no layout matches.
func ParseTime(s string) (time.Time, error) {
	for _, l := range humanLayouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf(
		"cannot parse %q — accepted formats: RFC3339, 2006-01-02T15:04:05, 2006-01-02 15:04:05, 2006-01-02", s,
	)
}

// FromFlags constructs a Filter from the raw string values of --after and
// --before CLI flags. Empty strings mean unbounded on that side.
func FromFlags(afterStr, beforeStr string) (*Filter, error) {
	var after, before time.Time
	var err error

	if afterStr != "" {
		after, err = ParseTime(afterStr)
		if err != nil {
			return nil, fmt.Errorf("--after: %w", err)
		}
	}
	if beforeStr != "" {
		before, err = ParseTime(beforeStr)
		if err != nil {
			return nil, fmt.Errorf("--before: %w", err)
		}
	}
	return New(after, before)
}
