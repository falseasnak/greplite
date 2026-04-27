package aggregate

import (
	"testing"
)

func TestFromFlagsValid(t *testing.T) {
	cfg, err := FromFlags("level", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Field != "level" {
		t.Errorf("expected field 'level', got %q", cfg.Field)
	}
	if cfg.TopN != 0 {
		t.Errorf("expected TopN 0, got %d", cfg.TopN)
	}
}

func TestFromFlagsEmptyField(t *testing.T) {
	_, err := FromFlags("", 5)
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestFromFlagsNegativeTopN(t *testing.T) {
	_, err := FromFlags("level", -1)
	if err == nil {
		t.Fatal("expected error for negative topN")
	}
}

func TestApplyTopN(t *testing.T) {
	cfg := &Config{Field: "level", TopN: 2}
	input := []Result{
		{Value: "info", Count: 10},
		{Value: "error", Count: 5},
		{Value: "debug", Count: 1},
	}
	out := cfg.Apply(input)
	if len(out) != 2 {
		t.Errorf("expected 2 results after Apply, got %d", len(out))
	}
}

func TestApplyTopNUnlimited(t *testing.T) {
	cfg := &Config{Field: "level", TopN: 0}
	input := []Result{
		{Value: "a", Count: 3},
		{Value: "b", Count: 2},
		{Value: "c", Count: 1},
	}
	out := cfg.Apply(input)
	if len(out) != 3 {
		t.Errorf("expected all 3 results, got %d", len(out))
	}
}
