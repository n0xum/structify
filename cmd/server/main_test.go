package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestHandleHealth(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	s.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", resp["status"])
	}
}

func TestHandleVersion(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodGet, "/api/version", nil)
	w := httptest.NewRecorder()

	s.handleVersion(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if resp["version"] != version {
		t.Fatalf("expected version %q, got %q", version, resp["version"])
	}
}

func TestHandleGenerateSQL_EmptySource(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodPost, "/api/generate/sql", strings.NewReader(`{"source":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleGenerateSQL(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleGenerateSQL_InvalidJSON(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodPost, "/api/generate/sql", strings.NewReader(`not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleGenerateSQL(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestHandleGenerateSQL_ValidSource(t *testing.T) {
	s := newServer()
	source := `package models

type User struct {
	ID   int64  ` + "`db:\"pk\"`" + `
	Name string
}
`
	body := `{"source":` + toJSON(source) + `}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate/sql", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleGenerateSQL(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !strings.Contains(resp["output"], "CREATE TABLE") {
		t.Fatalf("expected SQL output, got %q", resp["output"])
	}
}

func TestHandleGenerateCode_DefaultPackage(t *testing.T) {
	s := newServer()
	source := `package models

type User struct {
	ID   int64  ` + "`db:\"pk\"`" + `
	Name string
}
`
	body := `{"source":` + toJSON(source) + `}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate/code", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleGenerateCode(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !strings.Contains(resp["output"], "package models") {
		t.Fatalf("expected package models in output, got %q", resp["output"])
	}
}

func TestHandleGenerateCode_EmptySource(t *testing.T) {
	s := newServer()
	req := httptest.NewRequest(http.MethodPost, "/api/generate/code", strings.NewReader(`{"source":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	s.handleGenerateCode(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestRateLimiter(t *testing.T) {
	rl := newRateLimiter(3, time.Minute)

	for i := 0; i < 3; i++ {
		if !rl.allow("1.2.3.4") {
			t.Fatalf("request %d should be allowed", i+1)
		}
	}

	if rl.allow("1.2.3.4") {
		t.Fatal("4th request should be denied")
	}

	// Different IP should still be allowed
	if !rl.allow("5.6.7.8") {
		t.Fatal("different IP should be allowed")
	}
}

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "https://example.com")

	s := newServer()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := corsMiddleware(s.limiter, mux)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "https://example.com" {
		t.Fatalf("expected CORS origin https://example.com, got %q", got)
	}
}

func TestCORSMiddleware_DisallowedOrigin(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "https://example.com")

	s := newServer()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /test", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	handler := corsMiddleware(s.limiter, mux)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "https://evil.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if got := w.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("expected no CORS header, got %q", got)
	}
}

func TestCORSMiddleware_Preflight(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "https://example.com")

	s := newServer()
	handler := corsMiddleware(s.limiter, http.NewServeMux())

	req := httptest.NewRequest(http.MethodOptions, "/api/generate/sql", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestClientIP_XForwardedFor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("X-Forwarded-For", "10.0.0.1, 10.0.0.2")

	if got := clientIP(req); got != "10.0.0.1" {
		t.Fatalf("expected 10.0.0.1, got %q", got)
	}
}

func TestClientIP_RemoteAddr(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Del("X-Forwarded-For")
	req.RemoteAddr = "192.168.1.1:12345"

	if got := clientIP(req); got != "192.168.1.1" {
		t.Fatalf("expected 192.168.1.1, got %q", got)
	}
}

func TestGetAllowedOrigins_Default(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "")

	origins := getAllowedOrigins()
	if len(origins) != 2 {
		t.Fatalf("expected 2 default origins, got %d", len(origins))
	}
	if origins[0] != "https://structify.alexander-kruska.dev" {
		t.Fatalf("expected structify origin, got %q", origins[0])
	}
}

func TestGetAllowedOrigins_Env(t *testing.T) {
	t.Setenv("ALLOWED_ORIGINS", "https://a.com,https://b.com")

	origins := getAllowedOrigins()
	if len(origins) != 2 {
		t.Fatalf("expected 2 origins, got %d", len(origins))
	}
}

func toJSON(s string) string {
	b, _ := json.Marshal(s)
	return string(b)
}
