// Package highlight provides ANSI-based terminal highlighting for
// greplite search results.
//
// It supports highlighting arbitrary string patterns and pre-compiled
// regular expressions within a line of text. Color output can be
// disabled via the noColor flag to support piped or non-TTY output.
//
// Example usage:
//
//	h := highlight.New(highlight.BoldRed, false)
//	formatted := h.Line("level=error msg=\"disk full\"", "error")
//	fmt.Println(formatted)
//
// Available color constants: Reset, Bold, Red, Green, Yellow, Cyan,
// BoldRed, BoldGreen.
package highlight
