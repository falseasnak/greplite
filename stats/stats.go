// Package stats tracks search result metrics such as matched line counts,
// file counts, and elapsed time for greplite runs.
package stats

import (
	"fmt"
	"io"
	"time"
)

// Tracker accumulates statistics during a search run.
type Tracker struct {
	FilesSearched int
	FilesMatched  int
	LinesSearched int
	LinesMatched  int
	start         time.Time
}

// New returns a new Tracker with the start time set to now.
func New() *Tracker {
	return &Tracker{start: time.Now()}
}

// AddFile records that a file was searched and whether it had any matches.
func (t *Tracker) AddFile(hadMatch bool) {
	t.FilesSearched++
	if hadMatch {
		t.FilesMatched++
	}
}

// AddLine records that a line was examined and whether it matched.
func (t *Tracker) AddLine(matched bool) {
	t.LinesSearched++
	if matched {
		t.LinesMatched++
	}
}

// Elapsed returns the duration since the Tracker was created.
func (t *Tracker) Elapsed() time.Duration {
	return time.Since(t.start)
}

// Print writes a human-readable summary to w.
func (t *Tracker) Print(w io.Writer) {
	fmt.Fprintf(w, "Files searched : %d\n", t.FilesSearched)
	fmt.Fprintf(w, "Files matched  : %d\n", t.FilesMatched)
	fmt.Fprintf(w, "Lines searched : %d\n", t.LinesSearched)
	fmt.Fprintf(w, "Lines matched  : %d\n", t.LinesMatched)
	fmt.Fprintf(w, "Elapsed        : %s\n", t.Elapsed().Round(time.Millisecond))
}
