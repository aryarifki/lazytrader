package stocks

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Routes mounts the /api/stocks router.
func Routes(db *pgxpool.Pool) chi.Router {
	h := &Handler{Repo: &Repository{DB: db}}
	r := chi.NewRouter()
	r.Get("/{ticker}/history", h.GetHistory)
	return r
}
