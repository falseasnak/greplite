package parser

import (
	"testing"
)

func TestJSONParser(t *testing.T) {
	p := &JSONParser{}

	tests := []struct {
		input   string
		wantNil bool
		field   string
		value   string
	}{
		{`{"level":"info","msg":"hello"}`, false, "level", "info"},
		{`{"count":42}`, false, "count", "42"},
		{`not json`, true, "", ""},
		{``, true, "", ""},
	}

	for _, tt := range tests {
		entry := p.Parse(tt.input)
		if tt.wantNil && entry != nil {
			t.Errorf("Parse(%q): expected nil, got entry", tt.input)
			continue
		}
		if !tt.wantNil && entry == nil {
			t.Errorf("Parse(%q): expected entry, got nil", tt.input)
			continue
		}
		if entry != nil && tt.field != "" {
			if got := entry.Fields[tt.field]; got != tt.value {
				t.Errorf("Parse(%q) field %q = %q, want %q", tt.input, tt.field, got, tt.value)
			}
		}
	}
}

func TestLogfmtParser(t *testing.T) {
	p := &LogfmtParser{}

	tests := []struct {
		input   string
		wantNil bool
		field   string
		value   string
	}{
		{`level=info msg=hello`, false, "level", "info"},
		{`level=info msg="hello world"`, false, "msg", "hello world"},
		{`no equals here`, true, "", ""},
		{``, true, "", ""},
	}

	for _, tt := range tests {
		entry := p.Parse(tt.input)
		if tt.wantNil && entry != nil {
			t.Errorf("Parse(%q): expected nil, got entry", tt.input)
			continue
		}
		if !tt.wantNil && entry == nil {
			t.Errorf("Parse(%q): expected entry, got nil", tt.input)
			continue
		}
		if entry != nil && tt.field != "" {
			if got := entry.Fields[tt.field]; got != tt.value {
				t.Errorf("Parse(%q) field %q = %q, want %q", tt.input, tt.field, got, tt.value)
			}
		}
	}
}

func TestAutoDetect(t *testing.T) {
	if e := Auto(`{"a":"b"}`); e.Format != FormatJSON {
		t.Errorf("expected JSON format")
	}
	if e := Auto(`a=b c=d`); e.Format != FormatLogfmt {
		t.Errorf("expected logfmt format")
	}
	if e := Auto(`plain text`); e.Format != FormatPlain {
		t.Errorf("expected plain format")
	}
}
