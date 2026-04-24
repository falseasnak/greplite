// Package input provides utilities for reading log lines from various sources
// such as files, stdin, and gzip-compressed files.
package input

import (
	"bufio"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

// LineReader reads lines from an io.Reader one at a time.
type LineReader struct {
	scanner *bufio.Scanner
	lineNum int
}

// NewLineReader wraps the given reader in a buffered scanner.
func NewLineReader(r io.Reader) *LineReader {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	return &LineReader{scanner: scanner}
}

// Next advances to the next line. Returns false when there are no more lines
// or an error occurred.
func (lr *LineReader) Next() bool {
	if lr.scanner.Scan() {
		lr.lineNum++
		return true
	}
	return false
}

// Line returns the current line text.
func (lr *LineReader) Line() string {
	return lr.scanner.Text()
}

// LineNumber returns the 1-based current line number.
func (lr *LineReader) LineNumber() int {
	return lr.lineNum
}

// Err returns any error encountered during scanning.
func (lr *LineReader) Err() error {
	return lr.scanner.Err()
}

// OpenFile opens a file for reading, transparently decompressing .gz files.
// The caller is responsible for closing the returned ReadCloser.
func OpenFile(path string) (io.ReadCloser, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}

	if strings.HasSuffix(path, ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			f.Close()
			return nil, fmt.Errorf("gzip open %q: %w", path, err)
		}
		return &gzipReadCloser{gzip: gr, file: f}, nil
	}

	return f, nil
}

// gzipReadCloser closes both the gzip reader and the underlying file.
type gzipReadCloser struct {
	gzip *gzip.Reader
	file *os.File
}

func (g *gzipReadCloser) Read(p []byte) (int, error) {
	return g.gzip.Read(p)
}

func (g *gzipReadCloser) Close() error {
	err := g.gzip.Close()
	if ferr := g.file.Close(); ferr != nil && err == nil {
		err = ferr
	}
	return err
}
