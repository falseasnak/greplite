package timefilter

import (
	"testing"
	"time"
)

func mustTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestNoneAcceptsAll(t *testing.T) {
	f := None()
	fields := map[string]string{"ts": "2024-01-01T00:00:00Z"}
	if !f.Match(fields, "ts") {
		t.Fatal("None filter should accept every line")
	}
}

func TestAfterOnly(t *testing.T) {
	f, _ := New(mustTime("2024-06-01T00:00:00Z"), time.Time{})
	tests := []struct {
		ts   string
		want bool
	}{
		{"2024-05-31T23:59:59Z", false},
		{"2024-06-01T00:00:00Z", true},
		{"2024-07-01T00:00:00Z", true},
	}
	for _, tc := range tests {
		fields := map[string]string{"ts": tc.ts}
		if got := f.Match(fields, "ts"); got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.ts, got, tc.want)
		}
	}
}

func TestBeforeOnly(t *testing.T) {
	f, _ := New(time.Time{}, mustTime("2024-06-01T00:00:00Z"))
	tests := []struct {
		ts   string
		want bool
	}{
		{"2024-05-31T23:59:59Z", true},
		{"2024-06-01T00:00:00Z", false},
		{"2024-07-01T00:00:00Z", false},
	}
	for _, tc := range tests {
		fields := map[string]string{"ts": tc.ts}
		if got := f.Match(fields, "ts"); got != tc.want {
			t.Errorf("Match(%q) = %v, want %v", tc.ts, got, tc.want)
		}
	}
}

func TestWindow(t *testing.T) {
	f, _ := New(mustTime("2024-06-01T00:00:00Z"), mustTime("2024-06-02T00:00:00Z"))
	if !f.Match(map[string]string{"ts": "2024-06-01T12:00:00Z"}, "ts") {
		t.Error("timestamp inside window should match")
	}
	if f.Match(map[string]string{"ts": "2024-05-31T00:00:00Z"}, "ts") {
		t.Error("timestamp before window should not match")
	}
	if f.Match(map[string]string{"ts": "2024-06-02T00:00:00Z"}, "ts") {
		t.Error("timestamp equal to before bound should not match")
	}
}

func TestInvalidRange(t *testing.T) {
	_, err := New(mustTime("2024-06-02T00:00:00Z"), mustTime("2024-06-01T00:00:00Z"))
	if err == nil {
		t.Fatal("expected error for inverted range")
	}
}

func TestMissingFieldKeepsLine(t *testing.T) {
	f, _ := New(mustTime("2024-06-01T00:00:00Z"), mustTime("2024-06-02T00:00:00Z"))
	if !f.Match(map[string]string{"msg": "hello"}, "ts") {
		t.Error("missing timestamp field should keep the line")
	}
}

func TestUnparsableTimestampKeepsLine(t *testing.T) {
	f, _ := New(mustTime("2024-06-01T00:00:00Z"), mustTime("2024-06-02T00:00:00Z"))
	if !f.Match(map[string]string{"ts": "not-a-date"}, "ts") {
		t.Error("unparsable timestamp should keep the line")
	}
}

func TestAlternativeLayouts(t *testing.T) {
	f, _ := New(mustTime("2024-06-01T00:00:00Z"), mustTime("2024-06-02T00:00:00Z"))
	layouts := []string{
		"2024-06-01T10:00:00.000Z",
		"2024-06-01T10:00:00",
		"2024-06-01 10:00:00",
	}
	for _, ts := range layouts {
		if !f.Match(map[string]string{"ts": ts}, "ts") {
			t.Errorf("layout %q should be accepted inside window", ts)
		}
	}
}
