// Package tag provides key-value label attachment for load test requests,
// allowing results to be annotated and grouped by arbitrary metadata.
package tag

import "sync"

// Bag holds an immutable set of string key-value tags.
type Bag struct {
	tags map[string]string
}

// New returns a Bag populated from the provided alternating key, value pairs.
// Panics if an odd number of arguments is supplied.
func New(kvs ...string) Bag {
	if len(kvs)%2 != 0 {
		panic("tag.New: odd number of arguments")
	}
	m := make(map[string]string, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		m[kvs[i]] = kvs[i+1]
	}
	return Bag{tags: m}
}

// Get returns the value for key and whether it was present.
func (b Bag) Get(key string) (string, bool) {
	v, ok := b.tags[key]
	return v, ok
}

// All returns a copy of the underlying tag map.
func (b Bag) All() map[string]string {
	out := make(map[string]string, len(b.tags))
	for k, v := range b.tags {
		out[k] = v
	}
	return out
}

// Len returns the number of tags in the bag.
func (b Bag) Len() int { return len(b.tags) }

// Registry maps named Bags and is safe for concurrent use.
type Registry struct {
	mu   sync.RWMutex
	sets map[string]Bag
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{sets: make(map[string]Bag)}
}

// Set stores a Bag under name, replacing any previous value.
func (r *Registry) Set(name string, b Bag) {
	r.mu.Lock()
	r.sets[name] = b
	r.mu.Unlock()
}

// Get retrieves the Bag stored under name.
func (r *Registry) Get(name string) (Bag, bool) {
	r.mu.RLock()
	b, ok := r.sets[name]
	r.mu.RUnlock()
	return b, ok
}

// Names returns all registered names.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.sets))
	for k := range r.sets {
		out = append(out, k)
	}
	return out
}
