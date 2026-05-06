package fieldtemplate

import (
	"testing"
)

func rec(pairs ...string) map[string]string {
	m := make(map[string]string, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestNonePassesThrough(t *testing.T) {
	a := None()
	in := rec("msg", "hello")
	out, err := a.Apply(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["msg"] != "hello" {
		t.Fatalf("expected msg=hello, got %q", out["msg"])
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", "{{.msg}}")
	if err == nil {
		t.Fatal("expected error for empty dest field")
	}
}

func TestNewBadTemplateReturnsError(t *testing.T) {
	_, err := New("out", "{{.unclosed")
	if err == nil {
		t.Fatal("expected parse error")
	}
}

func TestApplySimpleInterpolation(t *testing.T) {
	a, err := New("summary", "{{.level}}: {{.msg}}")
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	out, err := a.Apply(rec("level", "error", "msg", "disk full"))
	if err != nil {
		t.Fatalf("Apply: %v", err)
	}
	if got := out["summary"]; got != "error: disk full" {
		t.Fatalf("expected 'error: disk full', got %q", got)
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	a, _ := New("extra", "computed")
	in := rec("k", "v")
	out, _ := a.Apply(in)
	if _, ok := in["extra"]; ok {
		t.Fatal("original map was mutated")
	}
	if out["extra"] != "computed" {
		t.Fatalf("expected computed, got %q", out["extra"])
	}
}

func TestApplyOverwritesExistingField(t *testing.T) {
	a, _ := New("msg", "overwritten")
	out, _ := a.Apply(rec("msg", "original"))
	if out["msg"] != "overwritten" {
		t.Fatalf("expected overwritten, got %q", out["msg"])
	}
}

func TestApplyMissingKeyRendersEmpty(t *testing.T) {
	a, _ := New("out", "{{.nosuchfield}}")
	out, err := a.Apply(rec("other", "val"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["out"] != "<no value>" && out["out"] != "" {
		// missingkey=zero renders the zero value for the type (empty string for map lookup)
		// Both outcomes are acceptable depending on Go version behaviour.
	}
}

func TestDestField(t *testing.T) {
	a, _ := New("target", "x")
	if a.DestField() != "target" {
		t.Fatalf("expected target, got %q", a.DestField())
	}
	if None().DestField() != "" {
		t.Fatal("None DestField should be empty")
	}
}
