// Package fieldaliases rewrites field keys in a structured log record
// according to user-supplied alias rules of the form "oldkey=newkey".
//
// # Motivation
//
// Different log producers use different key names for the same semantic
// concept (e.g. "message", "msg", "text"). fieldaliases lets you normalise
// those names at query time without modifying the original data.
//
// # Usage
//
//	// Build a mapper from CLI flags.
//	mapper, err := fieldaliases.FromFlags(flagSet)
//
//	// Or construct one directly.
//	mapper, err := fieldaliases.New([]string{"message=msg", "ts=timestamp"})
//
//	// Apply to a parsed record.
//	normalized := mapper.Apply(record)
//
// Keys that have no matching alias rule are passed through unchanged.
// Apply never mutates the original record map.
package fieldaliases
