package users

import "time"

// User mirrors public.users in PostgreSQL and schemas.UserRead in InvestOwl.
type User struct {
	ID         int64     `json:"id"`
	Email      string    `json:"email"`
	Name       *string   `json:"name"`
	Image      *string   `json:"image"`
	Provider   *string   `json:"provider"`
	Role       string    `json:"role"`
	IsApproved bool      `json:"is_approved"`
	Watchlist  *string   `json:"watchlist"`
	CreatedAt  time.Time `json:"created_at"`
}

// UserSyncRequest mirrors schemas.UserSyncRequest.
type UserSyncRequest struct {
	Email    string  `json:"email"`
	Name     *string `json:"name"`
	Image    *string `json:"image"`
	Provider *string `json:"provider"`
}

// UserApprovalRequest mirrors schemas.UserApprovalRequest.
type UserApprovalRequest struct {
	IsApproved bool `json:"is_approved"`
}

// UserWatchlistUpdate mirrors schemas.UserWatchlistUpdate.
type UserWatchlistUpdate struct {
	Watchlist []string `json:"watchlist"`
}

// AuthGateStatus mirrors schemas.AuthGateStatus.
type AuthGateStatus struct {
	RequireLoginFullFeature bool `json:"require_login_full_feature"`
}

// AuthGateUpdateRequest mirrors schemas.AuthGateUpdateRequest.
type AuthGateUpdateRequest struct {
	RequireLoginFullFeature bool   `json:"require_login_full_feature"`
	AdminEmail              string `json:"admin_email,omitempty"`
}
