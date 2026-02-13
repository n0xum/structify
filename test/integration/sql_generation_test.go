//go:build integration

package integration

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ak/structify/internal/application"
	"github.com/ak/structify/internal/application/command"
	"github.com/ak/structify/internal/application/query"
	"github.com/ak/structify/internal/generator"
)

var fixturesDir = filepath.Join("..", "fixtures")
var expectedDir = filepath.Join("..", "expected")

// helpers ──────────────────────────────────────────────────────────────────

func newHandlers(t *testing.T) (*query.Handler, *command.Handler) {
	t.Helper()
	parserWrapper := application.NewParserWrapper()
	qh := query.NewHandler(parserWrapper)
	gen := generator.NewCompositeGenerator()
	ch := command.NewHandler(gen)
	return qh, ch
}

func parseFixture(t *testing.T, qh *query.Handler, file string) *query.ParseResult {
	t.Helper()
	ctx := context.Background()
	result, err := qh.Parse(ctx, &query.ParseQuery{Files: []string{file}})
	if err != nil {
		t.Fatalf("Parse(%s) error: %v", file, err)
	}
	if result.Count == 0 {
		t.Fatalf("Parse(%s) returned 0 structs", file)
	}
	return result
}

// SQL generation ───────────────────────────────────────────────────────────

func TestSQLGenerationProducesValidSQL(t *testing.T) {
	qh, ch := newHandlers(t)

	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateSchema(ctx, &command.GenerateSchemaCommand{
		Entities: result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateSchema error: %v", err)
	}

	if !strings.Contains(got, "CREATE TABLE") {
		t.Errorf("output must contain CREATE TABLE:\n%s", got)
	}
	if !strings.Contains(got, "PRIMARY KEY") {
		t.Errorf("output must contain PRIMARY KEY:\n%s", got)
	}
}

func TestSQLGenerationContainsCREATETABLE(t *testing.T) {
	qh, ch := newHandlers(t)
	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateSchema(ctx, &command.GenerateSchemaCommand{
		Entities: result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateSchema error: %v", err)
	}

	for _, table := range []string{`CREATE TABLE "user"`, `CREATE TABLE "product"`} {
		if !strings.Contains(got, table) {
			t.Errorf("expected %q in output:\n%s", table, got)
		}
	}
}

func TestSQLPrimaryKey(t *testing.T) {
	qh, ch := newHandlers(t)
	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateSchema(ctx, &command.GenerateSchemaCommand{
		Entities: result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateSchema error: %v", err)
	}

	if !strings.Contains(got, "BIGINT PRIMARY KEY") {
		t.Errorf("expected PRIMARY KEY constraint in:\n%s", got)
	}
}

func TestSQLUniqueConstraint(t *testing.T) {
	qh, ch := newHandlers(t)
	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateSchema(ctx, &command.GenerateSchemaCommand{
		Entities: result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateSchema error: %v", err)
	}

	if !strings.Contains(got, "UNIQUE") {
		t.Errorf("expected UNIQUE constraint in:\n%s", got)
	}
}

func TestSQLIgnoredFieldNotGenerated(t *testing.T) {
	qh, ch := newHandlers(t)
	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateSchema(ctx, &command.GenerateSchemaCommand{
		Entities: result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateSchema error: %v", err)
	}

	// InStock has db:"-" → must not appear as column
	if strings.Contains(got, "in_stock") {
		t.Errorf("ignored field 'in_stock' (db:\"-\") should not appear in SQL:\n%s", got)
	}
}

func TestSQLTypeMapping(t *testing.T) {
	qh, ch := newHandlers(t)
	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateSchema(ctx, &command.GenerateSchemaCommand{
		Entities: result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateSchema error: %v", err)
	}

	typeMappings := map[string]string{
		"BIGINT":           "int64 maps to BIGINT",
		"VARCHAR(255)":     "string maps to VARCHAR(255)",
		"BOOLEAN":          "bool maps to BOOLEAN",
		"DOUBLE PRECISION": "float64 maps to DOUBLE PRECISION",
	}

	for sqlType, description := range typeMappings {
		if !strings.Contains(got, sqlType) {
			t.Errorf("expected %s (%s) in:\n%s", sqlType, description, got)
		}
	}
}
