package sampling

import (
	"testing"
)

func TestFromFlagsNone(t *testing.T) {
	s, err := FromFlags(Config{})
	if err != nil {
		t.Fatal(err)
	}
	if s.mode != ModeNone {
		t.Fatalf("expected ModeNone, got %v", s.mode)
	}
}

func TestFromFlagsRate(t *testing.T) {
	s, err := FromFlags(Config{Rate: 5})
	if err != nil {
		t.Fatal(err)
	}
	if s.mode != ModeRate || s.rate != 5 {
		t.Fatalf("unexpected sampler: mode=%v rate=%d", s.mode, s.rate)
	}
}

func TestFromFlagsRandom(t *testing.T) {
	s, err := FromFlags(Config{Prob: 0.3, Seed: 7})
	if err != nil {
		t.Fatal(err)
	}
	if s.mode != ModeRandom {
		t.Fatalf("expected ModeRandom, got %v", s.mode)
	}
}

func TestFromFlagsMutuallyExclusive(t *testing.T) {
	_, err := FromFlags(Config{Rate: 2, Prob: 0.5})
	if err == nil {
		t.Fatal("expected error when both Rate and Prob are set")
	}
}

func TestParseRate(t *testing.T) {
	n, err := ParseRate("10")
	if err != nil || n != 10 {
		t.Fatalf("ParseRate(\"10\") = %d, %v", n, err)
	}
	_, err = ParseRate("abc")
	if err == nil {
		t.Fatal("expected error for non-numeric rate")
	}
}

func TestParseProb(t *testing.T) {
	p, err := ParseProb("0.25")
	if err != nil || p != 0.25 {
		t.Fatalf("ParseProb(\"0.25\") = %f, %v", p, err)
	}
	_, err = ParseProb("bad")
	if err == nil {
		t.Fatal("expected error for non-numeric probability")
	}
}
