// Package context implements before/after line context buffering for
// greplite search results.
//
// It mirrors the behaviour of grep's -B (before), -A (after), and -C
// (combined) flags.  A circular [Buffer] keeps the most-recent N lines so
// that, when a match is found, the caller can flush the buffered lines as
// "before" context.  A [Tracker] pairs the before-buffer with a simple
// countdown for after-context lines that should be printed following a match.
//
// Usage:
//
//	tr := context.NewTracker(beforeN, afterN)
//	for each line {
//	    if match {
//	        for _, l := range tr.Before.Lines() { emit(l) }
//	        emit(matchLine)
//	        tr.Before.Reset()
//	        tr.After = afterN
//	    } else if tr.After > 0 {
//	        emit(line)
//	        tr.After--
//	    } else {
//	        tr.Before.Add(line, lineNum)
//	    }
//	}
package context
