// Package ratelimit provides output rate limiting to cap the number of
// matching lines emitted per second, preventing runaway output on large
// or continuously streaming log sources.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks how many lines have been emitted within the current
// time window and blocks or drops lines once the cap is reached.
type Limiter struct {
	mu       sync.Mutex
	max      int
	count    int
	windowAt time.Time
	window   time.Duration
	drop     bool // if true, drop instead of block
}

// New creates a Limiter that allows at most maxPerSec lines per second.
// When drop is true excess lines are silently discarded; when false the
// caller blocks until the next window opens.
func New(maxPerSec int, drop bool) *Limiter {
	return &Limiter{
		max:      maxPerSec,
		window:   time.Second,
		windowAt: time.Now(),
		drop:     drop,
	}
}

// None returns a no-op Limiter that never restricts output.
func None() *Limiter { return &Limiter{max: -1} }

// Allow reports whether the caller may emit a line right now.
// If the limiter is in blocking mode it sleeps until the window resets
// and then returns true. In drop mode it returns false immediately when
// the cap is exceeded.
func (l *Limiter) Allow() bool {
	if l.max < 0 {
		return true
	}
	for {
		l.mu.Lock()
		now := time.Now()
		if now.Sub(l.windowAt) >= l.window {
			l.count = 0
			l.windowAt = now
		}
		if l.count < l.max {
			l.count++
			l.mu.Unlock()
			return true
		}
		if l.drop {
			l.mu.Unlock()
			return false
		}
		// blocking mode: wait for window to expire
		wait := l.window - now.Sub(l.windowAt)
		l.mu.Unlock()
		time.Sleep(wait)
	}
}

// Reset resets the window counter, useful between files.
func (l *Limiter) Reset() {
	if l.max < 0 {
		return
	}
	l.mu.Lock()
	l.count = 0
	l.windowAt = time.Now()
	l.mu.Unlock()
}
