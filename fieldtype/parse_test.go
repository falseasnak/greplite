package fieldtype_test

import (
	"testing"

	"github.com/yourorg/greplite/fieldtype"
)

func TestParseFieldSpecNoKind(t *testing.T) {
	spec, err := fieldtype.ParseFieldSpec("latency")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Field != "latency" || spec.Kind != fieldtype.KindAuto {
		t.Fatalf("unexpected spec: %+v", spec)
	}
}

func TestParseFieldSpecWithKind(t *testing.T) {
	spec, err := fieldtype.ParseFieldSpec("latency:number")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Field != "latency" || spec.Kind != fieldtype.KindNumber {
		t.Fatalf("unexpected spec: %+v", spec)
	}
}

func TestParseFieldSpecEmptyField(t *testing.T) {
	_, err := fieldtype.ParseFieldSpec(":number")
	if err == nil {
		t.Fatal("expected error for empty field name")
	}
}

func TestParseFieldSpecBadKind(t *testing.T) {
	_, err := fieldtype.ParseFieldSpec("field:nope")
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
}

func TestParseFieldSpecsCSV(t *testing.T) {
	specs, err := fieldtype.ParseFieldSpecs("status:number, msg:string, ok:bool")
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) != 3 {
		t.Fatalf("expected 3 specs, got %d", len(specs))
	}
	if specs[0].Field != "status" || specs[0].Kind != fieldtype.KindNumber {
		t.Errorf("unexpected spec[0]: %+v", specs[0])
	}
}

func TestParseFieldSpecsEmpty(t *testing.T) {
	specs, err := fieldtype.ParseFieldSpecs("")
	if err != nil {
		t.Fatal(err)
	}
	if len(specs) != 0 {
		t.Fatalf("expected 0 specs, got %d", len(specs))
	}
}

func TestRegistryResolveKnownField(t *testing.T) {
	specs, _ := fieldtype.ParseFieldSpecs("latency:number")
	reg := fieldtype.NewRegistry(specs)

	v := reg.Resolve("latency", "123")
	if v.Kind != fieldtype.KindNumber || v.Num != 123 {
		t.Fatalf("unexpected value: %+v", v)
	}
}

func TestRegistryResolveUnknownFieldUsesAuto(t *testing.T) {
	reg := fieldtype.NewRegistry(nil)
	v := reg.Resolve("anything", "99")
	if v.Kind != fieldtype.KindNumber {
		t.Fatalf("expected auto-detected KindNumber, got %v", v.Kind)
	}
}

func TestRegistryResolveForcedString(t *testing.T) {
	specs, _ := fieldtype.ParseFieldSpecs("code:string")
	reg := fieldtype.NewRegistry(specs)
	v := reg.Resolve("code", "200")
	if v.Kind != fieldtype.KindString {
		t.Fatalf("expected KindString for forced string field, got %v", v.Kind)
	}
}
