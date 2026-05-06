package fieldgroup_test

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/user/greplite/fieldgroup"
)

func newFlags() *pflag.FlagSet {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	fieldgroup.RegisterFlags(fs)
	return fs
}

func TestFromFlagsNone(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse(nil)
	g, err := fieldgroup.FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	out := g.Apply(rec("a", "1"))
	if _, ok := out["dest"]; ok {
		t.Fatal("None grouper should not add fields")
	}
}

func TestFromFlagsValid(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--group-dest=full", "--group-src=first,last", "--group-sep= "})
	g, err := fieldgroup.FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	out := g.Apply(rec("first", "Ada", "last", "Lovelace"))
	if got := out["full"]; got != "Ada Lovelace" {
		t.Fatalf("expected 'Ada Lovelace', got %q", got)
	}
}

func TestFromFlagsMissingSourcesReturnsError(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--group-dest=out"})
	_, err := fieldgroup.FromFlags(fs)
	if err == nil {
		t.Fatal("expected error when --group-src is absent")
	}
}

func TestFromFlagsCustomSep(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--group-dest=addr", "--group-src=host,port", "--group-sep=:"})
	g, err := fieldgroup.FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	out := g.Apply(rec("host", "db", "port", "5432"))
	if got := out["addr"]; got != "db:5432" {
		t.Fatalf("expected 'db:5432', got %q", got)
	}
}
