// Package truncate provides line-length truncation for output lines.
// It supports truncating long lines to a configurable maximum byte length,
// optionally appending an ellipsis indicator.
package truncate

import "unicode/utf8"

const defaultEllipsis = "..."

// Truncator truncates strings that exceed a maximum length.
type Truncator struct {
	maxLen   int
	ellipsis string
	enabled  bool
}

// New creates a Truncator that truncates lines longer than maxLen runes.
// If maxLen <= 0, truncation is disabled.
func New(maxLen int) *Truncator {
	if maxLen <= 0 {
		return &Truncator{enabled: false}
	}
	return &Truncator{
		maxLen:   maxLen,
		ellipsis: defaultEllipsis,
		enabled:  true,
	}
}

// None returns a no-op Truncator that never truncates.
func None() *Truncator {
	return &Truncator{enabled: false}
}

// Apply truncates s if it exceeds the configured maximum rune count.
// Returns the (possibly truncated) string and whether truncation occurred.
func (t *Truncator) Apply(s string) (string, bool) {
	if !t.enabled {
		return s, false
	}
	count := utf8.RuneCountInString(s)
	if count <= t.maxLen {
		return s, false
	}
	// Find byte offset for maxLen runes minus ellipsis rune count.
	ellipsisCnt := utf8.RuneCountInString(t.ellipsis)
	cutAt := t.maxLen - ellipsisCnt
	if cutAt < 0 {
		cutAt = 0
	}
	byteIdx := 0
	for i := 0; i < cutAt; i++ {
		_, size := utf8.DecodeRuneInString(s[byteIdx:])
		byteIdx += size
	}
	return s[:byteIdx] + t.ellipsis, true
}

// Enabled reports whether truncation is active.
func (t *Truncator) Enabled() bool {
	return t.enabled
}

// MaxLen returns the configured maximum rune length (0 if disabled).
func (t *Truncator) MaxLen() int {
	return t.maxLen
}
