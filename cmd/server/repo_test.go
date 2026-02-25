package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleGenerateRepo_ValidSource(t *testing.T) {
	s := newServer()
	source := `package models

import "context"

type User struct {
	ID   int64  ` + "`db:\"pk\"`" + `
	Name string
}

type UserRepository interface {
	FindByID(ctx context.Context, id int64) (*User, error)
}
`
	body := `{"source":` + toJSON(source) + `,"package":"repository"}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate/repo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// This should fail initially because the route is not defined
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/generate/repo", s.handleGenerateRepo)
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d; body: %s", w.Code, w.Body.String())
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if !strings.Contains(resp["output"], "type UserRepositoryImpl struct") {
		t.Fatalf("expected repository implementation, got %q", resp["output"])
	}
	if !strings.Contains(resp["output"], "package repository") {
		t.Fatalf("expected package repository, got %q", resp["output"])
	}
}

func TestHandleGenerateRepo_NoInterface(t *testing.T) {
	s := newServer()
	source := `package models

type User struct {
	ID   int64  ` + "`db:\"pk\"`" + `
	Name string
}
`
	body := `{"source":` + toJSON(source) + `}`
	req := httptest.NewRequest(http.MethodPost, "/api/generate/repo", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/generate/repo", s.handleGenerateRepo)
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d; body: %s", w.Code, w.Body.String())
	}
}
