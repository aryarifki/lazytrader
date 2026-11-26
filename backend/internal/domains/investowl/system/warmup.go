package system

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var defaultWarmupTickers = []string{
	"BBCA", "BMRI", "ASII", "TLKM", "BBRI", "BBNI", "GOTO", "BRPT", "ADRO", "MDKA",
}

func collectWarmupTickers(ctx context.Context, db *pgxpool.Pool) []string {
	seen := make(map[string]struct{}, 32)
	for _, t := range defaultWarmupTickers {
		seen[t] = struct{}{}
	}
	if db != nil {
		rows, err := db.Query(ctx, `
			SELECT DISTINCT jsonb_array_elements_text(watchlist::jsonb)
			FROM users
			WHERE watchlist IS NOT NULL
			  AND watchlist::text != '[]'
			  AND watchlist::text != ''
		`)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var tk string
				if err := rows.Scan(&tk); err == nil {
					if s := strings.ToUpper(strings.TrimSpace(tk)); s != "" {
						seen[s] = struct{}{}
					}
				}
			}
		}
	}
	out := make([]string, 0, len(seen))
	for k := range seen {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// TriggerWarmup pre-warms core endpoints against the local HTTP server in background.
func TriggerWarmup(db *pgxpool.Pool, host string, port int) {
	if host == "" || host == "0.0.0.0" {
		host = "127.0.0.1"
	}
	baseURL := fmt.Sprintf("http://%s:%d", host, port)

	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		tickers := collectWarmupTickers(ctx, db)
		paths := []string{
			"/api/bandar/system-status",
			"/api/bandar/universe/all",
			"/api/bandar/universe/watchlist",
			"/api/bandar/dates/BBCA",
			"/api/bandar/daily-summary?universe_mode=watchlist",
			"/api/signal/daily?limit=100",
			"/api/signal/active?limit=30",
			"/api/foreign-flow/sectors/summary",
		}
		for _, tk := range tickers {
			esc := url.PathEscape(tk)
			paths = append(paths,
				fmt.Sprintf("/api/bandar/detail/%s?window_days=20", esc),
				fmt.Sprintf("/api/bandar/broker-flow/%s?window_days=20&flow_mode=Cumulative&dist_mode=Date%%20range", esc),
				fmt.Sprintf("/api/bandar/causality/%s?window_days=20", esc),
				fmt.Sprintf("/api/bandar/validation-v2/%s?window_days=20&horizon=10&min_events=5&universe_mode=watchlist", esc),
				fmt.Sprintf("/api/bandar/stocks/%s/raw-tables?window_days=20", esc),
				fmt.Sprintf("/api/bandar/dates/%s", esc),
			)
		}

		client := &http.Client{Timeout: 25 * time.Second}
		sem := make(chan struct{}, 4)
		var wg sync.WaitGroup
		for _, p := range paths {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL+path, nil)
				if err != nil {
					return
				}
				resp, err := client.Do(req)
				if err == nil {
					_, _ = io.Copy(io.Discard, resp.Body)
					_ = resp.Body.Close()
				}
			}(p)
		}
		wg.Wait()
	}()
}
