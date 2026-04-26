// Package dedupe implements output deduplication for greplite.
//
// It supports three modes:
//
//	ModeNone  – every line is passed through unchanged (default).
//	ModeExact – suppress consecutive or repeated lines with identical raw text.
//	ModeField – suppress lines where a specific structured-log field (e.g.
//	            "request_id" or "trace_id") has already been seen, regardless
//	            of the rest of the line content.
//
// Usage:
//
//	d := dedupe.New(dedupe.ModeExact, "")
//	for _, line := range lines {
//	    if !d.IsDuplicate(line.Raw, line.Fields) {
//	        output(line)
//	    }
//	}
//
// Activate via CLI flags:
//
//	--dedupe          exact line deduplication
//	--dedupe-field F  deduplicate on structured field F
package dedupe
