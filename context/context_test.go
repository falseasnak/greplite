package context_test

import (
	"testing"

	"github.com/yourorg/greplite/context"
)

func TestBufferEmpty(t *testing.T) {
	b := context.New(3)
	if got := b.Lines(); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestBufferUnderCapacity(t *testing.T) {
	b := context.New(3)
	b.Add("line1", 1)
	b.Add("line2", 2)
	lines := b.Lines()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "line1" || lines[1] != "line2" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestBufferOverCapacity(t *testing.T) {
	b := context.New(2)
	b.Add("line1", 1)
	b.Add("line2", 2)
	b.Add("line3", 3)
	lines := b.Lines()
	if len(lines) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(lines))
	}
	if lines[0] != "line2" || lines[1] != "line3" {
		t.Errorf("expected [line2 line3], got %v", lines)
	}
}

func TestBufferLineNums(t *testing.T) {
	b := context.New(3)
	b.Add("a", 10)
	b.Add("b", 11)
	nums := b.LineNums()
	if nums[0] != 10 || nums[1] != 11 {
		t.Errorf("unexpected line nums: %v", nums)
	}
}

func TestBufferReset(t *testing.T) {
	b := context.New(3)
	b.Add("x", 1)
	b.Reset()
	if got := b.Lines(); got != nil {
		t.Fatalf("expected nil after reset, got %v", got)
	}
}

func TestBufferZeroSize(t *testing.T) {
	b := context.New(0)
	b.Add("x", 1)
	if got := b.Lines(); got != nil {
		t.Fatalf("expected nil for zero-size buffer, got %v", got)
	}
}

func TestNewTracker(t *testing.T) {
	tr := context.NewTracker(3, 2)
	if tr.Before == nil {
		t.Fatal("expected non-nil Before buffer")
	}
	if tr.After != 0 {
		t.Errorf("expected After=0, got %d", tr.After)
	}
}
