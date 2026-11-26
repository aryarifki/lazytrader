package broker

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Routes mounts /api/broker-flow endpoints, registering static /whales before /{ticker}.
func Routes(db *pgxpool.Pool) chi.Router {
	repo := &Repository{DB: db}
	h := &Handler{Repo: repo, Svc: &Service{Repo: repo}}

	r := chi.NewRouter()
	r.Get("/whales", h.GetWhalesList)
	r.Get("/whales/{broker_code}", h.GetWhaleDetail)
	r.Get("/{ticker}/latest", h.GetLatest)
	r.Get("/{ticker}/history", h.GetHistory)
	r.Get("/{ticker}/summary", h.GetSummary)
	return r
}
