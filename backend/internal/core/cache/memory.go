package cache

import (
	"container/list"
	"path/filepath"
	"sync"
	"time"
)

type entry struct {
	key       string
	value     []byte
	expiresAt time.Time
}

// MemoryLRU is a thread-safe bounded L1 cache with TTL expiration and LRU eviction.
type MemoryLRU struct {
	mu      sync.Mutex
	maxsize int
	items   map[string]*list.Element
	order   *list.List
}

// NewMemoryLRU constructs a bounded in-memory cache (default maxsize=256).
func NewMemoryLRU(maxsize int) *MemoryLRU {
	if maxsize <= 0 {
		maxsize = 256
	}
	return &MemoryLRU{
		maxsize: maxsize,
		items:   make(map[string]*list.Element, maxsize),
		order:   list.New(),
	}
}

// Get returns cached bytes and remaining TTL if present and unexpired.
func (m *MemoryLRU) Get(key string) ([]byte, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	el, ok := m.items[key]
	if !ok {
		return nil, false
	}
	it := el.Value.(*entry)
	if time.Now().After(it.expiresAt) {
		m.removeElement(el)
		return nil, false
	}
	m.order.MoveToBack(el)
	cp := make([]byte, len(it.value))
	copy(cp, it.value)
	return cp, true
}

// Set inserts or updates a key with the given TTL, evicting the oldest entry if full.
func (m *MemoryLRU) Set(key string, val []byte, ttl time.Duration) {
	if ttl <= 0 {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	cp := make([]byte, len(val))
	copy(cp, val)

	if el, ok := m.items[key]; ok {
		it := el.Value.(*entry)
		it.value = cp
		it.expiresAt = time.Now().Add(ttl)
		m.order.MoveToBack(el)
		return
	}
	if m.order.Len() >= m.maxsize {
		if front := m.order.Front(); front != nil {
			m.removeElement(front)
		}
	}
	el := m.order.PushBack(&entry{
		key:       key,
		value:     cp,
		expiresAt: time.Now().Add(ttl),
	})
	m.items[key] = el
}

// Delete removes a single key.
func (m *MemoryLRU) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if el, ok := m.items[key]; ok {
		m.removeElement(el)
	}
}

// InvalidatePattern removes keys matching a glob pattern (e.g. "bandar_*").
func (m *MemoryLRU) InvalidatePattern(pattern string) int {
	m.mu.Lock()
	defer m.mu.Unlock()

	removed := 0
	for k, el := range m.items {
		if matchGlob(pattern, k) {
			m.removeElement(el)
			removed++
		}
	}
	return removed
}

// Len returns current item count.
func (m *MemoryLRU) Len() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.items)
}

func (m *MemoryLRU) removeElement(el *list.Element) {
	it := el.Value.(*entry)
	delete(m.items, it.key)
	m.order.Remove(el)
}

func matchGlob(pattern, key string) bool {
	if pattern == "*" || pattern == "" {
		return true
	}
	matched, err := filepath.Match(pattern, key)
	if err != nil {
		return false
	}
	return matched
}
