package command

import (
	"context"
	"errors"
	"testing"

	"github.com/n0xum/structify/internal/domain/entity"
)

type mockGenerator struct {
	schemaResult string
	codeResult   string
	schemaError  error
	codeError    error
}

func (m *mockGenerator) GenerateSchema(ctx context.Context, entities []*entity.Entity) (string, error) {
	return m.schemaResult, m.schemaError
}

func (m *mockGenerator) GenerateCode(ctx context.Context, packageName string, entities []*entity.Entity) (string, error) {
	return m.codeResult, m.codeError
}

func (m *mockGenerator) GenerateRepository(ctx context.Context, packageName string, ent *entity.Entity, repo *entity.RepositoryInterface) (string, error) {
	return m.codeResult, m.codeError
}

func TestHandlerGenerateSchema(t *testing.T) {
	t.Run("valid entities", func(t *testing.T) {
		gen := &mockGenerator{schemaResult: "CREATE TABLE..."}
		handler := NewHandler(gen)

		entities := []*entity.Entity{
			{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}}},
		}

		cmd := &GenerateSchemaCommand{Entities: entities}
		result, err := handler.GenerateSchema(context.Background(), cmd)

		if err != nil {
			t.Fatalf("GenerateSchema() error = %v", err)
		}
		if result != "CREATE TABLE..." {
			t.Errorf("GenerateSchema() result = %v, want CREATE TABLE...", result)
		}
	})

	t.Run("invalid entity", func(t *testing.T) {
		gen := &mockGenerator{}
		handler := NewHandler(gen)

		entities := []*entity.Entity{
			{Name: "", Fields: []entity.Field{}},
		}

		cmd := &GenerateSchemaCommand{Entities: entities}
		_, err := handler.GenerateSchema(context.Background(), cmd)

		if err == nil {
			t.Error("GenerateSchema() expected error for invalid entity")
		}
	})
}

func TestHandlerGenerateCode(t *testing.T) {
	t.Run("valid entities", func(t *testing.T) {
		gen := &mockGenerator{codeResult: "package main..."}
		handler := NewHandler(gen)

		entities := []*entity.Entity{
			{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}}},
		}

		cmd := &GenerateSchemaCommand{PackageName: "models", Entities: entities}
		result, err := handler.GenerateCode(context.Background(), cmd)

		if err != nil {
			t.Fatalf("GenerateCode() error = %v", err)
		}
		if result != "package main..." {
			t.Errorf("GenerateCode() result = %v, want package main...", result)
		}
	})
}

func TestHandlerValidate(t *testing.T) {
	t.Run("valid entities", func(t *testing.T) {
		gen := &mockGenerator{}
		handler := NewHandler(gen)

		entities := []*entity.Entity{
			{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}}},
		}

		cmd := &ValidateCommand{Entities: entities}
		err := handler.Validate(context.Background(), cmd)

		if err != nil {
			t.Errorf("Validate() error = %v", err)
		}
	})

	t.Run("invalid entity", func(t *testing.T) {
		gen := &mockGenerator{}
		handler := NewHandler(gen)

		entities := []*entity.Entity{
			{Name: "", Fields: []entity.Field{}},
		}

		cmd := &ValidateCommand{Entities: entities}
		err := handler.Validate(context.Background(), cmd)

		if err == nil {
			t.Error("Validate() expected error for invalid entity")
		}
	})
}

func TestHandlerGeneratorError(t *testing.T) {
	genErr := errors.New("generator failed")
	gen := &mockGenerator{schemaError: genErr}
	handler := NewHandler(gen)

	entities := []*entity.Entity{
		{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}}},
	}

	cmd := &GenerateSchemaCommand{Entities: entities}
	_, err := handler.GenerateSchema(context.Background(), cmd)

	if err != genErr {
		t.Errorf("GenerateSchema() error = %v, want %v", err, genErr)
	}
}

func TestHandlerGenerateRepository(t *testing.T) {
	t.Run("valid entity and interface", func(t *testing.T) {
		gen := &mockGenerator{codeResult: "package repository..."}
		handler := NewHandler(gen)

		ent := &entity.Entity{
			Name:   "User",
			Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}},
		}

		repo := &entity.RepositoryInterface{
			Name:       "UserRepository",
			EntityName: "User",
			Package:    "repository",
			Methods: []entity.RepositoryMethod{
				{Name: "GetByID", Kind: entity.MethodGetByID},
			},
		}

		cmd := &GenerateRepoCommand{
			Entity:      ent,
			Interface:   repo,
			PackageName: "repository",
		}

		result, err := handler.GenerateRepository(context.Background(), cmd)

		if err != nil {
			t.Fatalf("GenerateRepository() error = %v", err)
		}
		if result != "package repository..." {
			t.Errorf("GenerateRepository() result = %v, want package repository...", result)
		}
	})

	t.Run("invalid entity", func(t *testing.T) {
		gen := &mockGenerator{}
		handler := NewHandler(gen)

		ent := &entity.Entity{
			Name:   "",
			Fields: []entity.Field{},
		}

		repo := &entity.RepositoryInterface{
			Name:       "UserRepository",
			EntityName: "User",
			Package:    "repository",
		}

		cmd := &GenerateRepoCommand{
			Entity:      ent,
			Interface:   repo,
			PackageName: "repository",
		}

		_, err := handler.GenerateRepository(context.Background(), cmd)

		if err == nil {
			t.Error("GenerateRepository() expected error for invalid entity")
		}
	})

	t.Run("generator error", func(t *testing.T) {
		genErr := errors.New("generator failed")
		gen := &mockGenerator{codeError: genErr}
		handler := NewHandler(gen)

		ent := &entity.Entity{
			Name:   "User",
			Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}},
		}

		repo := &entity.RepositoryInterface{
			Name:       "UserRepository",
			EntityName: "User",
			Package:    "repository",
		}

		cmd := &GenerateRepoCommand{
			Entity:      ent,
			Interface:   repo,
			PackageName: "repository",
		}

		_, err := handler.GenerateRepository(context.Background(), cmd)

		if err != genErr {
			t.Errorf("GenerateRepository() error = %v, want %v", err, genErr)
		}
	})
}

func TestHandlerGenerateCode_InvalidEntity(t *testing.T) {
	gen := &mockGenerator{}
	handler := NewHandler(gen)

	entities := []*entity.Entity{
		{Name: "", Fields: []entity.Field{}},
	}

	cmd := &GenerateSchemaCommand{PackageName: "models", Entities: entities}
	_, err := handler.GenerateCode(context.Background(), cmd)

	if err == nil {
		t.Error("GenerateCode() expected error for invalid entity")
	}
}

func TestHandlerGenerateCodeError(t *testing.T) {
	genErr := errors.New("code generation failed")
	gen := &mockGenerator{codeError: genErr}
	handler := NewHandler(gen)

	entities := []*entity.Entity{
		{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}}},
	}

	cmd := &GenerateSchemaCommand{PackageName: "models", Entities: entities}
	_, err := handler.GenerateCode(context.Background(), cmd)

	if err != genErr {
		t.Errorf("GenerateCode() error = %v, want %v", err, genErr)
	}
}
