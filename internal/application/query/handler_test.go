package query

import (
	"context"
	"errors"
	"testing"

	"github.com/ak/structify/internal/domain/entity"
)

type mockParser struct {
	entities map[string][]*entity.Entity
	parseErr error
}

func (m *mockParser) ParseFiles(ctx context.Context, paths []string) (map[string][]*entity.Entity, error) {
	return m.entities, m.parseErr
}

func TestHandlerParse(t *testing.T) {
	t.Run("successful parse", func(t *testing.T) {
		entities := map[string][]*entity.Entity{
			"models": {
				{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}}},
				{Name: "Product", Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}}},
			},
		}

		parser := &mockParser{entities: entities}
		handler := NewHandler(parser)

		q := &ParseQuery{Files: []string{"user.go"}}
		result, err := handler.Parse(context.Background(), q)

		if err != nil {
			t.Fatalf("Parse() error = %v", err)
		}
		if result.Count != 2 {
			t.Errorf("Parse() Count = %d, want 2", result.Count)
		}
		if result.Package != "models" {
			t.Errorf("Parse() Package = %s, want models", result.Package)
		}
		if len(result.EntityList) != 2 {
			t.Errorf("Parse() EntityList length = %d, want 2", len(result.EntityList))
		}
	})

	t.Run("parse error", func(t *testing.T) {
		parseErr := errors.New("parse failed")
		parser := &mockParser{parseErr: parseErr}
		handler := NewHandler(parser)

		q := &ParseQuery{Files: []string{"invalid.go"}}
		_, err := handler.Parse(context.Background(), q)

		if err != parseErr {
			t.Errorf("Parse() error = %v, want %v", err, parseErr)
		}
	})
}

func TestHandlerFindEntity(t *testing.T) {
	t.Run("entity found", func(t *testing.T) {
		parser := &mockParser{}
		handler := NewHandler(parser)

		entities := map[string][]*entity.Entity{
			"models": {
				{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64"}}},
			},
		}

		q := &FindEntityQuery{Name: "User", Entities: entities}
		ent, found := handler.FindEntity(context.Background(), q)

		if !found {
			t.Error("FindEntity() found = false, want true")
		}
		if ent == nil {
			t.Fatal("FindEntity() returned nil entity")
		}
		if ent.Name != "User" {
			t.Errorf("FindEntity() entity.Name = %s, want User", ent.Name)
		}
	})

	t.Run("entity not found", func(t *testing.T) {
		parser := &mockParser{}
		handler := NewHandler(parser)

		entities := map[string][]*entity.Entity{
			"models": {
				{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64"}}},
			},
		}

		q := &FindEntityQuery{Name: "Product", Entities: entities}
		ent, found := handler.FindEntity(context.Background(), q)

		if found {
			t.Error("FindEntity() found = true, want false")
		}
		if ent != nil {
			t.Error("FindEntity() returned non-nil entity when not found")
		}
	})
}

func TestHandlerListEntities(t *testing.T) {
	parser := &mockParser{}
	handler := NewHandler(parser)

	entities := map[string][]*entity.Entity{
		"models": {
			{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64"}}},
			{Name: "Product", Fields: []entity.Field{{Name: "ID", Type: "int64"}}},
		},
		"admin": {
			{Name: "Admin", Fields: []entity.Field{{Name: "ID", Type: "int64"}}},
		},
	}

	q := &ListEntitiesQuery{Entities: entities}
	result := handler.ListEntities(context.Background(), q)

	if len(result) != 3 {
		t.Errorf("ListEntities() length = %d, want 3", len(result))
	}
}
