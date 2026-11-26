package flags

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"sync"
)

// Value represents a typed feature flag value matching Python's bool/int/string coercion.
type Value struct {
	Raw string
}

// AsBool returns true for "true", "yes", "on", "1" (case-insensitive).
func (v Value) AsBool() bool {
	switch strings.ToLower(strings.TrimSpace(v.Raw)) {
	case "true", "yes", "on", "1":
		return true
	default:
		return false
	}
}

// AsInt returns the integer value or fallback if not an integer.
func (v Value) AsInt(fallback int) int {
	n, err := strconv.Atoi(strings.TrimSpace(v.Raw))
	if err != nil {
		return fallback
	}
	return n
}

// Store caches feature flags parsed from SMtracker's feature_flags.yaml.
type Store struct {
	mu    sync.RWMutex
	path  string
	flags map[string]Value
}

// NewStore loads feature flags from path (gracefully empty if file is missing).
func NewStore(path string) *Store {
	s := &Store{path: path, flags: make(map[string]Value)}
	_ = s.Reload()
	return s
}

// Reload re-reads the YAML file from disk.
func (s *Store) Reload() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()

	parsed := make(map[string]Value)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = strings.TrimSpace(line[:idx])
		}
		colon := strings.Index(line, ":")
		if colon <= 0 {
			continue
		}
		key := strings.TrimSpace(line[:colon])
		raw := strings.TrimSpace(line[colon+1:])
		raw = strings.Trim(raw, `"'`)
		if key != "" {
			parsed[key] = Value{Raw: raw}
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	s.mu.Lock()
	s.flags = parsed
	s.mu.Unlock()
	return nil
}

// Get returns a flag value, honoring runtime FLAG_<KEY> environment overrides.
func (s *Store) Get(key string) (Value, bool) {
	if env := strings.TrimSpace(os.Getenv("FLAG_" + key)); env != "" {
		return Value{Raw: env}, true
	}
	if s == nil {
		return Value{}, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.flags[key]
	return v, ok
}

// GetBool returns a boolean feature flag or fallback if unset.
func (s *Store) GetBool(key string, fallback bool) bool {
	v, ok := s.Get(key)
	if !ok {
		return fallback
	}
	return v.AsBool()
}

// GetString returns a trimmed string feature flag or fallback if unset/empty.
func (s *Store) GetString(key, fallback string) string {
	v, ok := s.Get(key)
	if !ok || strings.TrimSpace(v.Raw) == "" {
		return fallback
	}
	return strings.TrimSpace(v.Raw)
}

// All returns a copy of all loaded flags.
func (s *Store) All() map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make(map[string]string, len(s.flags))
	for k, v := range s.flags {
		out[k] = v.Raw
	}
	return out
}
