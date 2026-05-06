package fieldparse

import (
	"testing"
)

func rec(pairs ...any) map[string]any {
	m := make(map[string]any, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i].(string)] = pairs[i+1]
	}
	return m
}

func TestNonePassesThrough(t *testing.T) {
	r := rec("msg", "hello")
	got := None.Apply(r)
	if got["msg"] != "hello" {
		t.Fatalf("expected msg=hello, got %v", got["msg"])
	}
}

func TestNewEmptyFieldReturnsError(t *testing.T) {
	_, err := New("", "auto")
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestNewBadFormatReturnsError(t *testing.T) {
	_, err := New("msg", "csv")
	if err == nil {
		t.Fatal("expected error for unknown format")
	}
}

func TestApplyJSONField(t *testing.T) {
	p, err := New("payload", "auto")
	if err != nil {
		t.Fatal(err)
	}
	r := rec("payload", `{"level":"info","code":42}`)
	got := p.Apply(r)
	if got["level"] != "info" {
		t.Errorf("expected level=info, got %v", got["level"])
	}
	if got["payload"] != `{"level":"info","code":42}` {
		t.Error("original field should be preserved")
	}
}

func TestApplyLogfmtField(t *testing.T) {
	p, err := New("meta", "logfmt")
	if err != nil {
		t.Fatal(err)
	}
	r := rec("meta", `host=web01 region=us-east`)
	got := p.Apply(r)
	if got["host"] != "web01" {
		t.Errorf("expected host=web01, got %v", got["host"])
	}
	if got["region"] != "us-east" {
		t.Errorf("expected region=us-east, got %v", got["region"])
	}
}

func TestApplyMissingFieldPassesThrough(t *testing.T) {
	p, _ := New("payload", "auto")
	r := rec("msg", "hello")
	got := p.Apply(r)
	if _, ok := got["payload"]; ok {
		t.Error("payload should not appear in output")
	}
	if got["msg"] != "hello" {
		t.Error("existing fields should be preserved")
	}
}

func TestApplyDoesNotOverwriteExistingKeys(t *testing.T) {
	p, _ := New("payload", "auto")
	r := rec("level", "warn", "payload", `{"level":"info"}`)
	got := p.Apply(r)
	if got["level"] != "warn" {
		t.Errorf("existing key should not be overwritten, got %v", got["level"])
	}
}

func TestFromFlagsNone(t *testing.T) {
	p, err := FromFlags("", "auto")
	if err != nil {
		t.Fatal(err)
	}
	if p != None {
		t.Error("expected None parser when field is empty")
	}
}

func TestFromFlagsValid(t *testing.T) {
	p, err := FromFlags("body", "json")
	if err != nil {
		t.Fatal(err)
	}
	if p.field != "body" {
		t.Errorf("expected field=body, got %q", p.field)
	}
}
