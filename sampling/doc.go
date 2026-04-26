// Package sampling implements configurable line-sampling strategies for
// greplite. When searching extremely large log files it is often useful
// to inspect only a representative subset of matching lines rather than
// printing every match.
//
// Two strategies are provided:
//
//   - Rate sampling  – keeps every Nth line (--sample-rate N).
//     Deterministic; useful for evenly spaced snapshots.
//
//   - Random sampling – keeps each line independently with probability
//     p (--sample-prob 0.1).  The seed can be fixed for reproducibility.
//
// Usage:
//
//	s, _ := sampling.NewRate(10)   // keep 1 in 10
//	for _, line := range lines {
//		if s.Keep() {
//			fmt.Println(line)
//		}
//	}
package sampling
