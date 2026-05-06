package fieldmatch

import (
	"testing"

	"github.com/spf13/pflag"
)

func newFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RegisterFlags(fs)
	return fs
}

func TestFromFlagsNone(t *testing.T) {
	fs := newFlags()
	m, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Accept(rec("any", "value")) {
		t.Error("no flags: should accept all")
	}
}

func TestFromFlagsExact(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--field-eq", "level=error"})
	m, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Accept(rec("level", "error")) {
		t.Error("expected accept")
	}
	if m.Accept(rec("level", "info")) {
		t.Error("expected reject")
	}
}

func TestFromFlagsContains(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--field-contains", "msg=fail"})
	m, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Accept(rec("msg", "operation failed")) {
		t.Error("expected accept")
	}
}

func TestFromFlagsRegex(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--field-regex", "status=^4"})
	m, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	if !m.Accept(rec("status", "404")) {
		t.Error("expected accept")
	}
	if m.Accept(rec("status", "200")) {
		t.Error("expected reject")
	}
}

func TestFromFlagsBadSpec(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--field-eq", "noequalssign"})
	_, err := FromFlags(fs)
	if err == nil {
		t.Fatal("expected error for missing '='")
	}
}

func TestFromFlagsBadRegex(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--field-regex", "f=[bad"})
	_, err := FromFlags(fs)
	if err == nil {
		t.Fatal("expected error for invalid regex")
	}
}
