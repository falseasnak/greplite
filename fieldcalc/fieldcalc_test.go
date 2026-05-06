package fieldcalc

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
	c := None()
	in := rec("a", "1")
	out := c.Apply(in)
	if out["a"] != "1" {
		t.Fatalf("expected passthrough, got %v", out)
	}
}

func TestAddition(t *testing.T) {
	c, err := New("total=x+y")
	if err != nil {
		t.Fatal(err)
	}
	out := c.Apply(rec("x", "3", "y", "4"))
	if out["total"] != "7" {
		t.Fatalf("expected 7, got %q", out["total"])
	}
}

func TestSubtraction(t *testing.T) {
	c, err := New("diff=a-b")
	if err != nil {
		t.Fatal(err)
	}
	out := c.Apply(rec("a", "10", "b", "3"))
	if out["diff"] != "7" {
		t.Fatalf("expected 7, got %q", out["diff"])
	}
}

func TestMultiplication(t *testing.T) {
	c, err := New("area=w*h")
	if err != nil {
		t.Fatal(err)
	}
	out := c.Apply(rec("w", "6", "h", "7"))
	if out["area"] != "42" {
		t.Fatalf("expected 42, got %q", out["area"])
	}
}

func TestDivision(t *testing.T) {
	c, err := New("rate=bytes/secs")
	if err != nil {
		t.Fatal(err)
	}
	out := c.Apply(rec("bytes", "100", "secs", "4"))
	if out["rate"] != "25" {
		t.Fatalf("expected 25, got %q", out["rate"])
	}
}

func TestDivisionByZeroPassesThrough(t *testing.T) {
	c, err := New("rate=bytes/secs")
	if err != nil {
		t.Fatal(err)
	}
	in := rec("bytes", "100", "secs", "0")
	out := c.Apply(in)
	if _, ok := out["rate"]; ok {
		t.Fatal("expected no rate field on division by zero")
	}
}

func TestMissingOperandPassesThrough(t *testing.T) {
	c, err := New("total=x+y")
	if err != nil {
		t.Fatal(err)
	}
	in := rec("x", "5")
	out := c.Apply(in)
	if _, ok := out["total"]; ok {
		t.Fatal("expected no total field when operand missing")
	}
}

func TestNonNumericOperandPassesThrough(t *testing.T) {
	c, err := New("total=x+y")
	if err != nil {
		t.Fatal(err)
	}
	out := c.Apply(rec("x", "foo", "y", "2"))
	if _, ok := out["total"]; ok {
		t.Fatal("expected no total field for non-numeric operand")
	}
}

func TestNewInvalidExprReturnsError(t *testing.T) {
	_, err := New("nodestination")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestNewEmptyDestReturnsError(t *testing.T) {
	_, err := New("=a+b")
	if err == nil {
		t.Fatal("expected error for empty destination")
	}
}

func TestDoesNotMutateOriginal(t *testing.T) {
	c, _ := New("sum=a+b")
	in := rec("a", "1", "b", "2")
	c.Apply(in)
	if _, ok := in["sum"]; ok {
		t.Fatal("Apply must not mutate the original map")
	}
}
