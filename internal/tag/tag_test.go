package tag_test

import (
	"sort"
	"testing"

	"github.com/example/grpcannon/internal/tag"
)

func TestNew_EvenArgs_ReturnsBag(t *testing.T) {
	b := tag.New("env", "prod", "region", "us-east-1")
	if b.Len() != 2 {
		t.Fatalf("expected 2 tags, got %d", b.Len())
	}
}

func TestNew_OddArgs_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fatal("expected panic for odd args")
		}
	}()
	tag.New("key")
}

func TestNew_NoArgs_ReturnsEmptyBag(t *testing.T) {
	b := tag.New()
	if b.Len() != 0 {
		t.Fatalf("expected 0 tags, got %d", b.Len())
	}
}

func TestGet_ExistingKey_ReturnsValue(t *testing.T) {
	b := tag.New("env", "staging")
	v, ok := b.Get("env")
	if !ok || v != "staging" {
		t.Fatalf("expected staging, got %q ok=%v", v, ok)
	}
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	b := tag.New("env", "prod")
	_, ok := b.Get("missing")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	b := tag.New("a", "1", "b", "2")
	m := b.All()
	m["a"] = "mutated"
	v, _ := b.Get("a")
	if v != "1" {
		t.Fatal("All() should return an independent copy")
	}
}

func TestRegistry_SetAndGet(t *testing.T) {
	r := tag.NewRegistry()
	b := tag.New("k", "v")
	r.Set("default", b)
	got, ok := r.Get("default")
	if !ok {
		t.Fatal("expected bag to be present")
	}
	if got.Len() != 1 {
		t.Fatalf("expected 1 tag, got %d", got.Len())
	}
}

func TestRegistry_Get_Missing_ReturnsFalse(t *testing.T) {
	r := tag.NewRegistry()
	_, ok := r.Get("nope")
	if ok {
		t.Fatal("expected false for missing name")
	}
}

func TestRegistry_Names_ReturnsAll(t *testing.T) {
	r := tag.NewRegistry()
	r.Set("a", tag.New())
	r.Set("b", tag.New())
	names := r.Names()
	sort.Strings(names)
	if len(names) != 2 || names[0] != "a" || names[1] != "b" {
		t.Fatalf("unexpected names: %v", names)
	}
}

func TestRegistry_Names_Empty_ReturnsEmptySlice(t *testing.T) {
	r := tag.NewRegistry()
	names := r.Names()
	if len(names) != 0 {
		t.Fatalf("expected empty names slice, got %v", names)
	}
}
