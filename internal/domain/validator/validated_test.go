package validator

import (
	"testing"

	"github.com/n0xum/structify/internal/domain/entity"
)

func TestNewValidatedEntity(t *testing.T) {
	t.Run("valid entity", func(t *testing.T) {
		ent := &entity.Entity{
			Name: "User",
			Fields: []entity.Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
			},
		}

		ve, err := NewValidatedEntity(ent)
		if err != nil {
			t.Fatalf("NewValidatedEntity() error = %v", err)
		}
		if !ve.IsValid() {
			t.Error("IsValid() returned false")
		}
		if ve.Get() != ent {
			t.Error("Get() returned wrong entity")
		}
	})

	t.Run("invalid entity", func(t *testing.T) {
		ent := &entity.Entity{
			Name: "",
		}

		_, err := NewValidatedEntity(ent)
		if err == nil {
			t.Fatal("NewValidatedEntity() expected error")
		}
	})
}
