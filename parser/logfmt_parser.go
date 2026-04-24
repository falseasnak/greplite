package parser

import (
	"strings"
)

// LogfmtParser parses logfmt-formatted log lines (key=value pairs).
type LogfmtParser struct{}

// Format returns FormatLogfmt.
func (p *LogfmtParser) Format() Format { return FormatLogfmt }

// Parse attempts to parse a logfmt log line.
// A valid logfmt line must contain at least one key=value pair.
func (p *LogfmtParser) Parse(line string) *LogEntry {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil
	}

	fields, ok := parseLogfmt(line)
	if !ok || len(fields) == 0 {
		return nil
	}

	return &LogEntry{
		Raw:    line,
		Fields: fields,
		Format: FormatLogfmt,
	}
}

func parseLogfmt(line string) (map[string]string, bool) {
	fields := make(map[string]string)
	rest := line
	hasKV := false

	for rest != "" {
		rest = strings.TrimLeft(rest, " \t")
		if rest == "" {
			break
		}
		eqIdx := strings.IndexByte(rest, '=')
		if eqIdx <= 0 {
			return fields, hasKV
		}
		key := rest[:eqIdx]
		rest = rest[eqIdx+1:]

		var value string
		if len(rest) > 0 && rest[0] == '"' {
			// quoted value
			end := strings.Index(rest[1:], "\"")
			if end < 0 {
				return fields, hasKV
			}
			value = rest[1 : end+1]
			rest = rest[end+2:]
		} else {
			spIdx := strings.IndexAny(rest, " \t")
			if spIdx < 0 {
				value = rest
				rest = ""
			} else {
				value = rest[:spIdx]
				rest = rest[spIdx:]
			}
		}
		fields[key] = value
		hasKV = true
	}
	return fields, hasKV
}
