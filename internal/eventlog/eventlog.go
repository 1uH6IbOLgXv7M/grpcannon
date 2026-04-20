// Package eventlog provides a bounded, thread-safe ring buffer for recording
// discrete load-test events (errors, retries, circuit-open notices, etc.)
// that can be drained for post-run reporting.
package eventlog

import (
	"sync"
	"time"
)

// Level classifies the severity of an event.
type Level uint8

const (
	LevelInfo Level = iota
	LevelWarn
	LevelError
)

// Event is a single entry stored in the log.
type Event struct {
	At      time.Time
	Level   Level
	Message string
}

// Log is a fixed-capacity ring buffer of Events.
type Log struct {
	mu       sync.Mutex
	buf      []Event
	head     int // next write position
	count    int // total events ever recorded
	capacity int
}

// New returns a Log that retains at most capacity events.
// If capacity is less than 1 it is clamped to 1.
func New(capacity int) *Log {
	if capacity < 1 {
		capacity = 1
	}
	return &Log{
		buf:      make([]Event, capacity),
		capacity: capacity,
	}
}

// Add appends an event to the ring buffer, overwriting the oldest entry
// once the buffer is full.
func (l *Log) Add(level Level, msg string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.buf[l.head] = Event{At: time.Now(), Level: level, Message: msg}
	l.head = (l.head + 1) % l.capacity
	l.count++
}

// Entries returns a snapshot of retained events in chronological order.
func (l *Log) Entries() []Event {
	l.mu.Lock()
	defer l.mu.Unlock()
	retained := l.count
	if retained > l.capacity {
		retained = l.capacity
	}
	out := make([]Event, retained)
	start := (l.head - retained + l.capacity) % l.capacity
	for i := 0; i < retained; i++ {
		out[i] = l.buf[(start+i)%l.capacity]
	}
	return out
}

// Total returns the total number of events ever recorded (including overwritten ones).
func (l *Log) Total() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.count
}
