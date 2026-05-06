// Package fieldmatch provides record-level filtering based on structured
// field values. It supports three matching modes:
//
//   - Exact  (--field-eq field=value)      — the field value must equal the pattern exactly.
//   - Contains (--field-contains field=sub) — the field value must contain the substring.
//   - Regex  (--field-regex field=pattern) — the field value must match the regular expression.
//
// Multiple rules may be supplied; a record is accepted only when ALL rules
// are satisfied (logical AND).
//
// Usage:
//
//	fs := pflag.NewFlagSet("app", pflag.ContinueOnError)
//	fieldmatch.RegisterFlags(fs)
//	fs.Parse(os.Args[1:])
//
//	matcher, err := fieldmatch.FromFlags(fs)
//	if err != nil { log.Fatal(err) }
//
//	if matcher.Accept(record) {
//	    // emit record
//	}
package fieldmatch
