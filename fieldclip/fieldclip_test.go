package fieldclip

import (
	"strings"
	"testing"
)

func rec(pairs ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestNonePassesThrough(t *testing.T) {
	c := None()
	in := rec("msg", "hello world")
	out := c.Apply(in)
	if out["msg"] != "hello world" {
		t.Fatalf("expected passthrough, got %v", out["msg"])
	}
}

func TestNewInvalidMaxRunes(t *testing.T) {
	_, err := New(map[string]int{"msg": 0}, "…")
	if err == nil {
		t.Fatal("expected error for maxRunes=0")
	}
}

func TestNewEmptyFields(t *testing.T) {
	_, err := New(map[string]int{}, "…")
	if err == nil {
		t.Fatal("expected error for empty fields map")
	}
}

func TestShortValueUnchanged(t *testing.T) {
	c, _ := New(map[string]int{"msg": 20}, "…")
	out := c.Apply(rec("msg", "hi"))
	if out["msg"] != "hi" {
		t.Fatalf("unexpected clip: %v", out["msg"])
	}
}

func TestLongValueClipped(t *testing.T) {
	c, _ := New(map[string]int{"msg": 5}, "…")
	out := c.Apply(rec("msg", "hello world"))
	got := out["msg"].(string)
	if got != "hello…" {
		t.Fatalf("expected \"hello…\", got %q", got)
	}
}

func TestNonStringFieldSkipped(t *testing.T) {
	c, _ := New(map[string]int{"count": 3}, "…")
	out := c.Apply(rec("count", 42))
	if out["count"] != 42 {
		t.Fatalf("expected non-string to pass through, got %v", out["count"])
	}
}

func TestMultibyteRunes(t *testing.T) {
	c, _ := New(map[string]int{"msg": 3}, "")
	// "日本語テスト" — 6 runes
	out := c.Apply(rec("msg", "日本語テスト"))
	got := out["msg"].(string)
	if got != "日本語" {
		t.Fatalf("expected \"日本語\", got %q", got)
	}
}

func TestUnrelatedFieldsPreserved(t *testing.T) {
	c, _ := New(map[string]int{"msg": 5}, "…")
	out := c.Apply(rec("msg", "hello world", "level", "info"))
	if out["level"] != "info" {
		t.Fatalf("unrelated field mutated: %v", out["level"])
	}
}

func TestFromFlagsNone(t *testing.T) {
	c, err := FromFlags(nil, "…")
	if err != nil || len(c.fields) != 0 {
		t.Fatalf("expected none clipper, got err=%v fields=%v", err, c.fields)
	}
}

func TestFromFlagsValid(t *testing.T) {
	c, err := FromFlags([]string{"msg=10", "body=50"}, "...")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.fields["msg"] != 10 || c.fields["body"] != 50 {
		t.Fatalf("unexpected field map: %v", c.fields)
	}
	if c.suffix != "..." {
		t.Fatalf("unexpected suffix: %q", c.suffix)
	}
}

func TestFromFlagsBadSpec(t *testing.T) {
	_, err := FromFlags([]string{"msgonly"}, "…")
	if err == nil || !strings.Contains(err.Error(), "invalid spec") {
		t.Fatalf("expected invalid spec error, got %v", err)
	}
}

func TestFromFlagsNonPositiveN(t *testing.T) {
	_, err := FromFlags([]string{"msg=0"}, "…")
	if err == nil {
		t.Fatal("expected error for N=0")
	}
}
