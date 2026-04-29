package pipeline_test

import (
	"bytes"
	"testing"

	"github.com/example/greplite/dedupe"
	"github.com/example/greplite/filter"
	"github.com/example/greplite/highlight"
	"github.com/example/greplite/output"
	"github.com/example/greplite/pipeline"
	"github.com/example/greplite/ratelimit"
	"github.com/example/greplite/sampling"
	"github.com/example/greplite/stats"
	"github.com/example/greplite/truncate"
)

func defaultCfg(buf *bytes.Buffer) pipeline.Config {
	return pipeline.Config{
		Filters:   nil,
		Sampler:   sampling.NewNone(),
		Deduper:   dedupe.New(dedupe.ModeNone, 0),
		RateLim:   ratelimit.None(),
		Selector:  nil,
		Truncator: truncate.None(),
		Highlight: highlight.New("", false),
		Formatter: output.PlainFormatter{},
		Stats:     stats.New(),
		Out:       buf,
	}
}

func TestRunEmitsMatchingLine(t *testing.T) {
	var buf bytes.Buffer
	p := pipeline.New(defaultCfg(&buf))
	emitted := p.Run(`{"msg":"hello"}`, "test.log", 1)
	if !emitted {
		t.Fatal("expected line to be emitted")
	}
	if buf.Len() == 0 {
		t.Fatal("expected output to be written")
	}
}

func TestRunFilterDropsLine(t *testing.T) {
	var buf bytes.Buffer
	cfg := defaultCfg(&buf)
	f, err := filter.Parse(`msg=world`)
	if err != nil {
		t.Fatalf("parse filter: %v", err)
	}
	cfg.Filters = []*filter.Filter{f}
	p := pipeline.New(cfg)
	emitted := p.Run(`{"msg":"hello"}`, "test.log", 1)
	if emitted {
		t.Fatal("expected line to be dropped by filter")
	}
	if buf.Len() != 0 {
		t.Fatal("expected no output")
	}
}

func TestRunSamplerDropsLine(t *testing.T) {
	var buf bytes.Buffer
	cfg := defaultCfg(&buf)
	// rate=2 keeps every 2nd line; first call should be dropped
	cfg.Sampler = sampling.NewRate(2)
	p := pipeline.New(cfg)
	emitted := p.Run(`{"msg":"hello"}`, "test.log", 1)
	if emitted {
		t.Fatal("expected first line to be dropped by rate sampler")
	}
}

func TestRunDedupeDropsDuplicate(t *testing.T) {
	var buf bytes.Buffer
	cfg := defaultCfg(&buf)
	cfg.Deduper = dedupe.New(dedupe.ModeExact, 100)
	p := pipeline.New(cfg)
	line := `{"msg":"dup"}`
	p.Run(line, "f", 1)
	buf.Reset()
	emitted := p.Run(line, "f", 2)
	if emitted {
		t.Fatal("expected duplicate line to be dropped")
	}
}

func TestRunStatsIncrement(t *testing.T) {
	var buf bytes.Buffer
	cfg := defaultCfg(&buf)
	tr := stats.New()
	cfg.Stats = tr
	p := pipeline.New(cfg)
	p.Run(`hello world`, "f", 1)
	p.Run(`second line`, "f", 2)
	if tr.Lines() != 2 {
		t.Fatalf("expected 2 lines counted, got %d", tr.Lines())
	}
}
