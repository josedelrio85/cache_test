package pkg

import "sync"

type SafeCounter interface {
	Inc(string)
	Value(string) int
}

type safeCounter struct {
	mu sync.Mutex
	V  map[string]int
}

func NewSafeCounter() *safeCounter {
	return &safeCounter{
		V: make(map[string]int),
	}
}

func (c *safeCounter) Inc(key string) {
	c.mu.Lock()
	c.V[key]++
	c.mu.Unlock()
}

// Value returns the current value of the counter for the given key.
func (c *safeCounter) Value(key string) int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.V[key]
}
