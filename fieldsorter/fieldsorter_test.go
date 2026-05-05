package fieldsorter_test

import (
	"testing"

	"github.com/user/greplite/fieldsorter"
)

func TestNonePreservesOrder(t *testing.T) {
	s := fieldsorter.None()
	keys := []string{"z", "a", "m", "b"}
	got := s.Sort(keys)
	for i, k := range keys {
		if got[i] != k {
			t.Fatalf("index %d: want %q got %q", i, k, got[i])
		}
	}
}

func TestNoneDoesNotMutateInput(t *testing.T) {
	s := fieldsorter.None()
	keys := []string{"z", "a"}
	orig := []string{"z", "a"}
	s.Sort(keys)
	for i := range keys {
		if keys[i] != orig[i] {
			t.Fatalf("input was mutated at index %d", i)
		}
	}
}

func TestAlphaSortsKeys(t *testing.T) {
	s := fieldsorter.NewAlpha()
	keys := []string{"z", "a", "m", "b"}
	got := s.Sort(keys)
	want := []string{"a", "b","m", "z"}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("index %d: want %q got %q", i, w, got[i])
		}
	}
}

func TestPriorityFieldsFirst(t *testing.T) {
	s := fieldsorter.NewPriority([]string{"time", "level", "msg"})
	keys := []string{"host", "msg", "level", "time", "app"}
	got := s.Sort(keys)
	want := []string{"time", "level", "msg", "app", "host"}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("index %d: want %q got %q", i, w, got[i])
		}
	}
}

func TestPriorityMissingPriorityFields(t *testing.T) {
	s := fieldsorter.NewPriority([]string{"time", "level"})
	keys := []string{"z", "a"}
	got := s.Sort(keys)
	want := []string{"a", "z"}
	for i, w := range want {
		if got[i] != w {
			t.Fatalf("index %d: want %q got %q", i, w, got[i])
		}
	}
}

func TestAlphaEmptyInput(t *testing.T) {
	s := fieldsorter.NewAlpha()
	got := s.Sort([]string{})
	if len(got) != 0 {
		t.Fatalf("expected empty slice, got %v", got)
	}
}
