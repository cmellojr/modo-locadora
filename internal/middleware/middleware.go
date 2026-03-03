package middleware

import (
	"context"
	"net/http"

	"github.com/cmellojr/modo-locadora/internal/auth"
	"github.com/cmellojr/modo-locadora/internal/database"
	"github.com/google/uuid"
)

type contextKey string

const memberIDKey contextKey = "member_id"

// MemberIDFromContext extracts the authenticated member ID from the request context.
func MemberIDFromContext(ctx context.Context) string {
	v, _ := ctx.Value(memberIDKey).(string)
	return v
}

// RequireAuth rejects requests without a valid signed session cookie.
func RequireAuth(secret string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memberID := auth.GetSessionMemberID(r, secret)
		if memberID == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), memberIDKey, memberID)
		next(w, r.WithContext(ctx))
	}
}

// RequireAdmin rejects requests unless the authenticated member's email
// matches the configured admin email.
func RequireAdmin(secret, adminEmail string, store database.Store, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		memberID := auth.GetSessionMemberID(r, secret)
		if memberID == "" {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if store == nil {
			http.Error(w, "Database not configured", http.StatusServiceUnavailable)
			return
		}

		id, err := uuid.Parse(memberID)
		if err != nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		member, err := store.GetMemberByID(r.Context(), id)
		if err != nil || member == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if member.Email != adminEmail {
			http.Error(w, "Acesso restrito ao Tio da Locadora", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), memberIDKey, memberID)
		next(w, r.WithContext(ctx))
	}
}
