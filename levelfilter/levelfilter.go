// Package levelfilter provides filtering of structured log lines by log level.
// It supports common level field names and standard severity values.
package levelfilter

import (
	"strings"
)

// Order maps level names to numeric severity (higher = more severe).
var Order = map[string]int{
	"trace": 0,
	"debug": 1,
	"info":  2,
	"warn":  3,
	"warning": 3,
	"error": 4,
	"err":   4,
	"fatal": 5,
	"panic": 5,
}

// candidateFields are common field names that carry log level.
var candidateFields = []string{"level", "lvl", "severity", "log_level"}

// Filter accepts or rejects log records based on minimum severity level.
type Filter struct {
	minSeverity int
	disabled    bool
}

// None returns a Filter that accepts every record.
func None() *Filter {
	return &Filter{disabled: true}
}

// New returns a Filter that accepts records at or above minLevel severity.
// Returns an error if minLevel is not a recognised level name.
func New(minLevel string) (*Filter, error) {
	norm := strings.ToLower(strings.TrimSpace(minLevel))
	sev, ok := Order[norm]
	if !ok {
		return nil, &UnknownLevelError{Level: minLevel}
	}
	return &Filter{minSeverity: sev}, nil
}

// Allow returns true when the record's level meets the minimum severity.
// fields is the parsed key-value map from a structured log line; raw is the
// original line used as fallback when no level field is found.
func (f *Filter) Allow(fields map[string]string) bool {
	if f.disabled {
		return true
	}
	for _, key := range candidateFields {
		val, ok := fields[key]
		if !ok {
			continue
		}
		norm := strings.ToLower(strings.TrimSpace(val))
		sev, known := Order[norm]
		if !known {
			// Field present but unrecognised value — let it through.
			return true
		}
		return sev >= f.minSeverity
	}
	// No level field found — let the record through.
	return true
}

// UnknownLevelError is returned when an unrecognised level name is supplied.
type UnknownLevelError struct {
	Level string
}

func (e *UnknownLevelError) Error() string {
	return "levelfilter: unknown log level: " + e.Level
}
