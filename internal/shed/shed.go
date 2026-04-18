// Package shed implements load shedding based on a configurable queue depth.
// When the number of in-flight requests exceeds the limit, new requests are
// rejected immediately rather than queued, protecting downstream services.
package shed

import (
	"errors"
	"sync/atomic"
)

// ErrShed is returned when a request is shed due to queue depth being exceeded.
var ErrShed = errors.New("shed: load too high, request dropped")

// Shed tracks in-flight requests and rejects new ones above a limit.
type Shed struct {
	limit   int64
	inflight atomic.Int64
}

// New returns a Shed that allows at most limit concurrent requests.
// A limit <= 0 means unlimited.
func New(limit int) *Shed {
	return &Shed{limit: int64(limit)}
}

// Acquire attempts to register a new in-flight request.
// It returns ErrShed if the limit is exceeded.
func (s *Shed) Acquire() error {
	if s.limit <= 0 {
		s.inflight.Add(1)
		return nil
	}
	for {
		cur := s.inflight.Load()
		if cur >= s.limit {
			return ErrShed
		}
		if s.inflight.CompareAndSwap(cur, cur+1) {
			return nil
		}
	}
}

// Release decrements the in-flight counter. Must be called after Acquire succeeds.
func (s *Shed) Release() {
	s.inflight.Add(-1)
}

// InFlight returns the current number of in-flight requests.
func (s *Shed) InFlight() int64 {
	return s.inflight.Load()
}
