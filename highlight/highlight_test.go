package highlight_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/user/greplite/highlight"
)

func TestLineNoMatch(t *testing.T) {
	h := highlight.New(highlight.BoldRed, false)
	result := h.Line("hello world", "xyz")
	if result != "hello world" {
		t.Errorf("expected no change, got %q", result)
	}
}

func TestLineMatch(t *testing.T) {
	h := highlight.New(highlight.BoldRed, false)
	result := h.Line("hello world", "world")
	if !strings.Contains(result, highlight.BoldRed) {
		t.Errorf("expected ANSI code in output, got %q", result)
	}
	if !strings.Contains(result, "world") {
		t.Errorf("expected 'world' in output, got %q", result)
	}
	if !strings.Contains(result, highlight.Reset) {
		t.Errorf("expected reset code in output, got %q", result)
	}
}

func TestLineNoColor(t *testing.T) {
	h := highlight.New(highlight.BoldRed, true)
	result := h.Line("hello world", "world")
	if result != "hello world" {
		t.Errorf("expected plain output when noColor=true, got %q", result)
	}
}

func TestLineEmptyPattern(t *testing.T) {
	h := highlight.New(highlight.BoldRed, false)
	result := h.Line("hello world", "")
	if result != "hello world" {
		t.Errorf("expected no change for empty pattern, got %q", result)
	}
}

func TestLineRegexp(t *testing.T) {
	h := highlight.New(highlight.Cyan, false)
	re := regexp.MustCompile(`\d+`)
	result := h.LineRegexp("error 404 not found", re)
	if !strings.Contains(result, highlight.Cyan) {
		t.Errorf("expected cyan color around number, got %q", result)
	}
	if !strings.Contains(result, "404") {
		t.Errorf("expected '404' preserved in output, got %q", result)
	}
}

func TestLineRegexpNilRegexp(t *testing.T) {
	h := highlight.New(highlight.Cyan, false)
	result := h.LineRegexp("hello", nil)
	if result != "hello" {
		t.Errorf("expected no change for nil regexp, got %q", result)
	}
}

func TestMultipleMatches(t *testing.T) {
	h := highlight.New(highlight.Yellow, false)
	result := h.Line("foo bar foo baz foo", "foo")
	count := strings.Count(result, highlight.Yellow)
	if count != 3 {
		t.Errorf("expected 3 highlights, got %d in %q", count, result)
	}
}

func TestLineRegexpNoColor(t *testing.T) {
	// When noColor is true, LineRegexp should return the original line unchanged.
	h := highlight.New(highlight.Cyan, true)
	re := regexp.MustCompile(`\d+`)
	result := h.LineRegexp("error 404 not found", re)
	if result != "error 404 not found" {
		t.Errorf("expected plain output when noColor=true, got %q", result)
	}
}
