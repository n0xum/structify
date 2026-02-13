package command

import (
	"context"

	"github.com/ak/structify/internal/domain/entity"
	"github.com/ak/structify/internal/domain/validator"
)

type Handler struct {
	generator Generator
}

type Generator interface {
	GenerateSchema(ctx context.Context, entities []*entity.Entity) (string, error)
	GenerateCode(ctx context.Context, packageName string, entities []*entity.Entity) (string, error)
}

func NewHandler(generator Generator) *Handler {
	return &Handler{
		generator: generator,
	}
}

type GenerateSchemaCommand struct {
	PackageName string
	Entities    []*entity.Entity
}

func (h *Handler) GenerateSchema(ctx context.Context, cmd *GenerateSchemaCommand) (string, error) {
	if err := h.validateEntities(cmd.Entities); err != nil {
		return "", err
	}
	return h.generator.GenerateSchema(ctx, cmd.Entities)
}

func (h *Handler) GenerateCode(ctx context.Context, cmd *GenerateSchemaCommand) (string, error) {
	if err := h.validateEntities(cmd.Entities); err != nil {
		return "", err
	}
	return h.generator.GenerateCode(ctx, cmd.PackageName, cmd.Entities)
}

func (h *Handler) validateEntities(entities []*entity.Entity) error {
	for _, ent := range entities {
		if _, err := validator.NewValidatedEntity(ent); err != nil {
			return err
		}
	}
	return nil
}

type ValidateCommand struct {
	Entities []*entity.Entity
}

func (h *Handler) Validate(ctx context.Context, cmd *ValidateCommand) error {
	return h.validateEntities(cmd.Entities)
}
