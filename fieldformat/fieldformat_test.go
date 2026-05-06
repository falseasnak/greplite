package fieldformat_test

import (
	"testing"

	"github.com/your-org/greplite/fieldformat"
)

func rec(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestNonePassesThrough(t *testing.T) {
	f := fieldformat.None()
	in := rec("x", 42)
	out := f.Apply(in)
	if out["x"] != 42 {
		t.Fatalf("expected 42, got %v", out["x"])
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := fieldformat.New("", "dst", "%d")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewEmptyFmtReturnsError(t *testing.T) {
	_, err := fieldformat.New("x", "dst", "")
	if err == nil {
		t.Fatal("expected error for empty format spec")
	}
}

func TestNewNoVerbReturnsError(t *testing.T) {
	_, err := fieldformat.New("x", "dst", "hello")
	if err == nil {
		t.Fatal("expected error when format spec has no verb")
	}
}

func TestApplyOverwritesSourceField(t *testing.T) {
	f, err := fieldformat.New("score", "", "%.2f")
	if err != nil {
		t.Fatal(err)
	}
	out := f.Apply(rec("score", 3.14159))
	if out["score"] != "3.14" {
		t.Fatalf("expected \"3.14\", got %q", out["score"])
	}
}

func TestApplyWritesToDestField(t *testing.T) {
	f, err := fieldformat.New("count", "count_fmt", "%05d")
	if err != nil {
		t.Fatal(err)
	}
	out := f.Apply(rec("count", 7))
	if out["count_fmt"] != "00007" {
		t.Fatalf("expected \"00007\", got %q", out["count_fmt"])
	}
	if out["count"] != 7 {
		t.Fatal("source field should be preserved")
	}
}

func TestApplyMissingFieldPassesThrough(t *testing.T) {
	f, err := fieldformat.New("missing", "", "%s")
	if err != nil {
		t.Fatal(err)
	}
	in := rec("other", "val")
	out := f.Apply(in)
	if _, exists := out["missing"]; exists {
		t.Fatal("missing field should not be created")
	}
}

func TestDoesNotMutateOriginal(t *testing.T) {
	f, err := fieldformat.New("n", "", "%d")
	if err != nil {
		t.Fatal(err)
	}
	in := rec("n", 1)
	_ = f.Apply(in)
	if in["n"] != 1 {
		t.Fatal("original record must not be mutated")
	}
}
