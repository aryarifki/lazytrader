package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all runtime configuration for the LazyTrader server.
type Config struct {
	Host               string
	Port               int
	StaticDistDir      string
	GzipMinSize        int
	DatabaseURL        string
	RedisHost          string
	RedisPort          int
	RedisDBInvestowl   int
	RedisDBTerminal    int
	RedisKeyPrefix     string
	CORSOrigins        []string
	AuthSecret         string
	GoogleClientID     string
	GoogleClientSecret string
	AppBaseURL         string
	AdminEmails        []string
	TerminalAdminKey   string
	FeatureFlagsPath   string
	SMtrackerDir       string
	CloudflareAPIToken string
	CloudflareZoneID   string
}

// Load reads optional .env files and environment variables with safe defaults.
func Load() (*Config, error) {
	_ = godotenv.Load(".env")
	_ = godotenv.Load("../.env")

	authSecret := firstNonEmpty(os.Getenv("AUTH_SECRET"), os.Getenv("NEXTAUTH_SECRET"))
	if authSecret == "" {
		authSecret = "dev-insecure-change-me-in-production"
	}

	cfg := &Config{
		Host:               getEnvOr("APP_HOST", "127.0.0.1"),
		Port:               getEnvInt("APP_PORT", 8099),
		StaticDistDir:      getEnvOr("STATIC_DIST_DIR", "../frontend/dist"),
		GzipMinSize:        getEnvInt("GZIP_MINIMUM_SIZE", 1000),
		DatabaseURL:        os.Getenv("DATABASE_URL"),
		RedisHost:          getEnvOr("REDIS_HOST", "127.0.0.1"),
		RedisPort:          getEnvInt("REDIS_PORT", 6379),
		RedisDBInvestowl:   getEnvInt("REDIS_DB_INVESTOWL", 0),
		RedisDBTerminal:    getEnvInt("REDIS_DB_TERMINAL", 1),
		RedisKeyPrefix:     strings.TrimSpace(os.Getenv("REDIS_KEY_PREFIX")),
		CORSOrigins:        parseAllowlist(os.Getenv("CORS_ORIGINS"), []string{"http://127.0.0.1:8099", "http://localhost:5173"}),
		AuthSecret:         authSecret,
		GoogleClientID:     strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_ID")),
		GoogleClientSecret: strings.TrimSpace(os.Getenv("GOOGLE_CLIENT_SECRET")),
		AppBaseURL:         strings.TrimRight(getEnvOr("APP_BASE_URL", "http://127.0.0.1:8099"), "/"),
		AdminEmails:        parseAllowlist(os.Getenv("ADMIN_EMAILS"), nil),
		TerminalAdminKey:   strings.TrimSpace(firstNonEmpty(os.Getenv("TERMINAL_ADMIN_KEY"), os.Getenv("ADMIN_API_KEY"))),
		FeatureFlagsPath:   getEnvOr("FEATURE_FLAGS_PATH", "/root/SMtracker/config/feature_flags.yaml"),
		SMtrackerDir:       getEnvOr("SMTRACKER_DIR", "/root/SMtracker"),
		CloudflareAPIToken: strings.TrimSpace(os.Getenv("CLOUDFLARE_API_TOKEN")),
		CloudflareZoneID:   strings.TrimSpace(os.Getenv("CLOUDFLARE_ZONE_ID")),
	}
	return cfg, nil
}

func getEnvOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return n
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}

func parseAllowlist(raw string, fallback []string) []string {
	var out []string
	for _, item := range strings.Split(raw, ",") {
		s := strings.TrimRight(strings.TrimSpace(item), "/")
		if s != "" && s != "*" {
			out = append(out, s)
		}
	}
	if len(out) == 0 && len(fallback) > 0 {
		return fallback
	}
	return out
}
