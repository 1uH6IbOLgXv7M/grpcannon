package profile

// Stage represents a single step in a concurrency profile, pairing a target
// worker count with the duration for which that count should be maintained.
type Stage struct {
	Workers  int
	Duration interface{ String() string }
}
