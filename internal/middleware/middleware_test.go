package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cmellojr/modo-locadora/internal/auth"
	"github.com/cmellojr/modo-locadora/internal/database"
	"github.com/cmellojr/modo-locadora/internal/models"
	"github.com/google/uuid"
)

type mockStore struct {
	database.Store
	member *models.Member
	err    error
}

func (m *mockStore) GetMemberByID(ctx context.Context, id uuid.UUID) (*models.Member, error) {
	return m.member, m.err
}

func TestRequireAuth(t *testing.T) {
	secret := "test-secret"
	validMemberID := "123e4567-e89b-12d3-a456-426614174000"

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		memberID := MemberIDFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(memberID))
	})

	handler := RequireAuth(secret, nextHandler)

	// 1. Missing cookie -> StatusSeeOther
	req1 := httptest.NewRequest("GET", "/protected", nil)
	w1 := httptest.NewRecorder()
	handler(w1, req1)

	if w1.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, w1.Code)
	}

	// 2. Invalid cookie -> StatusSeeOther
	req2 := httptest.NewRequest("GET", "/protected", nil)
	req2.AddCookie(&http.Cookie{Name: "session_member", Value: "invalid_signed_cookie"})
	w2 := httptest.NewRecorder()
	handler(w2, req2)

	if w2.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d, got %d", http.StatusSeeOther, w2.Code)
	}

	// 3. Valid cookie -> StatusOK and member ID in context
	req3 := httptest.NewRequest("GET", "/protected", nil)
	signedVal := auth.SignCookie(validMemberID, secret)
	req3.AddCookie(&http.Cookie{Name: "session_member", Value: signedVal})
	w3 := httptest.NewRecorder()
	handler(w3, req3)

	if w3.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, w3.Code)
	}
	if body := w3.Body.String(); body != validMemberID {
		t.Fatalf("expected body %s, got %s", validMemberID, body)
	}
}

func TestRequireAdmin(t *testing.T) {
	secret := "test-secret"
	adminEmail := "tio@locadora.com"
	adminUUID := uuid.New()

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("admin_ok"))
	})

	// 1. No cookie -> redirect
	handlerNilStore := RequireAdmin(secret, adminEmail, nil, nextHandler)
	req1 := httptest.NewRequest("GET", "/admin", nil)
	w1 := httptest.NewRecorder()
	handlerNilStore(w1, req1)
	if w1.Code != http.StatusSeeOther {
		t.Fatalf("expected status %d for missing cookie, got %d", http.StatusSeeOther, w1.Code)
	}

	// 2. Valid cookie but store is nil -> StatusServiceUnavailable
	req2 := httptest.NewRequest("GET", "/admin", nil)
	req2.AddCookie(&http.Cookie{Name: "session_member", Value: auth.SignCookie(adminUUID.String(), secret)})
	w2 := httptest.NewRecorder()
	handlerNilStore(w2, req2)
	if w2.Code != http.StatusServiceUnavailable {
		t.Fatalf("expected status %d for nil store, got %d", http.StatusServiceUnavailable, w2.Code)
	}

	// 3. Member is not admin -> StatusForbidden
	regularStore := &mockStore{
		member: &models.Member{
			ID:    adminUUID,
			Email: "regular@user.com",
		},
	}
	handlerRegular := RequireAdmin(secret, adminEmail, regularStore, nextHandler)
	req3 := httptest.NewRequest("GET", "/admin", nil)
	req3.AddCookie(&http.Cookie{Name: "session_member", Value: auth.SignCookie(adminUUID.String(), secret)})
	w3 := httptest.NewRecorder()
	handlerRegular(w3, req3)
	if w3.Code != http.StatusForbidden {
		t.Fatalf("expected status %d for non-admin, got %d", http.StatusForbidden, w3.Code)
	}

	// 4. Member is admin -> StatusOK
	adminStore := &mockStore{
		member: &models.Member{
			ID:    adminUUID,
			Email: adminEmail,
		},
	}
	handlerAdmin := RequireAdmin(secret, adminEmail, adminStore, nextHandler)
	req4 := httptest.NewRequest("GET", "/admin", nil)
	req4.AddCookie(&http.Cookie{Name: "session_member", Value: auth.SignCookie(adminUUID.String(), secret)})
	w4 := httptest.NewRecorder()
	handlerAdmin(w4, req4)
	if w4.Code != http.StatusOK {
		t.Fatalf("expected status %d for admin, got %d", http.StatusOK, w4.Code)
	}
}
