// Package dedupe provides line deduplication for greplite output.
// It supports exact and field-based deduplication for structured log formats.
package dedupe

import "crypto/md5"

// Mode controls how deduplication is performed.
type Mode int

const (
	// ModeNone disables deduplication.
	ModeNone Mode = iota
	// ModeExact deduplicates identical raw lines.
	ModeExact
	// ModeField deduplicates based on a specific parsed field value.
	ModeField
)

// Deduper tracks seen lines and reports whether a line is a duplicate.
type Deduper struct {
	mode  Mode
	field string
	seen  map[[16]byte]struct{}
}

// New creates a Deduper with the given mode. For ModeField, field is the
// structured log field name to deduplicate on.
func New(mode Mode, field string) *Deduper {
	return &Deduper{
		mode:  mode,
		field: field,
		seen:  make(map[[16]byte]struct{}),
	}
}

// IsDuplicate returns true if the line (or its field value) has been seen
// before. It records the line so future identical lines are detected.
// fields may be nil when the input is not structured.
func (d *Deduper) IsDuplicate(raw string, fields map[string]string) bool {
	if d.mode == ModeNone {
		return false
	}

	var key string
	switch d.mode {
	case ModeField:
		if fields != nil {
			key = fields[d.field]
		} else {
			key = raw
		}
	default: // ModeExact
		key = raw
	}

	h := md5.Sum([]byte(key)) //nolint:gosec // not used for security
	if _, exists := d.seen[h]; exists {
		return true
	}
	d.seen[h] = struct{}{}
	return false
}

// Reset clears the set of seen lines.
func (d *Deduper) Reset() {
	d.seen = make(map[[16]byte]struct{})
}

// Count returns the number of unique lines recorded so far.
func (d *Deduper) Count() int {
	return len(d.seen)
}
