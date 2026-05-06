package fieldgroup_test

import (
	"testing"

	"github.com/user/greplite/fieldgroup"
)

func rec(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestNonePassesThrough(t *testing.T) {
	g := fieldgroup.None()
	in := rec("a", "1", "b", "2")
	out := g.Apply(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(out))
	}
}

func TestNewEmptyDestReturnsError(t *testing.T) {
	_, err := fieldgroup.New("", []string{"a"}, "-")
	if err == nil {
		t.Fatal("expected error for empty dest")
	}
}

func TestNewNoSourcesReturnsError(t *testing.T) {
	_, err := fieldgroup.New("out", nil, "-")
	if err == nil {
		t.Fatal("expected error for nil sources")
	}
}

func TestNewEmptySourceElementReturnsError(t *testing.T) {
	_, err := fieldgroup.New("out", []string{"a", ""}, "-")
	if err == nil {
		t.Fatal("expected error for empty source element")
	}
}

func TestApplyConcatenatesFields(t *testing.T) {
	g, err := fieldgroup.New("full", []string{"first", "last"}, " ")
	if err != nil {
		t.Fatal(err)
	}
	out := g.Apply(rec("first", "Jane", "last", "Doe"))
	if got := out["full"]; got != "Jane Doe" {
		t.Fatalf("expected 'Jane Doe', got %q", got)
	}
}

func TestApplyMissingSourceTreatedAsEmpty(t *testing.T) {
	g, _ := fieldgroup.New("full", []string{"first", "last"}, "-")
	out := g.Apply(rec("first", "Jane"))
	if got := out["full"]; got != "Jane-" {
		t.Fatalf("expected 'Jane-', got %q", got)
	}
}

func TestApplyDoesNotMutateOriginal(t *testing.T) {
	g, _ := fieldgroup.New("combo", []string{"x", "y"}, "+")
	in := rec("x", "1", "y", "2")
	g.Apply(in)
	if _, ok := in["combo"]; ok {
		t.Fatal("original record was mutated")
	}
}

func TestApplyNumericFields(t *testing.T) {
	g, _ := fieldgroup.New("label", []string{"host", "port"}, ":")
	out := g.Apply(rec("host", "localhost", "port", 8080))
	if got := out["label"]; got != "localhost:8080" {
		t.Fatalf("expected 'localhost:8080', got %q", got)
	}
}

func TestApplyOverwritesExistingDest(t *testing.T) {
	g, _ := fieldgroup.New("full", []string{"a", "b"}, "|")
	out := g.Apply(rec("a", "x", "b", "y", "full", "old"))
	if got := out["full"]; got != "x|y" {
		t.Fatalf("expected 'x|y', got %q", got)
	}
}
