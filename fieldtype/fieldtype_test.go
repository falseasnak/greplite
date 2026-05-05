package fieldtype_test

import (
	"testing"

	"github.com/yourorg/greplite/fieldtype"
)

func TestCoerceAuto(t *testing.T) {
	cases := []struct {
		raw      string
		wantKind fieldtype.Kind
	}{
		{"hello", fieldtype.KindString},
		{"42", fieldtype.KindNumber},
		{"3.14", fieldtype.KindNumber},
		{"true", fieldtype.KindBool},
		{"false", fieldtype.KindBool},
		{"", fieldtype.KindNull},
		{"null", fieldtype.KindNull},
	}
	for _, tc := range cases {
		v := fieldtype.Coerce(tc.raw, fieldtype.KindAuto)
		if v.Kind != tc.wantKind {
			t.Errorf("Coerce(%q, Auto): got kind %v, want %v", tc.raw, v.Kind, tc.wantKind)
		}
	}
}

func TestCoerceForceString(t *testing.T) {
	v := fieldtype.Coerce("42", fieldtype.KindString)
	if v.Kind != fieldtype.KindString {
		t.Fatalf("expected KindString, got %v", v.Kind)
	}
	if v.Str != "42" {
		t.Fatalf("expected Str=42, got %q", v.Str)
	}
}

func TestCoerceNumber(t *testing.T) {
	v := fieldtype.Coerce("1.5", fieldtype.KindNumber)
	if v.Kind != fieldtype.KindNumber || v.Num != 1.5 {
		t.Fatalf("unexpected value: %+v", v)
	}
}

func TestCoerceNumberFallback(t *testing.T) {
	v := fieldtype.Coerce("notanumber", fieldtype.KindNumber)
	if v.Kind != fieldtype.KindString {
		t.Fatalf("expected fallback to KindString, got %v", v.Kind)
	}
}

func TestCoerceBool(t *testing.T) {
	for _, raw := range []string{"true", "false", "TRUE", "1", "0"} {
		v := fieldtype.Coerce(raw, fieldtype.KindBool)
		if v.Kind != fieldtype.KindBool {
			t.Errorf("Coerce(%q, Bool): expected KindBool, got %v", raw, v.Kind)
		}
	}
}

func TestCoerceNull(t *testing.T) {
	for _, raw := range []string{"", "null", "nil"} {
		v := fieldtype.Coerce(raw, fieldtype.KindAuto)
		if !v.IsNull {
			t.Errorf("Coerce(%q): expected IsNull=true", raw)
		}
	}
}

func TestValueString(t *testing.T) {
	cases := []struct {
		v    fieldtype.Value
		want string
	}{
		{fieldtype.Coerce("hello", fieldtype.KindAuto), "hello"},
		{fieldtype.Coerce("3.14", fieldtype.KindAuto), "3.14"},
		{fieldtype.Coerce("true", fieldtype.KindAuto), "true"},
		{fieldtype.Coerce("null", fieldtype.KindAuto), "<null>"},
	}
	for _, tc := range cases {
		if got := tc.v.String(); got != tc.want {
			t.Errorf("String(): got %q, want %q", got, tc.want)
		}
	}
}

func TestParseKind(t *testing.T) {
	cases := []struct {
		input string
		want  fieldtype.Kind
		wantErr bool
	}{
		{"", fieldtype.KindAuto, false},
		{"auto", fieldtype.KindAuto, false},
		{"string", fieldtype.KindString, false},
		{"str", fieldtype.KindString, false},
		{"number", fieldtype.KindNumber, false},
		{"bool", fieldtype.KindBool, false},
		{"null", fieldtype.KindNull, false},
		{"unknown", fieldtype.KindAuto, true},
	}
	for _, tc := range cases {
		got, err := fieldtype.ParseKind(tc.input)
		if tc.wantErr && err == nil {
			t.Errorf("ParseKind(%q): expected error", tc.input)
		}
		if !tc.wantErr && err != nil {
			t.Errorf("ParseKind(%q): unexpected error: %v", tc.input, err)
		}
		if !tc.wantErr && got != tc.want {
			t.Errorf("ParseKind(%q): got %v, want %v", tc.input, got, tc.want)
		}
	}
}
