package fieldmap

import (
	"testing"
)

func TestNonePassesThrough(t *testing.T) {
	m := None()
	record := map[string]interface{}{"level": "info", "msg": "hello"}
	out := m.Apply(record)
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Errorf("expected unchanged record, got %v", out)
	}
}

func TestNewValidMapping(t *testing.T) {
	m, err := New([]string{"level=severity", "msg=message"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	record := map[string]interface{}{"level": "warn", "msg": "oops", "ts": "now"}
	out := m.Apply(record)
	if _, ok := out["level"]; ok {
		t.Error("expected 'level' to be renamed, but it still exists")
	}
	if out["severity"] != "warn" {
		t.Errorf("expected severity=warn, got %v", out["severity"])
	}
	if out["message"] != "oops" {
		t.Errorf("expected message=oops, got %v", out["message"])
	}
	if out["ts"] != "now" {
		t.Errorf("expected ts to pass through unchanged, got %v", out["ts"])
	}
}

func TestNewInvalidMapping(t *testing.T) {
	cases := []string{
		"noequalssign",
		"=newname",
		"oldname=",
		"",
	}
	for _, c := range cases {
		_, err := New([]string{c})
		if err == nil {
			t.Errorf("expected error for mapping %q, got nil", c)
		}
	}
}

func TestApplyNoOverlapWithUnmapped(t *testing.T) {
	m, err := New([]string{"a=b"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	record := map[string]interface{}{"a": 1, "c": 3}
	out := m.Apply(record)
	if out["b"] != 1 {
		t.Errorf("expected b=1, got %v", out["b"])
	}
	if out["c"] != 3 {
		t.Errorf("expected c=3, got %v", out["c"])
	}
	if _, ok := out["a"]; ok {
		t.Error("expected 'a' to be removed after rename")
	}
}

func TestMappingsReturnsCopy(t *testing.T) {
	m, _ := New([]string{"x=y"})
	copy1 := m.Mappings()
	copy1["x"] = "z"
	copy2 := m.Mappings()
	if copy2["x"] != "y" {
		t.Error("Mappings should return a copy, not a reference")
	}
}

func TestNoneReturnsSameRecordOnEmpty(t *testing.T) {
	m := None()
	record := map[string]interface{}{}
	out := m.Apply(record)
	if len(out) != 0 {
		t.Errorf("expected empty record, got %v", out)
	}
}
