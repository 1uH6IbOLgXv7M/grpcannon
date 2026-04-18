// Package label provides key-value annotation for load test requests,
// allowing results to be grouped and filtered by arbitrary metadata.
package label

import "fmt"

// Bag is an immutable set of key-value string labels.
type Bag struct {
	pairs map[string]string
}

// New constructs a Bag from an even-length list of alternating key, value
// strings. It panics if an odd number of arguments is supplied.
func New(kvs ...string) Bag {
	if len(kvs)%2 != 0 {
		panic(fmt.Sprintf("label.New: odd number of arguments (%d)", len(kvs)))
	}
	m := make(map[string]string, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		m[kvs[i]] = kvs[i+1]
	}
	return Bag{pairs: m}
}

// Get returns the value for key and whether it was present.
func (b Bag) Get(key string) (string, bool) {
	v, ok := b.pairs[key]
	return v, ok
}

// All returns a shallow copy of the underlying map.
func (b Bag) All() map[string]string {
	out := make(map[string]string, len(b.pairs))
	for k, v := range b.pairs {
		out[k] = v
	}
	return out
}

// Len returns the number of labels in the bag.
func (b Bag) Len() int { return len(b.pairs) }

// Merge returns a new Bag containing labels from both b and other.
// Labels in other overwrite labels in b when keys collide.
func (b Bag) Merge(other Bag) Bag {
	out := b.All()
	for k, v := range other.pairs {
		out[k] = v
	}
	return Bag{pairs: out}
}
