package levelfilter

import (
	"testing"
)

func TestNoneAcceptsAll(t *testing.T) {
	f := None()
	cases := []map[string]string{
		{"level": "debug"},
		{"level": "error"},
		{},
	}
	for _, c := range cases {
		if !f.Allow(c) {
			t.Errorf("None() rejected record %v", c)
		}
	}
}

func TestNewUnknownLevel(t *testing.T) {
	_, err := New("verbose")
	if err == nil {
		t.Fatal("expected error for unknown level")
	}
}

func TestAllowMinLevelWarn(t *testing.T) {
	f, err := New("warn")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	allow := []map[string]string{
		{"level": "warn"},
		{"level": "warning"},
		{"level": "error"},
		{"level": "fatal"},
		{"level": "ERROR"}, // case-insensitive
	}
	deny := []map[string]string{
		{"level": "debug"},
		{"level": "info"},
		{"level": "trace"},
	}
	for _, c := range allow {
		if !f.Allow(c) {
			t.Errorf("expected Allow for %v", c)
		}
	}
	for _, c := range deny {
		if f.Allow(c) {
			t.Errorf("expected Deny for %v", c)
		}
	}
}

func TestAllowAlternativeFieldNames(t *testing.T) {
	f, _ := New("error")
	if !f.Allow(map[string]string{"lvl": "error"}) {
		t.Error("expected Allow for lvl=error")
	}
	if !f.Allow(map[string]string{"severity": "FATAL"}) {
		t.Error("expected Allow for severity=FATAL")
	}
	if f.Allow(map[string]string{"log_level": "info"}) {
		t.Error("expected Deny for log_level=info")
	}
}

func TestAllowNoLevelFieldPassesThrough(t *testing.T) {
	f, _ := New("error")
	if !f.Allow(map[string]string{"msg": "hello"}) {
		t.Error("record with no level field should pass through")
	}
}

func TestAllowUnrecognisedLevelValuePassesThrough(t *testing.T) {
	f, _ := New("error")
	if !f.Allow(map[string]string{"level": "critical"}) {
		t.Error("unrecognised level value should pass through")
	}
}

func TestCaseInsensitiveMinLevel(t *testing.T) {
	f, err := New("WARN")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Allow(map[string]string{"level": "error"}) {
		t.Error("expected Allow for error >= warn")
	}
}
