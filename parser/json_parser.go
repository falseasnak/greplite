package parser

import (
	"encoding/json"
	"strings"
)

// JSONParser parses JSON-formatted log lines.
type JSONParser struct{}

// Format returns FormatJSON.
func (p *JSONParser) Format() Format { return FormatJSON }

// Parse attempts to parse a JSON log line.
func (p *JSONParser) Parse(line string) *LogEntry {
	line = strings.TrimSpace(line)
	if len(line) == 0 || line[0] != '{' {
		return nil
	}

	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil
	}

	fields := make(map[string]string, len(raw))
	for k, v := range raw {
		switch val := v.(type) {
		case string:
			fields[k] = val
		case nil:
			fields[k] = ""
		default:
			b, _ := json.Marshal(val)
			fields[k] = string(b)
		}
	}

	return &LogEntry{
		Raw:    line,
		Fields: fields,
		Format: FormatJSON,
	}
}
