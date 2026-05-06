package fieldexist_test

import (
	"testing"

	"github.com/spf13/pflag"

	"greplite/fieldexist"
)

func rec(keys ...string) map[string]interface{} {
	m := make(map[string]interface{}, len(keys))
	for _, k := range keys {
		m[k] = "v"
	}
	return m
}

func TestNoneAcceptsAll(t *testing.T) {
	f := fieldexist.None()
	if !f.Apply(rec("a", "b")) {
		t.Fatal("None should accept every record")
	}
	if !f.Apply(rec()) {
		t.Fatal("None should accept empty record")
	}
}

func TestRequireFieldPresent(t *testing.T) {
	f, _ := fieldexist.New([]string{"level", "msg"}, nil)
	if !f.Apply(rec("level", "msg", "ts")) {
		t.Fatal("expected record with required fields to pass")
	}
}

func TestRequireFieldMissing(t *testing.T) {
	f, _ := fieldexist.New([]string{"level", "msg"}, nil)
	if f.Apply(rec("level")) {
		t.Fatal("expected record missing 'msg' to be dropped")
	}
}

func TestExcludeFieldAbsent(t *testing.T) {
	f, _ := fieldexist.New(nil, []string{"debug"})
	if !f.Apply(rec("level", "msg")) {
		t.Fatal("expected record without excluded field to pass")
	}
}

func TestExcludeFieldPresent(t *testing.T) {
	f, _ := fieldexist.New(nil, []string{"debug"})
	if f.Apply(rec("level", "debug")) {
		t.Fatal("expected record with excluded field to be dropped")
	}
}

func TestCaseInsensitiveMatch(t *testing.T) {
	f, _ := fieldexist.New([]string{"Level"}, []string{"DEBUG"})
	if !f.Apply(rec("level")) {
		t.Fatal("require check should be case-insensitive")
	}
	if f.Apply(rec("level", "debug")) {
		t.Fatal("exclude check should be case-insensitive")
	}
}

func TestRequiredAndExcludedFields(t *testing.T) {
	f, _ := fieldexist.New([]string{"msg"}, []string{"error"})
	if f.RequiredFields()[0] != "msg" {
		t.Fatalf("unexpected required field: %v", f.RequiredFields())
	}
	if f.ExcludedFields()[0] != "error" {
		t.Fatalf("unexpected excluded field: %v", f.ExcludedFields())
	}
}

func TestFromFlagsNone(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fieldexist.RegisterFlags(fs)
	f, err := fieldexist.FromFlags(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !f.Apply(rec()) {
		t.Fatal("FromFlags with no flags set should return None")
	}
}

func TestFromFlagsRequire(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fieldexist.RegisterFlags(fs)
	_ = fs.Parse([]string{"--require-fields=level,msg"})
	f, err := fieldexist.FromFlags(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if f.Apply(rec("level")) {
		t.Fatal("record missing 'msg' should be dropped")
	}
	if !f.Apply(rec("level", "msg")) {
		t.Fatal("record with both fields should pass")
	}
}
