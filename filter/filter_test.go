package filter_test

import (
	"testing"

	"github.com/user/greplite/filter"
)

func TestParseEqual(t *testing.T) {
	f, err := filter.Parse("level=error")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field != "level" || f.Op != filter.OpEqual || f.Value != "error" {
		t.Errorf("unexpected filter: %+v", f)
	}
}

func TestParseNotEqual(t *testing.T) {
	f, err := filter.Parse("level!=debug")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field != "level" || f.Op != filter.OpNotEqual || f.Value != "debug" {
		t.Errorf("unexpected filter: %+v", f)
	}
}

func TestParseContains(t *testing.T) {
	f, err := filter.Parse("msg~=timeout")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field != "msg" || f.Op != filter.OpContains || f.Value != "timeout" {
		t.Errorf("unexpected filter: %+v", f)
	}
}

func TestParseRegex(t *testing.T) {
	f, err := filter.Parse("msg/err.*/")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Field != "msg" || f.Op != filter.OpRegex {
		t.Errorf("unexpected filter: %+v", f)
	}
}

func TestParseInvalid(t *testing.T) {
	_, err := filter.Parse("justaplainword")
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestMatchEqual(t *testing.T) {
	f, _ := filter.Parse("level=error")
	fields := map[string]string{"level": "error", "msg": "something failed"}
	if !f.Match(fields) {
		t.Error("expected match")
	}
	fields["level"] = "info"
	if f.Match(fields) {
		t.Error("expected no match")
	}
}

func TestMatchMissingField(t *testing.T) {
	f, _ := filter.Parse("service=web")
	fields := map[string]string{"level": "info"}
	if f.Match(fields) {
		t.Error("expected no match for missing field")
	}
}

func TestMatchContains(t *testing.T) {
	f, _ := filter.Parse("msg~=timeout")
	if !f.Match(map[string]string{"msg": "connection timeout exceeded"}) {
		t.Error("expected contains match")
	}
	if f.Match(map[string]string{"msg": "all good"}) {
		t.Error("expected no match")
	}
}

func TestMatchRegex(t *testing.T) {
	f, _ := filter.Parse("msg/^error.*/")
	if !f.Match(map[string]string{"msg": "error: disk full"}) {
		t.Error("expected regex match")
	}
	if f.Match(map[string]string{"msg": "warning: low memory"}) {
		t.Error("expected no regex match")
	}
}

func TestMatchAll(t *testing.T) {
	f1, _ := filter.Parse("level=error")
	f2, _ := filter.Parse("service=api")
	fields := map[string]string{"level": "error", "service": "api", "msg": "oops"}
	if !filter.MatchAll([]*filter.Filter{f1, f2}, fields) {
		t.Error("expected all filters to match")
	}
	fields["service"] = "worker"
	if filter.MatchAll([]*filter.Filter{f1, f2}, fields) {
		t.Error("expected MatchAll to fail")
	}
}
