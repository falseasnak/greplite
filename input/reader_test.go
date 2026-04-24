package input

import (
	"bytes"
	"compress/gzip"
	"os"
	"strings"
	"testing"
)

func TestLineReaderBasic(t *testing.T) {
	input := "line one\nline two\nline three\n"
	lr := NewLineReader(strings.NewReader(input))

	var lines []string
	for lr.Next() {
		lines = append(lines, lr.Line())
	}
	if err := lr.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(lines))
	}
	if lines[1] != "line two" {
		t.Errorf("expected 'line two', got %q", lines[1])
	}
}

func TestLineReaderLineNumbers(t *testing.T) {
	input := "a\nb\nc"
	lr := NewLineReader(strings.NewReader(input))

	num := 0
	for lr.Next() {
		num++
		if lr.LineNumber() != num {
			t.Errorf("expected line number %d, got %d", num, lr.LineNumber())
		}
	}
}

func TestLineReaderEmpty(t *testing.T) {
	lr := NewLineReader(strings.NewReader(""))
	if lr.Next() {
		t.Error("expected no lines from empty reader")
	}
}

func TestOpenFilePlain(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "greplite-*.log")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = f.WriteString("hello\nworld\n")
	f.Close()

	rc, err := OpenFile(f.Name())
	if err != nil {
		t.Fatalf("OpenFile: %v", err)
	}
	defer rc.Close()

	lr := NewLineReader(rc)
	var lines []string
	for lr.Next() {
		lines = append(lines, lr.Line())
	}
	if len(lines) != 2 || lines[0] != "hello" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestOpenFileGzip(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/test.log.gz"

	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, _ = gw.Write([]byte("compressed\ndata\n"))
	gw.Close()

	if err := os.WriteFile(path, buf.Bytes(), 0600); err != nil {
		t.Fatal(err)
	}

	rc, err := OpenFile(path)
	if err != nil {
		t.Fatalf("OpenFile gzip: %v", err)
	}
	defer rc.Close()

	lr := NewLineReader(rc)
	var lines []string
	for lr.Next() {
		lines = append(lines, lr.Line())
	}
	if len(lines) != 2 || lines[0] != "compressed" {
		t.Errorf("unexpected lines: %v", lines)
	}
}

func TestOpenFileNotFound(t *testing.T) {
	_, err := OpenFile("/nonexistent/path/file.log")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
