package parser

// LogEntry represents a parsed log line with extracted fields.
type LogEntry struct {
	Raw    string
	Fields map[string]string
	Format Format
}

// Format indicates the detected log format.
type Format int

const (
	FormatPlain  Format = iota
	FormatJSON
	FormatLogfmt
)

// Parser defines the interface for log line parsers.
type Parser interface {
	// Parse attempts to parse a raw log line into a LogEntry.
	// Returns nil if the line cannot be parsed by this parser.
	Parse(line string) *LogEntry

	// Format returns the format this parser handles.
	Format() Format
}

// Auto detects the format of a log line and parses it.
func Auto(line string) *LogEntry {
	parsers := []Parser{
		&JSONParser{},
		&LogfmtParser{},
	}
	for _, p := range parsers {
		if entry := p.Parse(line); entry != nil {
			return entry
		}
	}
	return &LogEntry{
		Raw:    line,
		Fields: map[string]string{},
		Format: FormatPlain,
	}
}
