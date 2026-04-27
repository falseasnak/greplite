// Package aggregate provides field-based counting and grouping of structured log lines.
package aggregate

import (
	"fmt"
	"io"
	"sort"
	"text/tabwriter"
)

// Counter accumulates counts keyed by a field value.
type Counter struct {
	field  string
	counts map[string]int
	total  int
}

// New returns a Counter that groups by the named field.
func New(field string) *Counter {
	return &Counter{
		field:  field,
		counts: make(map[string]int),
	}
}

// Add records one occurrence for the given parsed fields map.
// If the target field is absent the special key "<missing>" is used.
func (c *Counter) Add(fields map[string]string) {
	val, ok := fields[c.field]
	if !ok {
		val = "<missing>"
	}
	c.counts[val]++
	c.total++
}

// Total returns the total number of lines processed.
func (c *Counter) Total() int { return c.total }

// Results returns key/count pairs sorted by count descending.
func (c *Counter) Results() []Result {
	out := make([]Result, 0, len(c.counts))
	for k, v := range c.counts {
		out = append(out, Result{Value: k, Count: v})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Count != out[j].Count {
			return out[i].Count > out[j].Count
		}
		return out[i].Value < out[j].Value
	})
	return out
}

// Result holds a single aggregation bucket.
type Result struct {
	Value string
	Count int
}

// Print writes a human-readable table to w.
func (c *Counter) Print(w io.Writer) {
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "VALUE\tCOUNT\n")
	for _, r := range c.Results() {
		fmt.Fprintf(tw, "%s\t%d\n", r.Value, r.Count)
	}
	fmt.Fprintf(tw, "---\t---\n")
	fmt.Fprintf(tw, "total\t%d\n", c.total)
	tw.Flush()
}
