package fieldcount

import (
	"testing"
)

func fields(n int) map[string]interface{} {
	m := make(map[string]interface{}, n)
	for i := 0; i < n; i++ {
		m[string(rune('a'+i))] = i
	}
	return m
}

func TestNoneAcceptsAll(t *testing.T) {
	f := None()
	for _, n := range []int{0, 1, 5, 100} {
		if !f.Accept(fields(n)) {
			t.Errorf("None should accept %d fields", n)
		}
	}
}

func TestMinOnly(t *testing.T) {
	f, err := New(3, -1)
	if err != nil {
		t.Fatal(err)
	}
	if f.Accept(fields(2)) {
		t.Error("should reject 2 fields when min=3")
	}
	if !f.Accept(fields(3)) {
		t.Error("should accept 3 fields when min=3")
	}
	if !f.Accept(fields(100)) {
		t.Error("should accept 100 fields when max is unlimited")
	}
}

func TestMaxOnly(t *testing.T) {
	f, err := New(0, 5)
	if err != nil {
		t.Fatal(err)
	}
	if !f.Accept(fields(5)) {
		t.Error("should accept 5 fields when max=5")
	}
	if f.Accept(fields(6)) {
		t.Error("should reject 6 fields when max=5")
	}
}

func TestMinMax(t *testing.T) {
	f, err := New(2, 4)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct {
		n    int
		want bool
	}{
		{1, false},
		{2, true},
		{3, true},
		{4, true},
		{5, false},
	} {
		if got := f.Accept(fields(tc.n)); got != tc.want {
			t.Errorf("n=%d: got %v, want %v", tc.n, got, tc.want)
		}
	}
}

func TestNewInvalidMin(t *testing.T) {
	if _, err := New(-1, 5); err == nil {
		t.Error("expected error for negative min")
	}
}

func TestNewInvalidMaxLessThanMin(t *testing.T) {
	if _, err := New(5, 3); err == nil {
		t.Error("expected error when max < min")
	}
}

func TestFromFlagsNone(t *testing.T) {
	f, err := FromFlags(Flags{})
	if err != nil {
		t.Fatal(err)
	}
	if !f.Accept(fields(0)) {
		t.Error("None filter should accept empty record")
	}
}

func TestFromFlagsMinOnly(t *testing.T) {
	f, err := FromFlags(Flags{MinFields: "2"})
	if err != nil {
		t.Fatal(err)
	}
	if f.Accept(fields(1)) {
		t.Error("should reject 1 field")
	}
	if !f.Accept(fields(2)) {
		t.Error("should accept 2 fields")
	}
}

func TestFromFlagsBadValue(t *testing.T) {
	if _, err := FromFlags(Flags{MinFields: "abc"}); err == nil {
		t.Error("expected error for non-integer value")
	}
}

func TestStringOutput(t *testing.T) {
	if got := None().String(); got != "fieldcount:none" {
		t.Errorf("unexpected: %s", got)
	}
	f, _ := New(2, 5)
	if got := f.String(); got != "fieldcount:min=2,max=5" {
		t.Errorf("unexpected: %s", got)
	}
}
