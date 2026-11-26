package cache

import (
	"testing"
	"time"
)

func TestMemoryLRU_EvictionAndTTL(t *testing.T) {
	m := NewMemoryLRU(2)
	m.Set("k1", []byte("v1"), time.Minute)
	m.Set("k2", []byte("v2"), time.Minute)

	// Touch k1 so k2 becomes least-recently-used
	if v, ok := m.Get("k1"); !ok || string(v) != "v1" {
		t.Fatalf("expected k1=v1, got %q ok=%v", string(v), ok)
	}

	m.Set("k3", []byte("v3"), time.Minute)
	if _, ok := m.Get("k2"); ok {
		t.Fatalf("expected k2 to be evicted")
	}
	if _, ok := m.Get("k1"); !ok {
		t.Fatalf("expected k1 to remain")
	}

	m.Set("expired", []byte("old"), time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if _, ok := m.Get("expired"); ok {
		t.Fatalf("expected expired key to miss")
	}
}

func TestIsProtectedKey(t *testing.T) {
	cases := []struct {
		key  string
		want bool
	}{
		{"rate_limit:data:127.0.0.1", true},
		{"setting:require_login_full_feature", true},
		{"bandar_summary:BBCA", false},
		{"ipo_hunter_list", false},
	}
	for _, tc := range cases {
		if got := IsProtectedKey(tc.key); got != tc.want {
			t.Errorf("IsProtectedKey(%q) = %v, want %v", tc.key, got, tc.want)
		}
	}
}
