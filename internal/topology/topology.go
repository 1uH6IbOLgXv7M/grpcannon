// Package topology provides a round-robin endpoint selector for distributing
// gRPC load across multiple target addresses.
package topology

import (
	"errors"
	"sync/atomic"
)

// ErrNoEndpoints is returned when the topology contains no endpoints.
var ErrNoEndpoints = errors.New("topology: no endpoints configured")

// Topology holds an ordered list of target addresses and selects among them
// using a lock-free round-robin counter.
type Topology struct {
	endpoints []string
	cursor    atomic.Uint64
}

// New creates a Topology from the provided list of endpoint addresses.
// It returns ErrNoEndpoints if the slice is empty.
func New(endpoints []string) (*Topology, error) {
	if len(endpoints) == 0 {
		return nil, ErrNoEndpoints
	}
	cp := make([]string, len(endpoints))
	copy(cp, endpoints)
	return &Topology{endpoints: cp}, nil
}

// Next returns the next endpoint address in round-robin order.
// It is safe for concurrent use.
func (t *Topology) Next() string {
	idx := t.cursor.Add(1) - 1
	return t.endpoints[idx%uint64(len(t.endpoints))]
}

// Len returns the number of endpoints in the topology.
func (t *Topology) Len() int {
	return len(t.endpoints)
}

// Endpoints returns a copy of the configured endpoint addresses.
func (t *Topology) Endpoints() []string {
	out := make([]string, len(t.endpoints))
	copy(out, t.endpoints)
	return out
}
