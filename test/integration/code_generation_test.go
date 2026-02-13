//go:build integration

package integration

import (
	"context"
	"path/filepath"
	"strings"
	"testing"

	"github.com/n0xum/structify/internal/application/command"
)

func TestCodeGenerationContainsCRUD(t *testing.T) {
	qh, ch := newHandlers(t)
	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateCode(ctx, &command.GenerateSchemaCommand{
		PackageName: result.Package,
		Entities:    result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateCode error: %v", err)
	}

	crudFunctions := []string{"Create", "Get", "Update", "Delete", "List"}
	for _, fn := range crudFunctions {
		if !strings.Contains(got, fn) {
			t.Errorf("expected CRUD function %q in generated code:\n%s", fn, got)
		}
	}
}

func TestCodeGenerationImportsSQL(t *testing.T) {
	qh, ch := newHandlers(t)
	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateCode(ctx, &command.GenerateSchemaCommand{
		PackageName: result.Package,
		Entities:    result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateCode error: %v", err)
	}

	if !strings.Contains(got, `"database/sql"`) {
		t.Errorf("expected database/sql import in generated code")
	}
	if !strings.Contains(got, `"context"`) {
		t.Errorf("expected context import in generated code")
	}
}

func TestCodeGenerationUsesContextParam(t *testing.T) {
	qh, ch := newHandlers(t)
	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateCode(ctx, &command.GenerateSchemaCommand{
		PackageName: result.Package,
		Entities:    result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateCode error: %v", err)
	}

	if !strings.Contains(got, "context.Context") {
		t.Errorf("generated code must accept context.Context in function signatures")
	}
}

func TestCodeGenerationUsesParameterizedQueries(t *testing.T) {
	qh, ch := newHandlers(t)
	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateCode(ctx, &command.GenerateSchemaCommand{
		PackageName: result.Package,
		Entities:    result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateCode error: %v", err)
	}

	// Must use $1 style placeholders (PostgreSQL prepared statements)
	if !strings.Contains(got, "$1") {
		t.Errorf("generated code must use parameterized queries ($1, $2, ...)")
	}
}

func TestCodeGenerationPerEntity(t *testing.T) {
	qh, ch := newHandlers(t)
	result := parseFixture(t, qh, filepath.Join(fixturesDir, "user.go"))

	ctx := context.Background()
	got, err := ch.GenerateCode(ctx, &command.GenerateSchemaCommand{
		PackageName: result.Package,
		Entities:    result.EntityList,
	})
	if err != nil {
		t.Fatalf("GenerateCode error: %v", err)
	}

	// Both User and Product structs must be present
	for _, entity := range []string{"User", "Product"} {
		if !strings.Contains(got, entity) {
			t.Errorf("expected entity %q in generated code", entity)
		}
	}
}
