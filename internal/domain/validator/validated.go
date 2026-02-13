package validator

import "github.com/ak/structify/internal/domain/entity"

type ValidatedEntity struct {
	Entity         *entity.Entity
	IsValidated   bool
}

func NewValidatedEntity(ent *entity.Entity) (*ValidatedEntity, error) {
	if err := ent.Validate(); err != nil {
		return nil, err
	}
	return &ValidatedEntity{
		Entity:       ent,
		IsValidated: true,
	}, nil
}

func (v *ValidatedEntity) Get() *entity.Entity {
	return v.Entity
}

func (v *ValidatedEntity) IsValid() bool {
	return v.IsValidated
}
