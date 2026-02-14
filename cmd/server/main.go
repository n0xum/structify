package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/n0xum/structify/internal/application"
	"github.com/n0xum/structify/internal/application/command"
	"github.com/n0xum/structify/internal/application/query"
	"github.com/n0xum/structify/internal/domain/entity"
	"github.com/n0xum/structify/internal/generator"
)

var version = "0.1.0"

func getAllowedOrigins() []string {
	if env := os.Getenv("ALLOWED_ORIGINS"); env != "" {
		return strings.Split(env, ",")
	}
	return []string{
		"https://n0xum.github.io",
		"http://localhost:3000",
	}
}

// rateLimiter tracks request counts per IP using a sliding window.
type rateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
	limit    int
	window   time.Duration
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (r *rateLimiter) allow(ip string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-r.window)

	reqs := r.requests[ip]
	var recent []time.Time
	for _, t := range reqs {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}

	if len(recent) >= r.limit {
		r.requests[ip] = recent
		return false
	}

	r.requests[ip] = append(recent, now)
	return true
}

type server struct {
	queryHandler *query.Handler
	cmdHandler   *command.Handler
	limiter      *rateLimiter
}

func newServer() *server {
	parserWrapper := application.NewParserWrapper()
	queryHandler := query.NewHandler(parserWrapper)
	compositeGen := generator.NewCompositeGenerator()
	cmdHandler := command.NewHandler(compositeGen)

	return &server{
		queryHandler: queryHandler,
		cmdHandler:   cmdHandler,
		limiter:      newRateLimiter(30, time.Minute),
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := newServer()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/version", s.handleVersion)
	mux.HandleFunc("POST /api/generate/sql", s.handleGenerateSQL)
	mux.HandleFunc("POST /api/generate/code", s.handleGenerateCode)

	handler := corsMiddleware(s.limiter, mux)

	log.Printf("structify server v%s listening on :%s", version, port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatal(err)
	}
}

// clientIP extracts the real client IP, respecting X-Forwarded-For behind a reverse proxy.
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		if ip := strings.TrimSpace(strings.SplitN(xff, ",", 2)[0]); ip != "" {
			return ip
		}
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func corsMiddleware(limiter *rateLimiter, next http.Handler) http.Handler {
	allowedOrigins := getAllowedOrigins()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		allowed := false
		for _, o := range allowedOrigins {
			if o == origin {
				allowed = true
				break
			}
		}

		if allowed {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ip := clientIP(r)
		if !limiter.allow(ip) {
			writeError(w, "rate limit exceeded — max 30 requests per minute", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) handleVersion(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, map[string]string{"version": version}, http.StatusOK)
}

func (s *server) handleGenerateSQL(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Source string `json:"source"`
	}
	if err := decodeBody(w, r, &req); err != nil {
		return
	}
	if req.Source == "" {
		writeError(w, "source is required", http.StatusBadRequest)
		return
	}

	entities, err := s.parseSource(r.Context(), req.Source)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	cmd := &command.GenerateSchemaCommand{Entities: entities}
	output, err := s.cmdHandler.GenerateSchema(r.Context(), cmd)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	writeJSON(w, map[string]string{"output": output}, http.StatusOK)
}

func (s *server) handleGenerateCode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Source  string `json:"source"`
		Package string `json:"package"`
	}
	if err := decodeBody(w, r, &req); err != nil {
		return
	}
	if req.Source == "" {
		writeError(w, "source is required", http.StatusBadRequest)
		return
	}
	if req.Package == "" {
		req.Package = "models"
	}

	entities, err := s.parseSource(r.Context(), req.Source)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	cmd := &command.GenerateSchemaCommand{PackageName: req.Package, Entities: entities}
	output, err := s.cmdHandler.GenerateCode(r.Context(), cmd)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	writeJSON(w, map[string]string{"output": output}, http.StatusOK)
}

// parseSource writes source to a temp file and parses it using the existing pipeline.
func (s *server) parseSource(ctx context.Context, source string) ([]*entity.Entity, error) {
	tmp, err := os.CreateTemp("", "structify_*.go")
	if err != nil {
		return nil, fmt.Errorf("internal error: could not create temp file")
	}
	name := tmp.Name()

	if _, err := tmp.WriteString(source); err != nil {
		tmp.Close()
		os.Remove(name)
		return nil, fmt.Errorf("internal error: could not write temp file")
	}
	tmp.Close()
	defer os.Remove(name)

	result, err := s.queryHandler.Parse(ctx, &query.ParseQuery{Files: []string{name}})
	if err != nil {
		return nil, err
	}
	if result.Count == 0 {
		return nil, fmt.Errorf("no exported structs found in the provided source")
	}

	return result.EntityList, nil
}

func decodeBody(w http.ResponseWriter, r *http.Request, dst any) error {
	if r.ContentLength > 500*1024 {
		writeError(w, "request body too large (max 500 KB)", http.StatusRequestEntityTooLarge)
		return fmt.Errorf("body too large")
	}
	r.Body = http.MaxBytesReader(w, r.Body, 500*1024)
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		writeError(w, "invalid JSON body", http.StatusBadRequest)
		return err
	}
	return nil
}

func writeJSON(w http.ResponseWriter, v any, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, map[string]string{"error": msg}, status)
}
