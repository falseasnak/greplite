// Package parser provides log line parsers for greplite.
//
// It supports three formats:
//
//   - Plain text: unstructured log lines matched as raw strings.
//   - JSON: lines beginning with '{' parsed as JSON objects.
//   - Logfmt: lines containing key=value pairs (optionally quoted).
//
// Use Auto to detect the format automatically, or instantiate a
// specific Parser (JSONParser, LogfmtParser) for known formats.
//
// Example:
//
//	entry := parser.Auto(`{"level":"error","msg":"disk full"}`)
//	fmt.Println(entry.Fields["level"]) // "error"
package parser
