package users

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/aryarifki/lazytrader/backend/internal/core/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const authGateKey = "require_login_full_feature"

// Repository provides PostgreSQL persistence for public.users and public.app_settings.
type Repository struct {
	DB          *pgxpool.Pool
	AdminEmails []string
}

// SyncOrGetByEmail implements auth.UserLookup and mirrors POST /api/users/sync.
func (r *Repository) SyncOrGetByEmail(ctx context.Context, email, name, image, provider string) (*auth.SessionUser, error) {
	u, err := r.SyncUser(ctx, UserSyncRequest{
		Email:    email,
		Name:     strPtrIfNotEmpty(name),
		Image:    strPtrIfNotEmpty(image),
		Provider: strPtrIfNotEmpty(provider),
	})
	if err != nil {
		return nil, err
	}
	return &auth.SessionUser{
		ID:         u.ID,
		Email:      u.Email,
		Name:       deref(u.Name),
		Image:      deref(u.Image),
		Role:       u.Role,
		IsApproved: u.IsApproved,
	}, nil
}

func (r *Repository) SyncUser(ctx context.Context, req UserSyncRequest) (*User, error) {
	u, err := r.GetByEmail(ctx, req.Email)
	if err == nil {
		return u, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	isOwner := auth.IsAdminEmail(req.Email, r.AdminEmails)
	role := "user"
	if isOwner {
		role = "admin"
	}

	var out User
	err = r.DB.QueryRow(ctx, `
		INSERT INTO public.users (email, name, image, provider, role, is_approved, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, email, name, image, provider, role, is_approved, watchlist, created_at
	`, req.Email, req.Name, req.Image, req.Provider, role, isOwner, time.Now().UTC()).Scan(
		&out.ID, &out.Email, &out.Name, &out.Image, &out.Provider, &out.Role, &out.IsApproved, &out.Watchlist, &out.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	var out User
	err := r.DB.QueryRow(ctx, `
		SELECT id, email, name, image, provider, role, is_approved, watchlist, created_at
		FROM public.users
		WHERE LOWER(email) = LOWER($1)
		LIMIT 1
	`, strings.TrimSpace(email)).Scan(
		&out.ID, &out.Email, &out.Name, &out.Image, &out.Provider, &out.Role, &out.IsApproved, &out.Watchlist, &out.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *Repository) GetByID(ctx context.Context, id int64) (*User, error) {
	var out User
	err := r.DB.QueryRow(ctx, `
		SELECT id, email, name, image, provider, role, is_approved, watchlist, created_at
		FROM public.users WHERE id = $1
	`, id).Scan(
		&out.ID, &out.Email, &out.Name, &out.Image, &out.Provider, &out.Role, &out.IsApproved, &out.Watchlist, &out.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *Repository) ListAll(ctx context.Context) ([]User, error) {
	rows, err := r.DB.Query(ctx, `
		SELECT id, email, name, image, provider, role, is_approved, watchlist, created_at
		FROM public.users ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.Image, &u.Provider, &u.Role, &u.IsApproved, &u.Watchlist, &u.CreatedAt); err != nil {
			return nil, err
		}
		list = append(list, u)
	}
	return list, rows.Err()
}

func (r *Repository) SetApproval(ctx context.Context, id int64, approved bool) (*User, error) {
	var out User
	err := r.DB.QueryRow(ctx, `
		UPDATE public.users SET is_approved = $2 WHERE id = $1
		RETURNING id, email, name, image, provider, role, is_approved, watchlist, created_at
	`, id, approved).Scan(
		&out.ID, &out.Email, &out.Name, &out.Image, &out.Provider, &out.Role, &out.IsApproved, &out.Watchlist, &out.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	_, err := r.DB.Exec(ctx, `DELETE FROM public.users WHERE id = $1`, id)
	return err
}

func (r *Repository) UpdateWatchlist(ctx context.Context, email, watchlistJSON string) error {
	_, err := r.DB.Exec(ctx, `UPDATE public.users SET watchlist = $2 WHERE LOWER(email) = LOWER($1)`, email, watchlistJSON)
	return err
}

func (r *Repository) GetAuthGate(ctx context.Context) (bool, error) {
	var val string
	err := r.DB.QueryRow(ctx, `SELECT value FROM public.app_settings WHERE key = $1`, authGateKey).Scan(&val)
	if errors.Is(err, pgx.ErrNoRows) {
		_, insErr := r.DB.Exec(ctx, `INSERT INTO public.app_settings (key, value, updated_at) VALUES ($1, 'true', $2) ON CONFLICT (key) DO NOTHING`, authGateKey, time.Now().UTC())
		if insErr != nil {
			return true, insErr
		}
		return true, nil
	}
	if err != nil {
		return true, err
	}
	return isTruthy(val), nil
}

func (r *Repository) SetAuthGate(ctx context.Context, enabled bool) error {
	strVal := "false"
	if enabled {
		strVal = "true"
	}
	_, err := r.DB.Exec(ctx, `
		INSERT INTO public.app_settings (key, value, updated_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (key) DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at
	`, authGateKey, strVal, time.Now().UTC())
	return err
}

func isTruthy(v string) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes":
		return true
	default:
		return false
	}
}

func strPtrIfNotEmpty(s string) *string {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return &s
}

func deref(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}
