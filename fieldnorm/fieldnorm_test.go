package fieldnorm_test

import (
	"testing"

	"github.com/user/greplite/fieldnorm"
)

func rec(kvs ...interface{}) map[string]interface{} {
	m := make(map[string]interface{}, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		m[kvs[i].(string)] = kvs[i+1]
	}
	return m
}

func TestNonePassesThrough(t *testing.T) {
	n := fieldnorm.None()
	in := rec("level", "  INFO  ", "msg", "hello")
	out := n.Apply(in)
	if out["level"] != "  INFO  " {
		t.Fatalf("expected value unchanged, got %q", out["level"])
	}
}

func TestTrimOp(t *testing.T) {
	n := fieldnorm.New([]fieldnorm.Rule{{Field: "level", Op: fieldnorm.OpTrim}})
	out := n.Apply(rec("level", "  WARN  "))
	if out["level"] != "WARN" {
		t.Fatalf("expected WARN, got %q", out["level"])
	}
}

func TestLowerOp(t *testing.T) {
	n := fieldnorm.New([]fieldnorm.Rule{{Field: "level", Op: fieldnorm.OpLower}})
	out := n.Apply(rec("level", "ERROR"))
	if out["level"] != "error" {
		t.Fatalf("expected error, got %q", out["level"])
	}
}

func TestUpperOp(t *testing.T) {
	n := fieldnorm.New([]fieldnorm.Rule{{Field: "level", Op: fieldnorm.OpUpper}})
	out := n.Apply(rec("level", "debug"))
	if out["level"] != "DEBUG" {
		t.Fatalf("expected DEBUG, got %q", out["level"])
	}
}

func TestTrimLowerOp(t *testing.T) {
	n := fieldnorm.New([]fieldnorm.Rule{{Field: "env", Op: fieldnorm.OpTrimLower}})
	out := n.Apply(rec("env", "  Production  "))
	if out["env"] != "production" {
		t.Fatalf("expected production, got %q", out["env"])
	}
}

func TestTrimUpperOp(t *testing.T) {
	n := fieldnorm.New([]fieldnorm.Rule{{Field: "env", Op: fieldnorm.OpTrimUpper}})
	out := n.Apply(rec("env", "  staging  "))
	if out["env"] != "STAGING" {
		t.Fatalf("expected STAGING, got %q", out["env"])
	}
}

func TestNonStringFieldSkipped(t *testing.T) {
	n := fieldnorm.New([]fieldnorm.Rule{{Field: "count", Op: fieldnorm.OpLower}})
	out := n.Apply(rec("count", 42))
	if out["count"] != 42 {
		t.Fatalf("expected 42 unchanged, got %v", out["count"])
	}
}

func TestMissingFieldIgnored(t *testing.T) {
	n := fieldnorm.New([]fieldnorm.Rule{{Field: "missing", Op: fieldnorm.OpTrim}})
	out := n.Apply(rec("level", "info"))
	if _, ok := out["missing"]; ok {
		t.Fatal("missing field should not be created")
	}
}

func TestDoesNotMutateOriginal(t *testing.T) {
	n := fieldnorm.New([]fieldnorm.Rule{{Field: "level", Op: fieldnorm.OpUpper}})
	in := rec("level", "info")
	_ = n.Apply(in)
	if in["level"] != "info" {
		t.Fatal("original record should not be mutated")
	}
}

func TestRulesReturnsCopy(t *testing.T) {
	rules := []fieldnorm.Rule{{Field: "level", Op: fieldnorm.OpLower}}
	n := fieldnorm.New(rules)
	copy := n.Rules()
	copy[0].Field = "mutated"
	if n.Rules()[0].Field != "level" {
		t.Fatal("Rules() should return a copy, not a reference")
	}
}
