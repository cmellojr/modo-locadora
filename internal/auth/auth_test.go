package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSignAndVerifyCookie(t *testing.T) {
	secret := "secret-key-123"
	memberID := "123e4567-e89b-12d3-a456-426614174000"

	signed := SignCookie(memberID, secret)
	if signed == "" {
		t.Fatal("expected non-empty signed cookie string")
	}

	// Verify valid cookie
	val, err := VerifyCookie(signed, secret)
	if err != nil {
		t.Fatalf("expected valid verification, got error: %v", err)
	}
	if val != memberID {
		t.Fatalf("expected value %s, got %s", memberID, val)
	}

	// Verify with wrong secret
	_, err = VerifyCookie(signed, "wrong-secret")
	if err != ErrInvalidSignature {
		t.Fatalf("expected ErrInvalidSignature for wrong secret, got: %v", err)
	}

	// Verify with tampered value
	tampered := "99999999-e89b-12d3-a456-426614174000." + signed[len(memberID)+1:]
	_, err = VerifyCookie(tampered, secret)
	if err != ErrInvalidSignature {
		t.Fatalf("expected ErrInvalidSignature for tampered value, got: %v", err)
	}

	// Verify malformed cookie (no separator)
	_, err = VerifyCookie("invalid_format_cookie", secret)
	if err != ErrInvalidSignature {
		t.Fatalf("expected ErrInvalidSignature for malformed cookie, got: %v", err)
	}
}

func TestSessionCookieHandlers(t *testing.T) {
	secret := "test-secret"
	memberID := "member-uuid-456"

	// 1. Test SetSessionCookie
	w := httptest.NewRecorder()
	SetSessionCookie(w, memberID, secret)

	res := w.Result()
	cookies := res.Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected set cookie header")
	}

	var sessionCookie *http.Cookie
	for _, c := range cookies {
		if c.Name == cookieName {
			sessionCookie = c
			break
		}
	}
	if sessionCookie == nil {
		t.Fatal("session_member cookie not found in response")
	}

	// 2. Test GetSessionMemberID
	req := httptest.NewRequest("GET", "/", nil)
	req.AddCookie(sessionCookie)

	extracted := GetSessionMemberID(req, secret)
	if extracted != memberID {
		t.Fatalf("expected member ID %s, got %s", memberID, extracted)
	}

	// Test GetSessionMemberID with missing cookie
	emptyReq := httptest.NewRequest("GET", "/", nil)
	if id := GetSessionMemberID(emptyReq, secret); id != "" {
		t.Fatalf("expected empty member ID for missing cookie, got %s", id)
	}

	// 3. Test ClearSessionCookie
	clearW := httptest.NewRecorder()
	ClearSessionCookie(clearW)

	clearCookies := clearW.Result().Cookies()
	var cleared *http.Cookie
	for _, c := range clearCookies {
		if c.Name == cookieName {
			cleared = c
			break
		}
	}
	if cleared == nil {
		t.Fatal("cleared session_member cookie not found")
	}
	if cleared.MaxAge != -1 {
		t.Fatalf("expected MaxAge -1 for cleared cookie, got %d", cleared.MaxAge)
	}
}
