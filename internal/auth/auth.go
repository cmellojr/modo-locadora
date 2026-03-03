package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
)

const cookieName = "session_member"

// ErrInvalidSignature is returned when a cookie signature does not match.
var ErrInvalidSignature = errors.New("invalid cookie signature")

// SignCookie returns "value.hmac_hex" for the given value and secret.
func SignCookie(value, secret string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(value))
	sig := hex.EncodeToString(mac.Sum(nil))
	return value + "." + sig
}

// VerifyCookie splits a signed cookie, verifies the HMAC, and returns the
// original value. Returns ErrInvalidSignature if the signature is invalid.
func VerifyCookie(signed, secret string) (string, error) {
	idx := strings.LastIndex(signed, ".")
	if idx < 0 {
		return "", ErrInvalidSignature
	}

	value := signed[:idx]
	sig := signed[idx+1:]

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(value))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(sig), []byte(expected)) {
		return "", ErrInvalidSignature
	}

	return value, nil
}

// SetSessionCookie writes a signed session cookie with the member ID.
func SetSessionCookie(w http.ResponseWriter, memberID, secret string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieName,
		Value:    SignCookie(memberID, secret),
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}

// GetSessionMemberID extracts and verifies the member ID from the session
// cookie. Returns an empty string if the cookie is missing or invalid.
func GetSessionMemberID(r *http.Request, secret string) string {
	c, err := r.Cookie(cookieName)
	if err != nil {
		return ""
	}

	value, err := VerifyCookie(c.Value, secret)
	if err != nil {
		return ""
	}

	return value
}

// ClearSessionCookie removes the session cookie.
func ClearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})
}
