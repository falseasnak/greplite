package ratelimit

import (
	"testing"
	"time"
)

func TestNoneAlwaysAllows(t *testing.T) {
	l := None()
	for i := 0; i < 1000; i++ {
		if !l.Allow() {
			t.Fatal("None limiter should always allow")
		}
	}
}

func TestDropModeCappsOutput(t *testing.T) {
	l := New(5, true)
	allowed := 0
	for i := 0; i < 20; i++ {
		if l.Allow() {
			allowed++
		}
	}
	if allowed != 5 {
		t.Fatalf("expected 5 allowed in drop mode, got %d", allowed)
	}
}

func TestResetClearsWindow(t *testing.T) {
	l := New(3, true)
	for i := 0; i < 3; i++ {
		l.Allow()
	}
	if l.Allow() {
		t.Fatal("should be capped after 3 in drop mode")
	}
	l.Reset()
	if !l.Allow() {
		t.Fatal("should be allowed after reset")
	}
}

func TestBlockModeEventuallyAllows(t *testing.T) {
	l := &Limiter{
		max:      2,
		window:   50 * time.Millisecond,
		windowAt: time.Now(),
		drop:     false,
	}
	// consume the window
	l.Allow()
	l.Allow()
	start := time.Now()
	ok := l.Allow() // should block until next window
	elapsed := time.Since(start)
	if !ok {
		t.Fatal("blocking limiter must return true")
	}
	if elapsed < 40*time.Millisecond {
		t.Fatalf("expected block of ~50ms, got %v", elapsed)
	}
}

func TestFromFlagsNone(t *testing.T) {
	l, err := FromFlags("")
	if err != nil {
		t.Fatal(err)
	}
	if l.max != -1 {
		t.Fatal("empty flag should produce no-op limiter")
	}
}

func TestFromFlagsRate(t *testing.T) {
	l, err := FromFlags("100")
	if err != nil {
		t.Fatal(err)
	}
	if l.max != 100 || l.drop {
		t.Fatalf("unexpected limiter config: max=%d drop=%v", l.max, l.drop)
	}
}

func TestFromFlagsDrop(t *testing.T) {
	l, err := FromFlags("50/drop")
	if err != nil {
		t.Fatal(err)
	}
	if l.max != 50 || !l.drop {
		t.Fatalf("unexpected limiter config: max=%d drop=%v", l.max, l.drop)
	}
}

func TestFromFlagsInvalidMode(t *testing.T) {
	_, err := FromFlags("10/fast")
	if err == nil {
		t.Fatal("expected error for unknown mode")
	}
}

func TestFromFlagsInvalidRate(t *testing.T) {
	_, err := FromFlags("abc")
	if err == nil {
		t.Fatal("expected error for non-numeric rate")
	}
}
