// Package fieldsorter provides ordering of structured log fields in output.
// Fields can be sorted alphabetically, by a fixed priority list, or left in
// their natural (insertion) order.
package fieldsorter

import "sort"

// Mode controls how fields are ordered.
type Mode int

const (
	// ModeNatural preserves the original field order.
	ModeNatural Mode = iota
	// ModeAlpha sorts fields alphabetically.
	ModeAlpha
	// ModePriority places listed fields first, then sorts the rest alphabetically.
	ModePriority
)

// Sorter reorders the keys of a map according to its configured mode.
type Sorter struct {
	mode     Mode
	priority []string
	priIdx   map[string]int
}

// None returns a Sorter that preserves natural order (no-op on already-ordered keys).
func None() *Sorter {
	return &Sorter{mode: ModeNatural}
}

// NewAlpha returns a Sorter that orders keys alphabetically.
func NewAlpha() *Sorter {
	return &Sorter{mode: ModeAlpha}
}

// NewPriority returns a Sorter that places the given fields first (in order),
// followed by any remaining fields sorted alphabetically.
func NewPriority(fields []string) *Sorter {
	idx := make(map[string]int, len(fields))
	for i, f := range fields {
		idx[f] = i
	}
	return &Sorter{mode: ModePriority, priority: fields, priIdx: idx}
}

// Sort returns a new slice of keys ordered according to the Sorter's mode.
// The input slice is not modified.
func (s *Sorter) Sort(keys []string) []string {
	out := make([]string, len(keys))
	copy(out, keys)

	switch s.mode {
	case ModeAlpha:
		sort.Strings(out)
	case ModePriority:
		sort.SliceStable(out, func(i, j int) bool {
			ii, iok := s.priIdx[out[i]]
			ij, jok := s.priIdx[out[j]]
			switch {
			case iok && jok:
				return ii < ij
			case iok:
				return true
			case jok:
				return false
			default:
				return out[i] < out[j]
			}
		})
	}
	return out
}
