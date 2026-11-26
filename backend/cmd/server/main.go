package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
	"github.com/aryarifki/lazytrader/backend/internal/core/cache"
	"github.com/aryarifki/lazytrader/backend/internal/core/config"
	"github.com/aryarifki/lazytrader/backend/internal/core/database"
	"github.com/aryarifki/lazytrader/backend/internal/core/flags"
	"github.com/aryarifki/lazytrader/backend/internal/core/middleware"
	"github.com/aryarifki/lazytrader/backend/internal/core/spa"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/bandarmology"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/broker"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/foreign_flow"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/forensic"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/health"
	signalsvc "github.com/aryarifki/lazytrader/backend/internal/domains/investowl/signal"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/stocks"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/system"
	"github.com/aryarifki/lazytrader/backend/internal/domains/investowl/users"
	"github.com/aryarifki/lazytrader/backend/internal/domains/terminal/portal"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[LazyTrader] config error: %v", err)
	}

	ctx := context.Background()
	dbPool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("[LazyTrader] postgres error: %v", err)
	}
	defer dbPool.Close()

	invCache := cache.NewStoreWithPrefix(cfg.RedisHost, cfg.RedisPort, cfg.RedisDBInvestowl, 256, cfg.RedisKeyPrefix)
	defer invCache.Close()
	termCache := cache.NewStoreWithPrefix(cfg.RedisHost, cfg.RedisPort, cfg.RedisDBTerminal, 256, cfg.RedisKeyPrefix)
	defer termCache.Close()

	flagStore := flags.NewStore(cfg.FeatureFlagsPath)

	userRepo := &users.Repository{DB: dbPool, AdminEmails: cfg.AdminEmails}
	guard := &auth.Guard{Secret: cfg.AuthSecret, AdminEmails: cfg.AdminEmails, Users: userRepo}
	oauthHandler := auth.NewOAuthHandler(guard, cfg.GoogleClientID, cfg.GoogleClientSecret, cfg.AppBaseURL)
	usersHandler := &users.Handler{Repo: userRepo, Redis: invCache.Redis(), AdminEmails: cfg.AdminEmails}
	healthHandler := &health.Handler{DB: dbPool}
	forensicSvc := &forensic.Service{
		Repo:  &forensic.Repository{DB: dbPool},
		Cache: invCache,
		Flags: flagStore,
	}

	r := chi.NewRouter()
	r.Use(chimw.Recoverer)
	r.Use(middleware.SecurityHeaders)
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.Use(middleware.Gzip)
	r.Use(middleware.RateLimit(invCache.Redis()))

	// InvestOwl Core & Auth APIs
	r.Get("/api/health", healthHandler.InvestowlHealth)
	r.Mount("/api/auth", oauthHandler.Routes())
	r.Mount("/api/users", users.Routes(usersHandler, guard))
	r.Mount("/api/stocks", stocks.Routes(dbPool))
	r.Mount("/api/broker-flow", broker.Routes(dbPool))
	r.Mount("/api/forensic", forensic.Routes(forensicSvc))
	r.Mount("/api/foreign-flow", foreign_flow.Routes(dbPool, invCache, forensicSvc))
	r.Mount("/api/signal", signalsvc.Routes(dbPool, invCache, flagStore))
	r.Mount("/api/bandar", bandarmology.Routes(bandarmology.Deps{
		DB:       dbPool,
		Cache:    invCache,
		Guard:    guard,
		CFZoneID: cfg.CloudflareZoneID,
	}))
	r.Mount("/api/logs", system.Routes(system.Deps{
		DB:           dbPool,
		Cache:        invCache,
		Guard:        guard,
		SMtrackerDir: cfg.SMtrackerDir,
		Host:         cfg.Host,
		Port:         cfg.Port,
		CFZoneID:     cfg.CloudflareZoneID,
	}))

	// Avenium Terminal APIs (no bei1persen)
	r.Route("/api/v1/terminal", func(tr chi.Router) {
		tr.Get("/health", healthHandler.TerminalHealth)
		tr.Mount("/portal", portal.Routes())
	})

	// Static Vite SPA + Legacy Redirects
	r.Handle("/*", spa.Handler(cfg.StaticDistDir))

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           r,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Printf("[LazyTrader] Listening on http://%s (static=%s)", addr, cfg.StaticDistDir)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[LazyTrader] server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Println("[LazyTrader] Shutdown complete.")
}
