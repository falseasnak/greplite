package fieldaliases

import (
	"testing"
)

func TestNonePassesThrough(t *testing.T) {
	m := None()
	rec := map[string]any{"msg": "hello", "level": "info"}
	out := m.Apply(rec)
	if out["msg"] != "hello" || out["level"] != "info" {
		t.Fatalf("None should not alter record, got %v", out)
	}
}

func TestRenamesSingleField(t *testing.T) {
	m, err := New([]string{"message=msg"})
	if err != nil {
		t.Fatal(err)
	}
	rec := map[string]any{"message": "hello", "level": "info"}
	out := m.Apply(rec)
	if _, ok := out["message"]; ok {
		t.Error("old key 'message' should have been removed")
	}
	if out["msg"] != "hello" {
		t.Errorf("expected new key 'msg'='hello', got %v", out["msg"])
	}
	if out["level"] != "info" {
		t.Error("unaliased key 'level' should be preserved")
	}
}

func TestRenamesMultipleFields(t *testing.T) {
	m, err := New([]string{"ts=timestamp", "lvl=level"})
	if err != nil {
		t.Fatal(err)
	}
	rec := map[string]any{"ts": "2024-01-01", "lvl": "warn", "msg": "hi"}
	out := m.Apply(rec)
	if out["timestamp"] != "2024-01-01" {
		t.Errorf("expected timestamp=2024-01-01, got %v", out["timestamp"])
	}
	if out["level"] != "warn" {
		t.Errorf("expected level=warn, got %v", out["level"])
	}
	if out["msg"] != "hi" {
		t.Error("unaliased key 'msg' should be preserved")
	}
}

func TestDoesNotMutateOriginal(t *testing.T) {
	m, _ := New([]string{"a=b"})
	rec := map[string]any{"a": 1}
	_ = m.Apply(rec)
	if _, ok := rec["a"]; !ok {
		t.Error("Apply should not mutate the original record")
	}
}

func TestDuplicateSourceReturnsError(t *testing.T) {
	_, err := New([]string{"a=b", "a=c"})
	if err == nil {
		t.Fatal("expected error for duplicate source field")
	}
}

func TestMissingEqualsReturnsError(t *testing.T) {
	_, err := New([]string{"nodivider"})
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestEmptySourceReturnsError(t *testing.T) {
	_, err := New([]string{"=newname"})
	if err == nil {
		t.Fatal("expected error for empty source field")
	}
}

func TestEmptyDestinationReturnsError(t *testing.T) {
	_, err := New([]string{"oldname="})
	if err == nil {
		t.Fatal("expected error for empty destination field")
	}
}

func TestAliasesReturnsCopy(t *testing.T) {
	m, _ := New([]string{"a=b"})
	cp := m.Aliases()
	cp["x"] = "y"
	if _, ok := m.Aliases()["x"]; ok {
		t.Error("Aliases should return an independent copy")
	}
}
