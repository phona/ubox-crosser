package contract

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func stubAuthMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v1/auth/login", stubLoginHandler)
	mux.HandleFunc("/api/v1/services", stubAuthMiddleware(stubListServicesHandler))
	return mux
}

func stubLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, 9999, "method not allowed")
		return
	}
	var req LoginRequest
	if err := readJSON(r, &req); err != nil {
		writeError(w, http.StatusBadRequest, 2003, "invalid request body")
		return
	}
	if req.Username == "admin" && req.Password == "correct-password" {
		writeSuccess(w, http.StatusOK, LoginData{
			Token:     "eyJhbGciOiJIUzI1NiJ9.stub-token",
			ExpiresAt: time.Now().Add(time.Hour).Format(time.RFC3339),
		})
		return
	}
	writeError(w, http.StatusUnauthorized, 1001, "invalid credentials")
}

func stubAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" {
			writeError(w, http.StatusUnauthorized, 1003, "missing authorization")
			return
		}
		if auth == "Bearer expired-token" {
			writeError(w, http.StatusUnauthorized, 1002, "token expired")
			return
		}
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, http.StatusUnauthorized, 1003, "invalid authorization format")
			return
		}
		next(w, r)
	}
}

// REQ-997-S1: POST /api/v1/auth/login with valid credentials returns JWT
func TestAuthLoginSuccess(t *testing.T) {
	srv := httptest.NewServer(stubAuthMux())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodPost, "/api/v1/auth/login",
		jsonBody(t, LoginRequest{Username: "admin", Password: "correct-password"}), nil)

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	ct := resp.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("expected application/json, got %q", ct)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 0 {
		t.Fatalf("expected code 0, got %d", ar.Code)
	}

	data := parseData[LoginData](t, ar)
	if data.Token == "" {
		t.Error("expected non-empty token")
	}
	if data.ExpiresAt == "" {
		t.Error("expected non-empty expires_at")
	}
	if _, err := time.Parse(time.RFC3339, data.ExpiresAt); err != nil {
		t.Errorf("expires_at is not valid RFC3339: %v", err)
	}
}

// REQ-997-S2: POST /api/v1/auth/login with invalid credentials returns 401
func TestAuthLoginInvalidCredentials(t *testing.T) {
	srv := httptest.NewServer(stubAuthMux())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodPost, "/api/v1/auth/login",
		jsonBody(t, LoginRequest{Username: "admin", Password: "wrong"}), nil)

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}

	ar := parseResponse(t, resp)
	if ar.Code != 1001 {
		t.Errorf("expected error code 1001, got %d", ar.Code)
	}
	if ar.Message == "" {
		t.Error("expected non-empty error message")
	}
}

// REQ-997-S3: Authenticated endpoints reject requests without valid JWT
func TestAuthRejectsNoToken(t *testing.T) {
	srv := httptest.NewServer(stubAuthMux())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/v1/services", nil, nil)
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 without auth, got %d", resp.StatusCode)
	}
	ar := parseResponse(t, resp)
	if ar.Code != 1003 {
		t.Errorf("expected code 1003, got %d", ar.Code)
	}
}

func TestAuthRejectsExpiredToken(t *testing.T) {
	srv := httptest.NewServer(stubAuthMux())
	defer srv.Close()

	resp := doRequest(t, srv, http.MethodGet, "/api/v1/services", nil,
		map[string]string{"Authorization": "Bearer expired-token"})
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401 for expired token, got %d", resp.StatusCode)
	}
	ar := parseResponse(t, resp)
	if ar.Code != 1002 {
		t.Errorf("expected code 1002, got %d", ar.Code)
	}
}
