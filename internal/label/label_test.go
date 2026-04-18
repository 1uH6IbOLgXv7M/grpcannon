package label

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew_EvenArgs_ReturnsBag(t *testing.T) {
	b := New("env", "prod", "region", "us-east-1")
	assert.Equal(t, 2, b.Len())
}

func TestNew_NoArgs_ReturnsEmptyBag(t *testing.T) {
	b := New()
	assert.Equal(t, 0, b.Len())
}

func TestNew_OddArgs_Panics(t *testing.T) {
	assert.Panics(t, func() { New("key") })
}

func TestGet_ExistingKey_ReturnsValue(t *testing.T) {
	b := New("env", "staging")
	v, ok := b.Get("env")
	require.True(t, ok)
	assert.Equal(t, "staging", v)
}

func TestGet_MissingKey_ReturnsFalse(t *testing.T) {
	b := New("env", "prod")
	_, ok := b.Get("missing")
	assert.False(t, ok)
}

func TestAll_ReturnsCopy(t *testing.T) {
	b := New("k", "v")
	m := b.All()
	m["k"] = "mutated"
	v, _ := b.Get("k")
	assert.Equal(t, "v", v, "original bag should not be mutated")
}

func TestMerge_OtherOverwrites(t *testing.T) {
	a := New("env", "prod", "region", "us-east-1")
	other := New("env", "staging", "host", "localhost")
	merged := a.Merge(other)

	v, _ := merged.Get("env")
	assert.Equal(t, "staging", v)

	v, _ = merged.Get("region")
	assert.Equal(t, "us-east-1", v)

	v, _ = merged.Get("host")
	assert.Equal(t, "localhost", v)

	assert.Equal(t, 3, merged.Len())
}

func TestMerge_DoesNotMutateOriginal(t *testing.T) {
	a := New("env", "prod")
	b := New("env", "staging")
	a.Merge(b)
	v, _ := a.Get("env")
	assert.Equal(t, "prod", v)
}
