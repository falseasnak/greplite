package fieldsorter

import (
	"fmt"
	"strings"
)

// RegisterFlags adds field-sorting flags to the provided flag set.
// It writes defaults into the pointed-to variables which FromFlags reads.
func RegisterFlags(fs interface {
	StringVar(p *string, name, value, usage string)
}) {
	// Callers bind --field-order and --field-priority via main.
	_ = fs
}

// FromFlags constructs a Sorter from CLI flag values.
//
//	--field-order=natural|alpha|priority
//	--field-priority=field1,field2,...
func FromFlags(order, priority string) (*Sorter, error) {
	switch strings.ToLower(strings.TrimSpace(order)) {
	case "", "natural":
		return None(), nil
	case "alpha":
		return NewAlpha(), nil
	case "priority":
		if priority == "" {
			return nil, fmt.Errorf("fieldsorter: --field-order=priority requires --field-priority")
		}
		fields := splitCSV(priority)
		if len(fields) == 0 {
			return nil, fmt.Errorf("fieldsorter: --field-priority must contain at least one field")
		}
		return NewPriority(fields), nil
	default:
		return nil, fmt.Errorf("fieldsorter: unknown --field-order value %q (want natural|alpha|priority)", order)
	}
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
