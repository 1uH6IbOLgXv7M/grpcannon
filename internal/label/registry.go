package label

import "sync"

// Registry maps string names to Bags, providing a thread-safe store of
// named label sets that can be reused across requests.
type Registry struct {
	mu   sync.RWMutex
	sets map[string]Bag
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{sets: make(map[string]Bag)}
}

// Register stores bag under name, overwriting any previous entry.
func (r *Registry) Register(name string, bag Bag) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.sets[name] = bag
}

// Lookup returns the Bag registered under name and whether it existed.
func (r *Registry) Lookup(name string) (Bag, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	b, ok := r.sets[name]
	return b, ok
}

// Names returns all registered names in undefined order.
func (r *Registry) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.sets))
	for k := range r.sets {
		out = append(out, k)
	}
	return out
}

// Delete removes the entry for name. It is a no-op if name is not registered.
func (r *Registry) Delete(name string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.sets, name)
}
