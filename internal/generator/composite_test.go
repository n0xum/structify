package generator

import (
	"context"
	"testing"

	"github.com/ak/structify/internal/domain/entity"
)

func TestNewCompositeGenerator(t *testing.T) {
	gen := NewCompositeGenerator()

	if gen == nil {
		t.Fatal("NewCompositeGenerator() returned nil")
	}
	if gen.sqlGenerator == nil {
		t.Error("sqlGenerator is nil")
	}
	if gen.codeGenerator == nil {
		t.Error("codeGenerator is nil")
	}
}

func TestCompositeGeneratorGenerateSchema(t *testing.T) {
	gen := NewCompositeGenerator()

	ctx := context.Background()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string", IsUnique: true},
			},
		},
	}

	result, err := gen.GenerateSchema(ctx, entities)
	if err != nil {
		t.Fatalf("GenerateSchema() error = %v", err)
	}

	if result == "" {
		t.Error("GenerateSchema() returned empty string")
	}

	if !contains(result, "CREATE TABLE") {
		t.Error("GenerateSchema() missing CREATE TABLE")
	}

	if !contains(result, "user") {
		t.Error("GenerateSchema() missing table name")
	}
}

func TestCompositeGeneratorGenerateCode(t *testing.T) {
	gen := NewCompositeGenerator()

	ctx := context.Background()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string"},
			},
		},
	}

	result, err := gen.GenerateCode(ctx, "models", entities)
	if err != nil {
		t.Fatalf("GenerateCode() error = %v", err)
	}

	if result == "" {
		t.Error("GenerateCode() returned empty string")
	}

	if !contains(result, "package models") {
		t.Error("GenerateCode() missing package declaration")
	}

	if !contains(result, "type User struct") {
		t.Error("GenerateCode() missing User struct")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
