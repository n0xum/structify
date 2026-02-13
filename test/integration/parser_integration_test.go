//go:build integration

package integration

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/n0xum/structify/internal/application"
	"github.com/n0xum/structify/internal/application/query"
)

func TestFullParseUserFixture(t *testing.T) {
	parserWrapper := application.NewParserWrapper()
	qh := query.NewHandler(parserWrapper)

	ctx := context.Background()
	result, err := qh.Parse(ctx, &query.ParseQuery{
		Files: []string{filepath.Join(fixturesDir, "user.go")},
	})
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result.Count != 2 {
		t.Errorf("expected 2 structs (User, Product), got %d", result.Count)
	}

	entityNames := make(map[string]bool)
	for _, e := range result.EntityList {
		entityNames[e.Name] = true
	}

	for _, name := range []string{"User", "Product"} {
		if !entityNames[name] {
			t.Errorf("expected entity %q in parsed result", name)
		}
	}
}

func TestParseNonExistentFile(t *testing.T) {
	parserWrapper := application.NewParserWrapper()
	qh := query.NewHandler(parserWrapper)

	ctx := context.Background()
	_, err := qh.Parse(ctx, &query.ParseQuery{
		Files: []string{"nonexistent.go"},
	})
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}

func TestParsePackageName(t *testing.T) {
	parserWrapper := application.NewParserWrapper()
	qh := query.NewHandler(parserWrapper)

	ctx := context.Background()
	result, err := qh.Parse(ctx, &query.ParseQuery{
		Files: []string{filepath.Join(fixturesDir, "user.go")},
	})
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if result.Package == "" {
		t.Error("expected non-empty package name")
	}
}

func TestParsePrimaryKeyTagDetected(t *testing.T) {
	parserWrapper := application.NewParserWrapper()
	qh := query.NewHandler(parserWrapper)

	ctx := context.Background()
	result, err := qh.Parse(ctx, &query.ParseQuery{
		Files: []string{filepath.Join(fixturesDir, "user.go")},
	})
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	var user interface{ HasPrimaryKey() bool }
	for _, e := range result.EntityList {
		if e.Name == "User" {
			user = e
			break
		}
	}

	if user == nil {
		t.Fatal("User entity not found")
	}
	if !user.HasPrimaryKey() {
		t.Error("User entity should have a primary key (db:\"pk\")")
	}
}

func TestParseIgnoredFieldExcluded(t *testing.T) {
	parserWrapper := application.NewParserWrapper()
	qh := query.NewHandler(parserWrapper)

	ctx := context.Background()
	result, err := qh.Parse(ctx, &query.ParseQuery{
		Files: []string{filepath.Join(fixturesDir, "user.go")},
	})
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	for _, e := range result.EntityList {
		if e.Name == "Product" {
			for _, f := range e.GetGenerateableFields() {
				if f.Name == "InStock" {
					t.Error("InStock (db:\"-\") should not be in generateable fields")
				}
			}
		}
	}
}
