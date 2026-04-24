// Package highlight provides utilities for highlighting matched text
// within log lines using ANSI escape codes.
package highlight

import (
	"regexp"
	"strings"
)

const (
	// ANSI color codes
	Reset     = "\033[0m"
	Bold      = "\033[1m"
	Red       = "\033[31m"
	Green     = "\033[32m"
	Yellow    = "\033[33m"
	Cyan      = "\033[36m"
	BoldRed   = "\033[1;31m"
	BoldGreen = "\033[1;32m"
)

// Highlighter holds configuration for text highlighting.
type Highlighter struct {
	color   string
	noColor bool
}

// New creates a new Highlighter with the given ANSI color code.
// If noColor is true, no ANSI codes are emitted.
func New(color string, noColor bool) *Highlighter {
	return &Highlighter{color: color, noColor: noColor}
}

// Line highlights all occurrences of pattern in line.
// Returns the original line if noColor is set or pattern is empty.
func (h *Highlighter) Line(line, pattern string) string {
	if h.noColor || pattern == "" {
		return line
	}
	re, err := regexp.Compile("(?i)" + regexp.QuoteMeta(pattern))
	if err != nil {
		return line
	}
	return h.applyRegexp(line, re)
}

// LineRegexp highlights all matches of a pre-compiled regexp in line.
func (h *Highlighter) LineRegexp(line string, re *regexp.Regexp) string {
	if h.noColor || re == nil {
		return line
	}
	return h.applyRegexp(line, re)
}

func (h *Highlighter) applyRegexp(line string, re *regexp.Regexp) string {
	var sb strings.Builder
	last := 0
	for _, loc := range re.FindAllStringIndex(line, -1) {
		sb.WriteString(line[last:loc[0]])
		sb.WriteString(h.color)
		sb.WriteString(line[loc[0]:loc[1]])
		sb.WriteString(Reset)
		last = loc[1]
	}
	sb.WriteString(line[last:])
	return sb.String()
}
