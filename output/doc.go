// Package output provides result formatting for greplite.
//
// It supports three output modes:
//
//   - FormatPlain: raw line text, optionally prefixed with line numbers.
//   - FormatJSON:  structured JSON objects built from parsed fields;
//     falls back to a "_raw" key when no fields are available.
//   - FormatColor: ANSI-colored output with optional keyword highlighting
//     and line-number prefixes.
//
// Formatter is safe for sequential use but is not goroutine-safe; callers
// that write from multiple goroutines must synchronise access externally.
//
// Usage:
//
//	f := &output.Formatter{
//		Writer:    os.Stdout,
//		Format:    output.FormatColor,
//		ShowLine:  true,
//		Highlight: "error",
//	}
//	f.Write(lineIndex, rawLine, parsedFields)
package output
