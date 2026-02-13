package sql

import (
	"context"
	"strings"
	"testing"

	"github.com/ak/structify/internal/domain/entity"
)

func TestSchemaGeneratorGenerate(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string", IsUnique: true},
				{Name: "Email", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	t.Logf("Generated SQL:\n%s", result)

	if !strings.Contains(result, "CREATE TABLE") {
		t.Error("Generate() result missing CREATE TABLE")
	}
	if !strings.Contains(result, "user") {
		t.Error("Generate() result missing table name")
	}
	if !strings.Contains(result, "id") {
		t.Error("Generate() result missing id column")
	}
	if !strings.Contains(result, "PRIMARY KEY") {
		t.Error("Generate() result missing PRIMARY KEY")
	}
}

func TestSchemaGeneratorGenerateWithIgnoredField(t *testing.T) {
	gen := NewSchemaGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string"},
				{Name: "Password", Type: "string", IsIgnored: true},
			},
		},
	}

	result, err := gen.Generate(context.Background(), entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if strings.Contains(result, "password") {
		t.Error("Generate() result contains ignored field password")
	}
}
