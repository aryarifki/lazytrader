package cache

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
	l1     *MemoryLRU
	redis  *redis.Client
	prefix string
}

// NewStore connects optionally to Redis and initializes a bounded L1 LRU cache.
func NewStore(host string, port, db, l1MaxSize int) *Store {
	return NewStoreWithPrefix(host, port, db, l1MaxSize, "")
}

// NewStoreWithPrefix initializes a Store with an optional Redis key namespace prefix.
func NewStoreWithPrefix(host string, port, db, l1MaxSize int, prefix string) *Store {
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
	return &Store{l1: l1, redis: client, prefix: prefix}
}

func (s *Store) key(k string) string {
	if s == nil || s.prefix == "" {
		return k
	}
	return s.prefix + k
}

// Redis returns the underlying Redis client (may be nil if Redis is unavailable).
func (s *Store) Redis() *redis.Client {
	if s == nil {
		return nil
	}
	return s.redis
}

// Close closes the underlying Redis connection if open.
func (s *Store) Close() {
	if s != nil && s.redis != nil {
		_ = s.redis.Close()
	}
}

// Get reads from L1 first, then L2 Redis, backfilling L1 with remaining TTL.
func (s *Store) Get(ctx context.Context, key string) ([]byte, bool) {
	if s == nil {
		return nil, false
	}
	if val, ok := s.l1.Get(key); ok {
		return val, true
	}
	if s.redis == nil {
		return nil, false
	}
	rKey := s.key(key)
	raw, err := s.redis.Get(ctx, rKey).Bytes()
	if err != nil {
		return nil, false
	}
	ttl, err := s.redis.TTL(ctx, rKey).Result()
	if err != nil || ttl <= 0 {
		ttl = 5 * time.Minute
	}
	s.l1.Set(key, raw, ttl)
	return raw, true
}

// Set writes to L1 and L2 Redis with identical TTL.
func (s *Store) Set(ctx context.Context, key string, val []byte, ttl time.Duration) {
	if s == nil {
		return
	}
	s.l1.Set(key, val, ttl)
	if s.redis != nil {
		_ = s.redis.Set(ctx, s.key(key), val, ttl).Err()
	}
}

// Delete removes a single key from both L1 and L2 Redis.
func (s *Store) Delete(ctx context.Context, key string) {
	if s == nil {
		return
	}
	s.l1.Delete(key)
	if s.redis != nil && !IsProtectedKey(key) {
		_ = s.redis.Unlink(ctx, s.key(key)).Err()
	}
}

// Stats returns cache diagnostic counts matching InvestOwl `/api/logs/cache/stats`.
func (s *Store) Stats(ctx context.Context) map[string]any {
	if s == nil {
		return map[string]any{
			"redis_connected":      false,
			"redis_keys_count":     int64(0),
			"fallback_memory_keys": 0,
		}
	}
	connected := s.redis != nil
	var redisKeys int64
	if connected {
		n, err := s.redis.DBSize(ctx).Result()
		if err != nil {
			connected = false
		} else {
			redisKeys = n
		}
	}
	return map[string]any{
		"redis_connected":      connected,
		"redis_keys_count":     redisKeys,
		"fallback_memory_keys": s.l1.Len(),
	}
}

// SafePurge deletes matching keys using SCAN + UNLINK, never FLUSHDB/FLUSHALL,
// and strictly skips ProtectedPrefixes ("rate_limit:", "setting:").
func (s *Store) SafePurge(ctx context.Context, pattern string) int {
	if s == nil {
		return 0
	}
	if pattern == "" {
		pattern = "*"
	}
	l1Count := s.l1.InvalidatePattern(pattern)
	if s.redis == nil {
		return l1Count
	}

	scanPattern := s.key(pattern)
	deleted := 0
	var cursor uint64
	var batch []string
	for {
		keys, nextCursor, err := s.redis.Scan(ctx, cursor, scanPattern, 500).Result()
		if err != nil {
			break
		}
		for _, k := range keys {
			unprefixed := strings.TrimPrefix(k, s.prefix)
			if IsProtectedKey(k) || IsProtectedKey(unprefixed) {
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

// PurgeCloudflare purges Cloudflare Edge CDN cache matching InvestOwl's core/cache.py.
func PurgeCloudflare(ctx context.Context, zoneID string, hosts []string, purgeEverything bool) map[string]any {
	token := strings.TrimSpace(os.Getenv("CLOUDFLARE_API_TOKEN"))
	zid := strings.TrimSpace(zoneID)
	if zid == "" {
		zid = strings.TrimSpace(os.Getenv("CLOUDFLARE_ZONE_ID"))
	}

	if token == "" || zid == "" {
		if f, err := os.Open("/root/.cloudflare_credentials"); err == nil {
			scanner := bufio.NewScanner(f)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if strings.HasPrefix(line, "CLOUDFLARE_API_TOKEN=") && token == "" {
					token = strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
				} else if strings.HasPrefix(line, "CLOUDFLARE_ZONE_ID=") && zid == "" {
					zid = strings.TrimSpace(strings.SplitN(line, "=", 2)[1])
				}
			}
			_ = f.Close()
		}
	}

	if token == "" || zid == "" {
		return map[string]any{"success": false, "reason": "No Cloudflare credentials found"}
	}

	var body map[string]any
	if purgeEverything {
		body = map[string]any{"purge_everything": true}
	} else {
		targetHosts := hosts
		if len(targetHosts) == 0 {
			targetHosts = []string{"investowl.avenia.my.id", "avenia.my.id", "www.avenia.my.id"}
		}
		body = map[string]any{"hosts": targetHosts}
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return map[string]any{"success": false, "error": err.Error()}
	}

	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/purge_cache", zid)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return map[string]any{"success": false, "error": err.Error()}
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return map[string]any{"success": false, "error": err.Error()}
	}
	defer resp.Body.Close()

	var decoded map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&decoded); err != nil {
		return map[string]any{"success": false, "error": err.Error()}
	}
	return map[string]any{
		"success": decoded["success"],
		"result":  decoded["result"],
	}
}
