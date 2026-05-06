package fieldsplit

import (
	"testing"
)

func rec(pairs ...string) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestNonePassesThrough(t *testing.T) {
	r := rec("msg", "hello")
	got := None.Apply(r)
	if got["msg"] != "hello" {
		t.Fatalf("expected hello, got %v", got["msg"])
	}
}

func TestNewEmptySourceReturnsError(t *testing.T) {
	_, err := New("", []string{"a"}, ",")
	if err == nil {
		t.Fatal("expected error for empty source")
	}
}

func TestNewNoDestsReturnsError(t *testing.T) {
	_, err := New("src", nil, ",")
	if err == nil {
		t.Fatal("expected error for no dests")
	}
}

func TestNewEmptyDestElementReturnsError(t *testing.T) {
	_, err := New("src", []string{"a", ""}, ",")
	if err == nil {
		t.Fatal("expected error for empty dest element")
	}
}

func TestNewEmptySepReturnsError(t *testing.T) {
	_, err := New("src", []string{"a"}, "")
	if err == nil {
		t.Fatal("expected error for empty separator")
	}
}

func TestApplySplitsField(t *testing.T) {
	s, err := New("addr", []string{"host", "port"}, ":")
	if err != nil {
		t.Fatal(err)
	}
	out := s.Apply(rec("addr", "localhost:8080"))
	if out["host"] != "localhost" {
		t.Errorf("host: got %v", out["host"])
	}
	if out["port"] != "8080" {
		t.Errorf("port: got %v", out["port"])
	}
	if out["addr"] != "localhost:8080" {
		t.Error("source field should be preserved")
	}
}

func TestApplyMissingSourcePreservesRecord(t *testing.T) {
	s, _ := New("addr", []string{"host", "port"}, ":")
	out := s.Apply(rec("msg", "hello"))
	if _, ok := out["host"]; ok {
		t.Error("host should not be present when source is missing")
	}
}

func TestApplyFewerPartsThanDests(t *testing.T) {
	s, _ := New("tag", []string{"a", "b", "c"}, "-")
	out := s.Apply(rec("tag", "x-y"))
	if out["a"] != "x" || out["b"] != "y" {
		t.Errorf("unexpected values: %v", out)
	}
	if out["c"] != "" {
		t.Errorf("expected empty string for missing part, got %v", out["c"])
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	s, _ := New("kv", []string{"k", "v"}, "=")
	orig := rec("kv", "foo=bar")
	s.Apply(orig)
	if _, ok := orig["k"]; ok {
		t.Error("original record should not be mutated")
	}
}
