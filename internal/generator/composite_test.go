package generator

import (
	"context"
	"testing"

	"github.com/n0xum/structify/internal/domain/entity"
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

func TestCompositeGeneratorGenerateRepository(t *testing.T) {
	t.Run("valid entity and interface", func(t *testing.T) {
		gen := NewCompositeGenerator()
		ctx := context.Background()

		ent := &entity.Entity{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Username", Type: "string"},
			},
		}

		repo := &entity.RepositoryInterface{
			Name:       "UserRepository",
			EntityName: "User",
			Package:    "repository",
			Methods: []entity.RepositoryMethod{
				{Name: "GetByID", Kind: entity.MethodGetByID},
				{Name: "Create", Kind: entity.MethodCreate},
			},
		}

		result, err := gen.GenerateRepository(ctx, "repository", ent, repo)
		if err != nil {
			t.Fatalf("GenerateRepository() error = %v", err)
		}

		if result == "" {
			t.Error("GenerateRepository() returned empty string")
		}

		if !contains(result, "package repository") {
			t.Error("GenerateRepository() missing package declaration")
		}

		if !contains(result, "type UserRepositoryImpl struct") {
			t.Error("GenerateRepository() missing impl struct")
		}

		if !contains(result, "func NewUserRepository") {
			t.Error("GenerateRepository() missing constructor")
		}

		if !contains(result, "func (r *UserRepositoryImpl)") {
			t.Error("GenerateRepository() missing receiver methods")
		}
	})

	t.Run("empty package name uses main", func(t *testing.T) {
		gen := NewCompositeGenerator()
		ctx := context.Background()

		ent := &entity.Entity{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
			},
		}

		repo := &entity.RepositoryInterface{
			Name:       "UserRepository",
			EntityName: "User",
			Package:    "repository",
			Methods:    []entity.RepositoryMethod{},
		}

		result, err := gen.GenerateRepository(ctx, "", ent, repo)
		if err != nil {
			t.Fatalf("GenerateRepository() error = %v", err)
		}

		if !contains(result, "package main") {
			t.Error("GenerateRepository() should use 'main' package when empty")
		}
	})

	t.Run("all method kinds", func(t *testing.T) {
		gen := NewCompositeGenerator()
		ctx := context.Background()

		ent := &entity.Entity{
			Name: "Product",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
				{Name: "Name", Type: "string"},
				{Name: "Price", Type: "float64"},
			},
		}

		repo := &entity.RepositoryInterface{
			Name:       "ProductRepository",
			EntityName: "Product",
			Package:    "repository",
			Methods: []entity.RepositoryMethod{
				{Name: "Create", Kind: entity.MethodCreate},
				{Name: "GetByID", Kind: entity.MethodGetByID},
				{Name: "Update", Kind: entity.MethodUpdate},
				{Name: "Delete", Kind: entity.MethodDelete},
				{Name: "List", Kind: entity.MethodList},
				{Name: "FindByName", Kind: entity.MethodFindBy},
				{Name: "CustomQuery", Kind: entity.MethodCustomSQL, CustomSQL: "SELECT * FROM products WHERE price > $1"},
			},
		}

		result, err := gen.GenerateRepository(ctx, "repo", ent, repo)
		if err != nil {
			t.Fatalf("GenerateRepository() error = %v", err)
		}

		// Verify all method types are generated
		if !contains(result, "func (r *ProductRepositoryImpl) Create") {
			t.Error("GenerateRepository() missing Create method")
		}
		if !contains(result, "func (r *ProductRepositoryImpl) GetByID") {
			t.Error("GenerateRepository() missing GetByID method")
		}
		if !contains(result, "func (r *ProductRepositoryImpl) Update") {
			t.Error("GenerateRepository() missing Update method")
		}
		if !contains(result, "func (r *ProductRepositoryImpl) Delete") {
			t.Error("GenerateRepository() missing Delete method")
		}
		if !contains(result, "func (r *ProductRepositoryImpl) List") {
			t.Error("GenerateRepository() missing List method")
		}
		if !contains(result, "func (r *ProductRepositoryImpl) FindByName") {
			t.Error("GenerateRepository() missing FindByName method")
		}
		if !contains(result, "func (r *ProductRepositoryImpl) CustomQuery") {
			t.Error("GenerateRepository() missing CustomQuery method")
		}
		if !contains(result, "SELECT * FROM products WHERE price > $1") {
			t.Error("GenerateRepository() missing custom SQL")
		}
	})
}
