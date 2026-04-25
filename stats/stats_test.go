package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/greplite/stats"
)

func TestNewTrackerZero(t *testing.T) {
	tr := stats.New()
	if tr.FilesSearched != 0 || tr.FilesMatched != 0 {
		t.Errorf("expected zero counts, got files=%d matched=%d", tr.FilesSearched, tr.FilesMatched)
	}
}

func TestAddFile(t *testing.T) {
	tr := stats.New()
	tr.AddFile(false)
	tr.AddFile(true)
	tr.AddFile(true)

	if tr.FilesSearched != 3 {
		t.Errorf("expected 3 files searched, got %d", tr.FilesSearched)
	}
	if tr.FilesMatched != 2 {
		t.Errorf("expected 2 files matched, got %d", tr.FilesMatched)
	}
}

func TestAddLine(t *testing.T) {
	tr := stats.New()
	for i := 0; i < 10; i++ {
		tr.AddLine(i%3 == 0) // lines 0,3,6,9 match => 4 matches
	}

	if tr.LinesSearched != 10 {
		t.Errorf("expected 10 lines searched, got %d", tr.LinesSearched)
	}
	if tr.LinesMatched != 4 {
		t.Errorf("expected 4 lines matched, got %d", tr.LinesMatched)
	}
}

func TestElapsed(t *testing.T) {
	tr := stats.New()
	time.Sleep(5 * time.Millisecond)
	if tr.Elapsed() < 5*time.Millisecond {
		t.Error("elapsed should be at least 5ms")
	}
}

func TestPrint(t *testing.T) {
	tr := stats.New()
	tr.AddFile(true)
	tr.AddFile(false)
	tr.AddLine(true)
	tr.AddLine(false)
	tr.AddLine(true)

	var buf bytes.Buffer
	tr.Print(&buf)
	out := buf.String()

	for _, want := range []string{
		"Files searched : 2",
		"Files matched  : 1",
		"Lines searched : 3",
		"Lines matched  : 2",
		"Elapsed",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("Print output missing %q\ngot:\n%s", want, out)
		}
	}
}
