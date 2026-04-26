package sampling

import (
	"testing"
)

func TestNoneSamplerKeepsAll(t *testing.T) {
	s := NewNone()
	for i := 0; i < 1000; i++ {
		if !s.Keep() {
			t.Fatalf("ModeNone dropped line %d", i)
		}
	}
}

func TestRateSamplerEveryN(t *testing.T) {
	s, err := NewRate(3)
	if err != nil {
		t.Fatal(err)
	}
	kept := 0
	for i := 0; i < 9; i++ {
		if s.Keep() {
			kept++
		}
	}
	if kept != 3 {
		t.Fatalf("expected 3 kept lines for rate=3 over 9 calls, got %d", kept)
	}
}

func TestRateSamplerReset(t *testing.T) {
	s, _ := NewRate(2)
	s.Keep() // counter=1
	s.Reset()
	// After reset counter is 0; next Keep increments to 1 → not kept.
	if s.Keep() {
		t.Fatal("expected first Keep after Reset to be dropped")
	}
}

func TestRateSamplerInvalidRate(t *testing.T) {
	_, err := NewRate(0)
	if err == nil {
		t.Fatal("expected error for rate=0")
	}
}

func TestRandomSamplerApproximate(t *testing.T) {
	s, err := NewRandom(0.5, 42)
	if err != nil {
		t.Fatal(err)
	}
	kept := 0
	const total = 10000
	for i := 0; i < total; i++ {
		if s.Keep() {
			kept++
		}
	}
	ratio := float64(kept) / total
	if ratio < 0.45 || ratio > 0.55 {
		t.Fatalf("expected ~50%% kept, got %.2f%%", ratio*100)
	}
}

func TestRandomSamplerInvalidProb(t *testing.T) {
	_, err := NewRandom(0, 1)
	if err == nil {
		t.Fatal("expected error for prob=0")
	}
	_, err = NewRandom(1.1, 1)
	if err == nil {
		t.Fatal("expected error for prob=1.1")
	}
}

func TestRandomSamplerProbOne(t *testing.T) {
	s, err := NewRandom(1.0, 0)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 100; i++ {
		if !s.Keep() {
			t.Fatalf("prob=1.0 dropped line %d", i)
		}
	}
}
