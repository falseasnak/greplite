// Package fieldtemplate adds a computed field to every parsed record
// using a Go text/template expression.
//
// # Usage
//
// Register the flags on your flag set:
//
//	fieldtemplate.RegisterFlags(fs)
//
// Then build an Applier after parsing:
//
//	applier, err := fieldtemplate.FromFlags(fs)
//
// Call Apply for each record in your pipeline:
//
//	rec, err = applier.Apply(rec)
//
// # Template syntax
//
// The template receives the record's field map as its dot value.
// Fields are referenced with {{.fieldName}}. Missing keys evaluate to
// an empty string (missingkey=zero semantics).
//
// Example — combine level and msg into a single summary field:
//
//	--template-field summary --template-expr '{{.level}}: {{.msg}}'
package fieldtemplate
