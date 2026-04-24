package output

import (
	"bytes"
	"strings"
	"testing"
)

func TestWritePlain(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatPlain}
	if err := f.Write(-1, "hello world", nil); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(buf.String()); got != "hello world" {
		t.Errorf("expected 'hello world', got %q", got)
	}
}

func TestWritePlainWithLineNumber(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatPlain, ShowLine: true}
	if err := f.Write(4, "log line", nil); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(buf.String()); got != "5:log line" {
		t.Errorf("expected '5:log line', got %q", got)
	}
}

func TestWriteJSON(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatJSON}
	fields := map[string]string{"level": "error", "msg": "oops"}
	if err := f.Write(-1, "", fields); err != nil {
		t.Fatal(err)
	}
	out := buf.String()
	if !strings.Contains(out, `"level":"error"`) {
		t.Errorf("expected level field in JSON output, got %q", out)
	}
	if !strings.Contains(out, `"msg":"oops"`) {
		t.Errorf("expected msg field in JSON output, got %q", out)
	}
}

func TestWriteJSONFallbackRaw(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatJSON}
	if err := f.Write(-1, "raw text", map[string]string{}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), `"_raw":"raw text"`) {
		t.Errorf("expected _raw field, got %q", buf.String())
	}
}

func TestWriteColorHighlight(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatColor, Highlight: "error"}
	if err := f.Write(-1, "an error occurred", nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "error") {
		t.Errorf("expected highlighted output to contain 'error', got %q", buf.String())
	}
	if !strings.Contains(buf.String(), "\033[") {
		t.Errorf("expected ANSI escape codes in color output")
	}
}

func TestWriteColorWithLineNumber(t *testing.T) {
	var buf bytes.Buffer
	f := &Formatter{Writer: &buf, Format: FormatColor, ShowLine: true}
	if err := f.Write(0, "first line", nil); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(buf.String(), "1:") {
		t.Errorf("expected line number in color output, got %q", buf.String())
	}
}
