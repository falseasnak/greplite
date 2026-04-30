package fieldmask

import (
	"testing"
)

func record() map[string]string {
	return map[string]string{
		"level": "info",
		"msg":   "hello",
		"svc":   "api",
		"ts":    "2024-01-01",
	}
}

func TestNonePassesThrough(t *testing.T) {
	m := None()
	r := record()
	out := m.Apply(r)
	if len(out) != len(r) {
		t.Fatalf("expected %d fields, got %d", len(r), len(out))
	}
}

func TestAllowKeepsOnlyListed(t *testing.T) {
	m := NewAllow([]string{"level", "msg"})
	out := m.Apply(record())
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
	if out["level"] != "info" || out["msg"] != "hello" {
		t.Fatalf("unexpected values: %v", out)
	}
	if _, ok := out["svc"]; ok {
		t.Fatal("svc should have been removed")
	}
}

func TestDenyRemovesListed(t *testing.T) {
	m := NewDeny([]string{"ts", "svc"})
	out := m.Apply(record())
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
	if _, ok := out["ts"]; ok {
		t.Fatal("ts should have been removed")
	}
	if _, ok := out["svc"]; ok {
		t.Fatal("svc should have been removed")
	}
}

func TestAllowEmptyFieldsReturnsEmpty(t *testing.T) {
	m := NewAllow([]string{"nonexistent"})
	out := m.Apply(record())
	if len(out) != 0 {
		t.Fatalf("expected 0 fields, got %d", len(out))
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	m := NewDeny([]string{"level"})
	r := record()
	m.Apply(r)
	if _, ok := r["level"]; !ok {
		t.Fatal("original record was mutated")
	}
}

func TestParseCSV(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{"level,msg", []string{"level", "msg"}},
		{" level , msg ", []string{"level", "msg"}},
		{"", []string{}},
		{"only", []string{"only"}},
	}
	for _, tc := range cases {
		got := ParseCSV(tc.input)
		if len(got) != len(tc.expected) {
			t.Errorf("ParseCSV(%q): got %v, want %v", tc.input, got, tc.expected)
			continue
		}
		for i := range tc.expected {
			if got[i] != tc.expected[i] {
				t.Errorf("ParseCSV(%q)[%d]: got %q, want %q", tc.input, i, got[i], tc.expected[i])
			}
		}
	}
}
