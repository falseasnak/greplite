package aggregate

import (
	"bytes"
	"strings"
	"testing"
)

func TestCounterMissingField(t *testing.T) {
	c := New("level")
	c.Add(map[string]string{"msg": "hello"})
	results := c.Results()
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Value != "<missing>" {
		t.Errorf("expected <missing>, got %q", results[0].Value)
	}
}

func TestCounterBasic(t *testing.T) {
	c := New("level")
	for _, lvl := range []string{"info", "info", "error", "warn", "info"} {
		c.Add(map[string]string{"level": lvl})
	}
	if c.Total() != 5 {
		t.Errorf("expected total 5, got %d", c.Total())
	}
	results := c.Results()
	if results[0].Value != "info" || results[0].Count != 3 {
		t.Errorf("unexpected top result: %+v", results[0])
	}
}

func TestCounterSortStable(t *testing.T) {
	c := New("status")
	for _, s := range []string{"200", "404", "500"} {
		c.Add(map[string]string{"status": s})
	}
	results := c.Results()
	// all counts equal → alphabetical
	if results[0].Value != "200" {
		t.Errorf("expected alphabetical order, got %q first", results[0].Value)
	}
}

func TestPrintOutput(t *testing.T) {
	c := New("env")
	c.Add(map[string]string{"env": "prod"})
	c.Add(map[string]string{"env": "prod"})
	c.Add(map[string]string{"env": "staging"})
	var buf bytes.Buffer
	c.Print(&buf)
	out := buf.String()
	if !strings.Contains(out, "prod") || !strings.Contains(out, "2") {
		t.Errorf("unexpected print output: %s", out)
	}
	if !strings.Contains(out, "total") {
		t.Errorf("missing total line in output")
	}
}

func TestCounterEmpty(t *testing.T) {
	c := New("level")
	if c.Total() != 0 {
		t.Errorf("expected 0 total")
	}
	if len(c.Results()) != 0 {
		t.Errorf("expected empty results")
	}
}
