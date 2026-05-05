// Package fieldtype provides type-coercion utilities for structured log fields.
// It allows filtering and display logic to treat field values as strings,
// numbers, booleans, or null rather than always using raw string comparison.
package fieldtype

import (
	"fmt"
	"strconv"
	"strings"
)

// Kind represents the detected or requested type of a log field value.
type Kind int

const (
	KindAuto    Kind = iota // detect at runtime
	KindString             // always treat as string
	KindNumber             // parse as float64
	KindBool               // parse as bool
	KindNull               // explicit null / missing
)

// Value holds a typed representation of a log field.
type Value struct {
	Kind    Kind
	Str     string
	Num     float64
	Bool    bool
	IsNull  bool
}

// Coerce converts a raw string value from a parsed log record into a typed
// Value. When kind is KindAuto the best-fitting type is inferred.
func Coerce(raw string, kind Kind) Value {
	if raw == "" || raw == "null" || raw == "nil" {
		return Value{Kind: KindNull, IsNull: true, Str: raw}
	}
	switch kind {
	case KindString:
		return Value{Kind: KindString, Str: raw}
	case KindNumber:
		n, err := strconv.ParseFloat(raw, 64)
		if err != nil {
			return Value{Kind: KindString, Str: raw}
		}
		return Value{Kind: KindNumber, Num: n, Str: raw}
	case KindBool:
		b, err := strconv.ParseBool(raw)
		if err != nil {
			return Value{Kind: KindString, Str: raw}
		}
		return Value{Kind: KindBool, Bool: b, Str: raw}
	default: // KindAuto
		if b, err := strconv.ParseBool(raw); err == nil {
			return Value{Kind: KindBool, Bool: b, Str: raw}
		}
		if n, err := strconv.ParseFloat(raw, 64); err == nil {
			return Value{Kind: KindNumber, Num: n, Str: raw}
		}
		return Value{Kind: KindString, Str: raw}
	}
}

// String returns a human-readable representation of the value.
func (v Value) String() string {
	if v.IsNull {
		return "<null>"
	}
	switch v.Kind {
	case KindNumber:
		return strconv.FormatFloat(v.Num, 'f', -1, 64)
	case KindBool:
		return strconv.FormatBool(v.Bool)
	default:
		return v.Str
	}
}

// ParseKind converts a user-supplied type name (string, number, bool, auto)
// into a Kind constant. Returns an error for unknown names.
func ParseKind(s string) (Kind, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "", "auto":
		return KindAuto, nil
	case "string", "str":
		return KindString, nil
	case "number", "num", "float", "int":
		return KindNumber, nil
	case "bool", "boolean":
		return KindBool, nil
	case "null":
		return KindNull, nil
	}
	return KindAuto, fmt.Errorf("fieldtype: unknown kind %q (want auto|string|number|bool|null)", s)
}
