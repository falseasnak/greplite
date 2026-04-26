// Package sampling provides log-line sampling to reduce output volume
// when searching large files. Lines are selected based on a fixed rate
// (1-in-N) or a random probability (0.0–1.0).
package sampling

import (
	"fmt"
	"math/rand"
)

// Mode controls how sampling is applied.
type Mode int

const (
	// ModeNone disables sampling; every line is forwarded.
	ModeNone Mode = iota
	// ModeRate keeps every Nth line (deterministic).
	ModeRate
	// ModeRandom keeps each line with probability P.
	ModeRandom
)

// Sampler decides whether an individual line should be kept.
type Sampler struct {
	mode Mode
	rate int     // used by ModeRate
	prob float64 // used by ModeRandom
	counter int
	rng     *rand.Rand
}

// NewRate returns a Sampler that keeps every nth line (n >= 1).
func NewRate(n int) (*Sampler, error) {
	if n < 1 {
		return nil, fmt.Errorf("sampling: rate must be >= 1, got %d", n)
	}
	return &Sampler{mode: ModeRate, rate: n, rng: rand.New(rand.NewSource(0))}, nil
}

// NewRandom returns a Sampler that keeps each line with probability p (0 < p <= 1).
func NewRandom(p float64, seed int64) (*Sampler, error) {
	if p <= 0 || p > 1 {
		return nil, fmt.Errorf("sampling: probability must be in (0,1], got %f", p)
	}
	return &Sampler{mode: ModeRandom, prob: p, rng: rand.New(rand.NewSource(seed))}, nil
}

// NewNone returns a pass-through Sampler (no sampling).
func NewNone() *Sampler {
	return &Sampler{mode: ModeNone}
}

// Keep returns true if the current line should be included in output.
func (s *Sampler) Keep() bool {
	switch s.mode {
	case ModeNone:
		return true
	case ModeRate:
		s.counter++
		if s.counter >= s.rate {
			s.counter = 0
			return true
		}
		return false
	case ModeRandom:
		return s.rng.Float64() < s.prob
	}
	return true
}

// Reset resets the internal counter (useful between files).
func (s *Sampler) Reset() {
	s.counter = 0
}
