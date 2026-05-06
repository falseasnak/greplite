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
// Returns nil if the line is empty, does not start with '{', or is invalid JSON.
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
		fields[k] = marshalFieldValue(v)
	}

	return &LogEntry{
		Raw:    line,
		Fields: fields,
		Format: FormatJSON,
	}
}

// marshalFieldValue converts a JSON field value to its string representation.
func marshalFieldValue(v interface{}) string {
	switch val := v.(type) {
	case string:
		return val
	case nil:
		return ""
	default:
		b, _ := json.Marshal(val)
		return string(b)
	}
}
