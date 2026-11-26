package forensic

import "github.com/go-chi/chi/v5"

// Routes mounts /api/forensic endpoints using a pre-built Service (shared with foreign_flow).
func Routes(svc *Service) chi.Router {
	h := &Handler{Svc: svc}
	r := chi.NewRouter()
	r.Post("/refresh", h.Refresh)
	r.Get("/foreign-exodus", h.GetForeignExodus)
	return r
}
