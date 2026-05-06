package fieldvalidate

import (
	"flag"
	"testing"
)

func rec(kvs ...string) map[string]string {
	m := make(map[string]string, len(kvs)/2)
	for i := 0; i+1 < len(kvs); i += 2 {
		m[kvs[i]] = kvs[i+1]
	}
	return m
}

func TestNoneAlwaysValid(t *testing.T) {
	v := None()
	if !v.Valid(rec()) {
		t.Fatal("None should accept empty record")
	}
	if !v.Valid(rec("level", "info")) {
		t.Fatal("None should accept any record")
	}
}

func TestNonemptyPass(t *testing.T) {
	v := New([]Rule{{Field: "msg", Kind: "nonempty"}})
	if !v.Valid(rec("msg", "hello")) {
		t.Fatal("expected valid")
	}
}

func TestNonemptyFail(t *testing.T) {
	v := New([]Rule{{Field: "msg", Kind: "nonempty"}})
	if v.Valid(rec("msg", "")) {
		t.Fatal("empty value should fail nonempty rule")
	}
}

func TestMissingFieldFails(t *testing.T) {
	v := New([]Rule{{Field: "ts", Kind: "nonempty"}})
	if v.Valid(rec("level", "info")) {
		t.Fatal("missing field should fail")
	}
}

func TestNumberPass(t *testing.T) {
	v := New([]Rule{{Field: "latency", Kind: "number"}})
	if !v.Valid(rec("latency", "3.14")) {
		t.Fatal("expected valid number")
	}
}

func TestNumberFail(t *testing.T) {
	v := New([]Rule{{Field: "latency", Kind: "number"}})
	if v.Valid(rec("latency", "fast")) {
		t.Fatal("non-numeric value should fail number rule")
	}
}

func TestBoolPass(t *testing.T) {
	v := New([]Rule{{Field: "ok", Kind: "bool"}})
	for _, val := range []string{"true", "false", "1", "0"} {
		if !v.Valid(rec("ok", val)) {
			t.Fatalf("expected %q to be valid bool", val)
		}
	}
}

func TestBoolFail(t *testing.T) {
	v := New([]Rule{{Field: "ok", Kind: "bool"}})
	if v.Valid(rec("ok", "yes")) {
		t.Fatal("'yes' should not be a valid bool")
	}
}

func TestRegexPass(t *testing.T) {
	r, err := ParseRule("env:regex:^(prod|staging)$")
	if err != nil {
		t.Fatal(err)
	}
	v := New([]Rule{r})
	if !v.Valid(rec("env", "prod")) {
		t.Fatal("expected prod to match")
	}
}

func TestRegexFail(t *testing.T) {
	r, _ := ParseRule("env:regex:^(prod|staging)$")
	v := New([]Rule{r})
	if v.Valid(rec("env", "dev")) {
		t.Fatal("dev should not match prod|staging")
	}
}

func TestParseRuleUnknownKind(t *testing.T) {
	_, err := ParseRule("field:unknown")
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
}

func TestParseCSV(t *testing.T) {
	v, err := ParseCSV("msg:nonempty,latency:number")
	if err != nil {
		t.Fatal(err)
	}
	if !v.Valid(rec("msg", "hi", "latency", "42")) {
		t.Fatal("expected valid")
	}
	if v.Valid(rec("msg", "hi", "latency", "nan")) {
		t.Fatal("expected invalid")
	}
}

func TestFromFlagsNone(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	RegisterFlags(fs)
	_ = fs.Parse([]string{})
	v, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	if !v.Valid(rec()) {
		t.Fatal("None validator should accept everything")
	}
}

func TestFromFlagsWithRules(t *testing.T) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	RegisterFlags(fs)
	_ = fs.Parse([]string{"--validate", "level:nonempty"})
	v, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	if v.Valid(rec("msg", "hello")) {
		t.Fatal("missing level should fail")
	}
	if !v.Valid(rec("level", "info")) {
		t.Fatal("present level should pass")
	}
}
