package cache

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

// ProtectedPrefixes are runtime state keys that must NEVER be purged by cache invalidation.
var ProtectedPrefixes = []string{"rate_limit:", "setting:"}

// IsProtectedKey reports whether key starts with a protected runtime prefix.
func IsProtectedKey(key string) bool {
	for _, p := range ProtectedPrefixes {
		if strings.HasPrefix(key, p) {
			return true
		}
	}
	return false
}

// Store provides two-tier L1 (in-memory LRU) + L2 (Redis) caching with safe SCAN+UNLINK purge.
type Store struct {
	l1    *MemoryLRU
	redis *redis.Client
}

// NewStore connects optionally to Redis and initializes a bounded L1 LRU cache.
func NewStore(host string, port, db, l1MaxSize int) *Store {
	l1 := NewMemoryLRU(l1MaxSize)
	addr := fmt.Sprintf("%s:%d", host, port)
	client := redis.NewClient(&redis.Options{
		Addr:         addr,
		DB:           db,
		DialTimeout:  2 * time.Second,
		ReadTimeout:  2 * time.Second,
		WriteTimeout: 2 * time.Second,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		_ = client.Close()
		client = nil
	}
	return &Store{l1: l1, redis: client}
}

// Redis returns the underlying Redis client (may be nil if Redis is unavailable).
func (s *Store) Redis() *redis.Client {
	return s.redis
}

// Close closes the underlying Redis connection if open.
func (s *Store) Close() {
	if s.redis != nil {
		_ = s.redis.Close()
	}
}

// Get reads from L1 first, then L2 Redis, backfilling L1 with remaining TTL.
func (s *Store) Get(ctx context.Context, key string) ([]byte, bool) {
	if val, ok := s.l1.Get(key); ok {
		return val, true
	}
	if s.redis == nil {
		return nil, false
	}
	raw, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, false
	}
	ttl, err := s.redis.TTL(ctx, key).Result()
	if err != nil || ttl <= 0 {
		ttl = 5 * time.Minute
	}
	s.l1.Set(key, raw, ttl)
	return raw, true
}

// Set writes to L1 and L2 Redis with identical TTL.
func (s *Store) Set(ctx context.Context, key string, val []byte, ttl time.Duration) {
	s.l1.Set(key, val, ttl)
	if s.redis != nil {
		_ = s.redis.Set(ctx, key, val, ttl).Err()
	}
}

// SafePurge deletes matching keys using SCAN + UNLINK, never FLUSHDB/FLUSHALL,
// and strictly skips ProtectedPrefixes ("rate_limit:", "setting:").
func (s *Store) SafePurge(ctx context.Context, pattern string) int {
	if pattern == "" {
		pattern = "*"
	}
	l1Count := s.l1.InvalidatePattern(pattern)
	if s.redis == nil {
		return l1Count
	}

	deleted := 0
	var cursor uint64
	var batch []string
	for {
		keys, nextCursor, err := s.redis.Scan(ctx, cursor, pattern, 500).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			if IsProtectedKey(k) {
				continue
			}
			batch = append(batch, k)
			if len(batch) >= 200 {
				if n, err := s.redis.Unlink(ctx, batch...).Result(); err == nil {
					deleted += int(n)
				}
				batch = batch[:0]
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	if len(batch) > 0 {
		if n, err := s.redis.Unlink(ctx, batch...).Result(); err == nil {
			deleted += int(n)
		}
	}
	if l1Count > deleted {
		return l1Count
	}
	return deleted
}
