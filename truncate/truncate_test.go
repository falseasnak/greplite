package truncate

import (
	"strings"
	"testing"
)

func TestNoneNeverTruncates(t *testing.T) {
	tr := None()
	long := strings.Repeat("a", 1000)
	out, truncated := tr.Apply(long)
	if truncated {
		t.Error("None truncator should never truncate")
	}
	if out != long {
		t.Error("None truncator should return original string")
	}
}

func TestDisabledWhenZero(t *testing.T) {
	tr := New(0)
	if tr.Enabled() {
		t.Error("maxLen=0 should disable truncation")
	}
	out, truncated := tr.Apply("hello world")
	if truncated || out != "hello world" {
		t.Error("disabled truncator should pass through unchanged")
	}
}

func TestShortLineUnchanged(t *testing.T) {
	tr := New(20)
	out, truncated := tr.Apply("short")
	if truncated {
		t.Error("short line should not be truncated")
	}
	if out != "short" {
		t.Errorf("expected %q, got %q", "short", out)
	}
}

func TestExactLengthUnchanged(t *testing.T) {
	tr := New(5)
	out, truncated := tr.Apply("hello")
	if truncated {
		t.Error("exact-length line should not be truncated")
	}
	if out != "hello" {
		t.Errorf("expected %q, got %q", "hello", out)
	}
}

func TestLongLineIsTruncated(t *testing.T) {
	tr := New(10)
	input := "hello world this is a long line"
	out, truncated := tr.Apply(input)
	if !truncated {
		t.Error("expected truncation")
	}
	if len([]rune(out)) > 10 {
		t.Errorf("output rune count %d exceeds max 10", len([]rune(out)))
	}
	if !strings.HasSuffix(out, "...") {
		t.Errorf("expected ellipsis suffix, got %q", out)
	}
}

func TestUnicodeMultibyte(t *testing.T) {
	tr := New(5)
	// Each emoji is 1 rune but multiple bytes
	input := "😀😁😂😃😄😅"
	out, truncated := tr.Apply(input)
	if !truncated {
		t.Error("expected truncation for 6-rune unicode string with max 5")
	}
	if len([]rune(out)) > 5 {
		t.Errorf("output rune count %d exceeds max 5", len([]rune(out)))
	}
}

func TestMaxLenReturned(t *testing.T) {
	tr := New(42)
	if tr.MaxLen() != 42 {
		t.Errorf("expected MaxLen 42, got %d", tr.MaxLen())
	}
}

func TestEnabledTrue(t *testing.T) {
	tr := New(50)
	if !tr.Enabled() {
		t.Error("expected Enabled() == true for positive maxLen")
	}
}
