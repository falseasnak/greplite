// Package aggregate provides field-based value counting for greplite.
//
// When --agg-field is supplied on the command line, greplite switches from
// printing matching lines to accumulating a frequency table keyed by the
// value of the named field in each structured log record.
//
// Usage:
//
//	greplite --agg-field level app.log
//	# prints a table like:
//	#   VALUE    COUNT
//	#   info     1042
//	#   error      17
//	#   warn        5
//	#   ---      ----
//	#   total    1064
//
// Combine with --agg-top N to show only the N most frequent values:
//
//	greplite --agg-field status --agg-top 3 access.log
//
// The package is format-agnostic: it operates on the map[string]string
// produced by the parser package, so it works equally well with JSON
// and logfmt input.
package aggregate
