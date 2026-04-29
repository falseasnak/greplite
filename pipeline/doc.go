// Package pipeline provides a composable processing pipeline for greplite.
//
// A Pipeline chains together all log-processing stages in the correct order:
//
//  1. Sampling  – probabilistic or rate-based line dropping
//  2. Filtering – field/value predicate matching
//  3. Deduplication – exact or field-key based duplicate suppression
//  4. Rate limiting – output throttling (drop or block mode)
//  5. Field selection – projection / renaming of structured fields
//  6. Truncation – cap line length before output
//  7. Formatting – plain, JSON, or colour-highlighted output
//  8. Stats – bookkeeping of matched line counts
//
// Usage:
//
//	cfg := pipeline.Config{
//		Filters:   myFilters,
//		Sampler:   sampling.NewNone(),
//		Deduper:   dedupe.New(dedupe.ModeNone, 0),
//		RateLim:   ratelimit.None(),
//		Truncator: truncate.None(),
//		Highlight: highlight.New(pattern, useColor),
//		Formatter: output.PlainFormatter{},
//		Stats:     stats.New(),
//		Out:       os.Stdout,
//	}
//	p := pipeline.New(cfg)
//	for _, line := range lines {
//		p.Run(line, filename, lineNum)
//	}
package pipeline
