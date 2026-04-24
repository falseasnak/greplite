// Package context provides before/after line context for search matches,
// similar to grep's -B, -A, and -C flags.
package context

// Buffer holds a circular buffer of lines for before-context tracking.
type Buffer struct {
	lines []string
	nums  []int
	size  int
	head  int
	count int
}

// New creates a new Buffer that retains up to n lines.
func New(n int) *Buffer {
	if n <= 0 {
		n = 0
	}
	return &Buffer{
		lines: make([]string, n),
		nums:  make([]int, n),
		size:  n,
	}
}

// Add appends a line to the circular buffer.
func (b *Buffer) Add(line string, lineNum int) {
	if b.size == 0 {
		return
	}
	b.lines[b.head] = line
	b.nums[b.head] = lineNum
	b.head = (b.head + 1) % b.size
	if b.count < b.size {
		b.count++
	}
}

// Lines returns the buffered lines in order (oldest first).
func (b *Buffer) Lines() []string {
	if b.count == 0 {
		return nil
	}
	out := make([]string, b.count)
	start := (b.head - b.count + b.size) % b.size
	for i := 0; i < b.count; i++ {
		out[i] = b.lines[(start+i)%b.size]
	}
	return out
}

// LineNums returns the line numbers corresponding to Lines().
func (b *Buffer) LineNums() []int {
	if b.count == 0 {
		return nil
	}
	out := make([]int, b.count)
	start := (b.head - b.count + b.size) % b.size
	for i := 0; i < b.count; i++ {
		out[i] = b.nums[(start+i)%b.size]
	}
	return out
}

// Reset clears the buffer.
func (b *Buffer) Reset() {
	b.head = 0
	b.count = 0
}

// Tracker manages before/after context state for a stream of lines.
type Tracker struct {
	Before *Buffer
	After  int // number of after-context lines remaining to emit
}

// NewTracker creates a Tracker with the given before/after sizes.
func NewTracker(before, after int) *Tracker {
	return &Tracker{
		Before: New(before),
		After:  0,
	}
}
