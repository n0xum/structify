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
