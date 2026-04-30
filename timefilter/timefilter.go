// Package timefilter provides time-range filtering for structured log lines.
// It extracts a timestamp field from parsed log fields and checks whether
// the timestamp falls within an optional [After, Before) window.
package timefilter

import (
	"fmt"
	"time"
)

// Filter holds an optional time window. A nil bound means unbounded.
type Filter struct {
	after  *time.Time
	before *time.Time
}

// None returns a Filter that accepts every line.
func None() *Filter { return &Filter{} }

// New returns a Filter that accepts lines whose timestamp field value
// falls within [after, before). Either bound may be zero to indicate
// unbounded on that side.
func New(after, before time.Time) (*Filter, error) {
	f := &Filter{}
	if !after.IsZero() {
		t := after
		f.after = &t
	}
	if !before.IsZero() {
		t := before
		f.before = &t
	}
	if f.after != nil && f.before != nil && !f.before.After(*f.after) {
		return nil, fmt.Errorf("timefilter: --before must be after --after")
	}
	return f, nil
}

// Match returns true when fields contains a recognisable timestamp under
// tsField and that timestamp satisfies the window, or when no window is
// configured. If the field is absent or unparseable the line is kept so
// that non-timestamped lines are never silently dropped.
func (f *Filter) Match(fields map[string]string, tsField string) bool {
	if f.after == nil && f.before == nil {
		return true
	}
	raw, ok := fields[tsField]
	if !ok {
		return true
	}
	t, err := parseTime(raw)
	if err != nil {
		return true
	}
	if f.after != nil && t.Before(*f.after) {
		return false
	}
	if f.before != nil && !t.Before(*f.before) {
		return false
	}
	return true
}

// common timestamp layouts tried in order.
var layouts = []string{
	time.RFC3339Nano,
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
}

func parseTime(s string) (time.Time, error) {
	for _, l := range layouts {
		if t, err := time.Parse(l, s); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("timefilter: cannot parse %q", s)
}
