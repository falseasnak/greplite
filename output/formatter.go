// Package output provides formatting utilities for greplite search results.
package output

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

// Format represents the output format for matched lines.
type Format int

const (
	FormatPlain Format = iota
	FormatJSON
	FormatColor
)

// Formatter writes matched records to an output writer.
type Formatter struct {
	Writer    io.Writer
	Format    Format
	ShowLine  bool
	Highlight string
}

// Write outputs a matched record. lineNum is 0-indexed; pass -1 to omit.
func (f *Formatter) Write(lineNum int, raw string, fields map[string]string) error {
	switch f.Format {
	case FormatJSON:
		return f.writeJSON(lineNum, fields, raw)
	case FormatColor:
		return f.writeColor(lineNum, raw)
	default:
		return f.writePlain(lineNum, raw)
	}
}

func (f *Formatter) writePlain(lineNum int, raw string) error {
	if f.ShowLine && lineNum >= 0 {
		_, err := fmt.Fprintf(f.Writer, "%d:%s\n", lineNum+1, raw)
		return err
	}
	_, err := fmt.Fprintln(f.Writer, raw)
	return err
}

func (f *Formatter) writeJSON(lineNum int, fields map[string]string, raw string) error {
	out := make(map[string]interface{}, len(fields)+1)
	for k, v := range fields {
		out[k] = v
	}
	if f.ShowLine && lineNum >= 0 {
		out["_line"] = lineNum + 1
	}
	if len(out) == 0 {
		out["_raw"] = raw
	}
	b, err := json.Marshal(out)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(f.Writer, "%s\n", b)
	return err
}

func (f *Formatter) writeColor(lineNum int, raw string) error {
	const (
		colorReset  = "\033[0m"
		colorYellow = "\033[33m"
		colorCyan   = "\033[36m"
	)
	line := raw
	if f.Highlight != "" {
		line = strings.ReplaceAll(raw, f.Highlight, colorYellow+f.Highlight+colorReset)
	}
	if f.ShowLine && lineNum >= 0 {
		_, err := fmt.Fprintf(f.Writer, "%s%d:%s%s\n", colorCyan, lineNum+1, colorReset, line)
		return err
	}
	_, err := fmt.Fprintln(f.Writer, line)
	return err
}
