package fieldsplit

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
	s, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	if s != None {
		t.Error("expected None when --split-field is not set")
	}
}

func TestFromFlagsValid(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--split-field=addr", "--split-into=host,port", "--split-sep=:"})
	s, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	out := s.Apply(rec("addr", "127.0.0.1:9200"))
	if out["host"] != "127.0.0.1" {
		t.Errorf("host: got %v", out["host"])
	}
	if out["port"] != "9200" {
		t.Errorf("port: got %v", out["port"])
	}
}

func TestFromFlagsMissingDestsReturnsError(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--split-field=addr"})
	_, err := FromFlags(fs)
	if err == nil {
		t.Fatal("expected error when --split-into is not set")
	}
}

func TestFromFlagsDefaultSep(t *testing.T) {
	fs := newFlags()
	_ = fs.Parse([]string{"--split-field=tags", "--split-into=first,second"})
	s, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	out := s.Apply(rec("tags", "alpha,beta"))
	if out["first"] != "alpha" || out["second"] != "beta" {
		t.Errorf("unexpected split result: %v", out)
	}
}
