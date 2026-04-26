package dedupe

import "testing"

func TestNoneNeverDedupe(t *testing.T) {
	d := New(ModeNone, "")
	for i := 0; i < 5; i++ {
		if d.IsDuplicate("same line", nil) {
			t.Fatal("ModeNone should never report duplicate")
		}
	}
}

func TestExactDuplicateDetected(t *testing.T) {
	d := New(ModeExact, "")
	if d.IsDuplicate("hello world", nil) {
		t.Fatal("first occurrence should not be duplicate")
	}
	if !d.IsDuplicate("hello world", nil) {
		t.Fatal("second occurrence should be duplicate")
	}
}

func TestExactDifferentLines(t *testing.T) {
	d := New(ModeExact, "")
	if d.IsDuplicate("line one", nil) {
		t.Fatal("first line should not be duplicate")
	}
	if d.IsDuplicate("line two", nil) {
		t.Fatal("different line should not be duplicate")
	}
	if d.Count() != 2 {
		t.Fatalf("expected 2 unique lines, got %d", d.Count())
	}
}

func TestFieldDuplicateDetected(t *testing.T) {
	d := New(ModeField, "request_id")
	fields1 := map[string]string{"request_id": "abc123", "msg": "start"}
	fields2 := map[string]string{"request_id": "abc123", "msg": "end"}
	fields3 := map[string]string{"request_id": "xyz999", "msg": "other"}

	if d.IsDuplicate("raw1", fields1) {
		t.Fatal("first request_id should not be duplicate")
	}
	if !d.IsDuplicate("raw2", fields2) {
		t.Fatal("same request_id should be duplicate")
	}
	if d.IsDuplicate("raw3", fields3) {
		t.Fatal("different request_id should not be duplicate")
	}
}

func TestFieldFallsBackToRawWhenNoFields(t *testing.T) {
	d := New(ModeField, "request_id")
	if d.IsDuplicate("raw line", nil) {
		t.Fatal("first raw line should not be duplicate")
	}
	if !d.IsDuplicate("raw line", nil) {
		t.Fatal("same raw line should be duplicate when fields nil")
	}
}

func TestReset(t *testing.T) {
	d := New(ModeExact, "")
	d.IsDuplicate("line", nil)
	if d.Count() != 1 {
		t.Fatalf("expected 1, got %d", d.Count())
	}
	d.Reset()
	if d.Count() != 0 {
		t.Fatalf("expected 0 after reset, got %d", d.Count())
	}
	if d.IsDuplicate("line", nil) {
		t.Fatal("after reset, line should not be duplicate")
	}
}
