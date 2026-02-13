package code

import (
	"context"
	"strings"
	"testing"

	"github.com/ak/structify/internal/domain/entity"
)

func TestRepositoryGeneratorGenerate(t *testing.T) {
	gen := NewRepositoryGenerator()

	entities := []*entity.Entity{
		{
			Name:    "User",
			Package: "models",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string", IsUnique: true},
				{Name: "Email", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "package models") {
		t.Error("Generate() result missing package declaration")
	}
	if !strings.Contains(result, "import") {
		t.Error("Generate() result missing import block")
	}
	if !strings.Contains(result, "database/sql") {
		t.Error("Generate() result missing database/sql import")
	}
	if !strings.Contains(result, "type User struct") {
		t.Error("Generate() result missing User struct")
	}
	if !strings.Contains(result, "func CreateUser") {
		t.Error("Generate() result missing CreateUser function")
	}
	if !strings.Contains(result, "func GetUserByID") {
		t.Error("Generate() result missing GetUserByID function")
	}
	if !strings.Contains(result, "func UpdateUser") {
		t.Error("Generate() result missing UpdateUser function")
	}
	if !strings.Contains(result, "func DeleteUser") {
		t.Error("Generate() result missing DeleteUser function")
	}
	if !strings.Contains(result, "func ListUser") {
		t.Error("Generate() result missing ListUser function")
	}
}

func TestRepositoryGeneratorGenerateWithIgnoredField(t *testing.T) {
	gen := NewRepositoryGenerator()

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

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if strings.Contains(result, "Password string") {
		t.Error("Generate() result contains ignored Password field in struct")
	}
}

func TestRepositoryGeneratorCreateMethodUsesCorrectAPI(t *testing.T) {
	gen := NewRepositoryGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, ".Scan(&id)") {
		t.Error("Generate() Create method missing .Scan() call")
	}
	if !strings.Contains(result, "QueryRowContext(ctx, query,") {
		t.Error("Generate() Create method using wrong API")
	}
}

func TestRepositoryGeneratorGetByIDMethodUsesCorrectAPI(t *testing.T) {
	gen := NewRepositoryGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, ".Scan(") {
		t.Error("Generate() GetByID method missing .Scan() call")
	}
	if !strings.Contains(result, "QueryRowContext(ctx, query, id)") {
		t.Error("Generate() GetByID method using wrong API")
	}
}

func TestRepositoryGeneratorListMethodUsesCorrectAPI(t *testing.T) {
	gen := NewRepositoryGenerator()

	entities := []*entity.Entity{
		{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string"},
			},
		},
	}

	result, err := gen.Generate(context.Background(), "models", entities)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if !strings.Contains(result, "rows.Scan(") {
		t.Error("Generate() List method missing .Scan() call")
	}
	if !strings.Contains(result, "defer rows.Close()") {
		t.Error("Generate() List method missing defer rows.Close()")
	}
}
