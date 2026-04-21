package scatter_test

import (
	"sync"
	"testing"

	"github.com/example/grpcannon/internal/scatter"
)

// collector is a thread-safe Sink that records every value it receives.
type collector[T any] struct {
	mu   sync.Mutex
	items []T
}

func (c *collector[T]) Send(v T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = append(c.items, v)
}

func (c *collector[T]) All() []T {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := make([]T, len(c.items))
	copy(out, c.items)
	return out
}

func key(s string) string { return s }

func TestNew_NilKeyFn_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for nil keyFn")
		}
	}()
	scatter.New[string](nil)
}

func TestSend_MatchingSink_RoutesCorrectly(t *testing.T) {
	r := scatter.New(key)
	a := &collector[string]{}
	r.Register("a", a)

	r.Send("a")
	r.Send("a")

	if got := len(a.All()); got != 2 {
		t.Fatalf("want 2 items in sink a, got %d", got)
	}
}

func TestSend_NoMatch_ForwardedToFallback(t *testing.T) {
	r := scatter.New(key)
	fb := &collector[string]{}
	r.SetFallback(fb)

	r.Send("unknown")

	if got := len(fb.All()); got != 1 {
		t.Fatalf("want 1 item in fallback, got %d", got)
	}
}

func TestSend_NoMatch_NoFallback_Dropped(t *testing.T) {
	r := scatter.New(key)
	// No sinks, no fallback — Send must not panic.
	r.Send("ghost")
}

func TestSend_MultipleRoutes_EachSinkReceivesOwn(t *testing.T) {
	r := scatter.New(key)
	a, b := &collector[string]{}, &collector[string]{}
	r.Register("a", a)
	r.Register("b", b)

	for i := 0; i < 5; i++ {
		r.Send("a")
		r.Send("b")
	}

	if got := len(a.All()); got != 5 {
		t.Errorf("sink a: want 5, got %d", got)
	}
	if got := len(b.All()); got != 5 {
		t.Errorf("sink b: want 5, got %d", got)
	}
}

func TestLen_ReflectsRegistrations(t *testing.T) {
	r := scatter.New(key)
	if r.Len() != 0 {
		t.Fatalf("want 0, got %d", r.Len())
	}
	r.Register("x", &collector[string]{})
	if r.Len() != 1 {
		t.Fatalf("want 1, got %d", r.Len())
	}
}

func TestRegister_Overwrite_ReplacesOldSink(t *testing.T) {
	r := scatter.New(key)
	old := &collector[string]{}
	new_ := &collector[string]{}
	r.Register("a", old)
	r.Register("a", new_)

	r.Send("a")

	if len(old.All()) != 0 {
		t.Error("old sink should not have received anything")
	}
	if len(new_.All()) != 1 {
		t.Error("new sink should have received one item")
	}
}
