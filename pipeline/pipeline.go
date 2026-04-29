// Package pipeline wires together the greplite processing stages
// into a single reusable execution unit.
package pipeline

import (
	"io"

	"github.com/example/greplite/dedupe"
	"github.com/example/greplite/filter"
	"github.com/example/greplite/highlight"
	"github.com/example/greplite/output"
	"github.com/example/greplite/parser"
	"github.com/example/greplite/ratelimit"
	"github.com/example/greplite/sampling"
	"github.com/example/greplite/stats"
	"github.com/example/greplite/transform"
	"github.com/example/greplite/truncate"
)

// Config holds all stage dependencies for a Pipeline.
type Config struct {
	Filters   []*filter.Filter
	Sampler   sampling.Sampler
	Deduper   *dedupe.Deduper
	RateLim   *ratelimit.Limiter
	Selector  *transform.Selector
	Truncator *truncate.Truncator
	Highlight *highlight.Highlighter
	Formatter output.Formatter
	Stats     *stats.Tracker
	Out       io.Writer
}

// Pipeline processes a single parsed log line through all configured stages.
// It returns true if the line was emitted to output.
type Pipeline struct {
	cfg Config
}

// New constructs a Pipeline from cfg.
func New(cfg Config) *Pipeline {
	return &Pipeline{cfg: cfg}
}

// Run processes raw text for one line, parsing fields and applying every stage.
// filename and lineNum are forwarded to the formatter.
func (p *Pipeline) Run(raw string, filename string, lineNum int) bool {
	fields := parser.Auto(raw)

	// sampling
	if !p.cfg.Sampler.Keep() {
		return false
	}

	// filter
	for _, f := range p.cfg.Filters {
		if !f.Match(fields, raw) {
			return false
		}
	}

	// deduplication
	if p.cfg.Deduper.IsDuplicate(fields, raw) {
		return false
	}

	// rate limiting
	if !p.cfg.RateLim.Allow() {
		return false
	}

	// field selection / renaming
	if p.cfg.Selector != nil {
		fields = p.cfg.Selector.Apply(fields)
	}

	// truncation
	raw = p.cfg.Truncator.Apply(raw)

	// emit
	p.cfg.Formatter.Write(p.cfg.Out, raw, fields, filename, lineNum, p.cfg.Highlight)

	if p.cfg.Stats != nil {
		p.cfg.Stats.AddLine()
	}
	return true
}
