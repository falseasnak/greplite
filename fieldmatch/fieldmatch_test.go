package fieldmatch

import (
	"testing"
)

func rec(pairs ...string) map[string]interface{} {
	m := make(map[string]interface{}, len(pairs)/2)
	for i := 0; i+1 < len(pairs); i += 2 {
		m[pairs[i]] = pairs[i+1]
	}
	return m
}

func TestNoneAcceptsAll(t *testing.T) {
	m := None()
	if !m.Accept(rec("level", "error")) {
		t.Fatal("None should accept every record")
	}
	if !m.Accept(rec()) {
		t.Fatal("None should accept empty record")
	}
}

func TestExactMatch(t *testing.T) {
	m, err := New([]Rule{{Field: "level", Mode: ModeExact, Pattern: "error"}})
	if err != nil {
		t.Fatal(err)
	}
	if !m.Accept(rec("level", "error")) {
		t.Error("expected accept")
	}
	if m.Accept(rec("level", "warn")) {
		t.Error("expected reject")
	}
}

func TestMissingFieldRejects(t *testing.T) {
	m, _ := New([]Rule{{Field: "level", Mode: ModeExact, Pattern: "error"}})
	if m.Accept(rec("msg", "hello")) {
		t.Error("missing field should reject")
	}
}

func TestContainsMatch(t *testing.T) {
	m, err := New([]Rule{{Field: "msg", Mode: ModeContains, Pattern: "timeout"}})
	if err != nil {
		t.Fatal(err)
	}
	if !m.Accept(rec("msg", "connection timeout occurred")) {
		t.Error("expected accept")
	}
	if m.Accept(rec("msg", "all good")) {
		t.Error("expected reject")
	}
}

func TestRegexMatch(t *testing.T) {
	m, err := New([]Rule{{Field: "code", Mode: ModeRegex, Pattern: "^5\\d{2}$"}})
	if err != nil {
		t.Fatal(err)
	}
	if !m.Accept(rec("code", "500")) {
		t.Error("expected accept for 500")
	}
	if m.Accept(rec("code", "200")) {
		t.Error("expected reject for 200")
	}
}

func TestBadRegexReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "f", Mode: ModeRegex, Pattern: "[invalid"}})
	if err == nil {
		t.Fatal("expected error for bad regex")
	}
}

func TestEmptyFieldReturnsError(t *testing.T) {
	_, err := New([]Rule{{Field: "", Mode: ModeExact, Pattern: "x"}})
	if err == nil {
		t.Fatal("expected error for empty field")
	}
}

func TestMultipleRulesAllMustMatch(t *testing.T) {
	m, err := New([]Rule{
		{Field: "level", Mode: ModeExact, Pattern: "error"},
		{Field: "service", Mode: ModeContains, Pattern: "auth"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !m.Accept(rec("level", "error", "service", "auth-svc")) {
		t.Error("both rules satisfied, expected accept")
	}
	if m.Accept(rec("level", "error", "service", "billing")) {
		t.Error("second rule fails, expected reject")
	}
}
