package fieldredact

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestNonePassesThrough(t *testing.T) {
	r := None()
	rec := map[string]any{"password": "secret", "msg": "hello"}
	out := r.Apply(rec)
	if out["password"] != "secret" {
		t.Fatalf("expected value unchanged, got %v", out["password"])
	}
}

func TestRedactsNamedField(t *testing.T) {
	r := New([]string{"password", "token"}, "")
	rec := map[string]any{"password": "s3cr3t", "token": "abc123", "user": "alice"}
	out := r.Apply(rec)
	if out["password"] != defaultPlaceholder {
		t.Fatalf("expected password redacted, got %v", out["password"])
	}
	if out["token"] != defaultPlaceholder {
		t.Fatalf("expected token redacted, got %v", out["token"])
	}
	if out["user"] != "alice" {
		t.Fatalf("expected user unchanged, got %v", out["user"])
	}
}

func TestCustomPlaceholder(t *testing.T) {
	r := New([]string{"secret"}, "***")
	out := r.Apply(map[string]any{"secret": "value"})
	if out["secret"] != "***" {
		t.Fatalf("expected ***, got %v", out["secret"])
	}
}

func TestDoesNotMutateOriginal(t *testing.T) {
	r := New([]string{"pw"}, "")
	orig := map[string]any{"pw": "original"}
	r.Apply(orig)
	if orig["pw"] != "original" {
		t.Fatal("original record was mutated")
	}
}

func TestFromFlagsNone(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RegisterFlags(fs)
	_ = fs.Parse([]string{})
	r, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Fields()) != 0 {
		t.Fatal("expected no-op redactor")
	}
}

func TestFromFlagsValid(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RegisterFlags(fs)
	_ = fs.Parse([]string{"--redact-fields", "password,token", "--redact-placeholder", "<hidden>"})
	r, err := FromFlags(fs)
	if err != nil {
		t.Fatal(err)
	}
	out := r.Apply(map[string]any{"password": "x", "token": "y", "msg": "z"})
	if out["password"] != "<hidden>" || out["token"] != "<hidden>" {
		t.Fatalf("unexpected output: %v", out)
	}
	if out["msg"] != "z" {
		t.Fatalf("msg should be unchanged: %v", out["msg"])
	}
}

func TestFromFlagsEmptyField(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	RegisterFlags(fs)
	_ = fs.Parse([]string{"--redact-fields", "password,,token"})
	_, err := FromFlags(fs)
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}
