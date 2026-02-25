package generator

import (
	"context"

	"github.com/n0xum/structify/internal/domain/entity"
	"github.com/n0xum/structify/internal/generator/code"
	"github.com/n0xum/structify/internal/generator/sql"
)

type CompositeGenerator struct {
	sqlGenerator  *sql.SchemaGenerator
	codeGenerator *code.RepositoryGenerator
}

func NewCompositeGenerator() *CompositeGenerator {
	return &CompositeGenerator{
		sqlGenerator:  sql.NewSchemaGenerator(),
		codeGenerator: code.NewRepositoryGenerator(),
	}
}

func (g *CompositeGenerator) GenerateSchema(ctx context.Context, entities []*entity.Entity) (string, error) {
	return g.sqlGenerator.Generate(ctx, entities)
}

func (g *CompositeGenerator) GenerateCode(ctx context.Context, packageName string, entities []*entity.Entity) (string, error) {
	return g.codeGenerator.Generate(ctx, packageName, entities)
}

func (g *CompositeGenerator) GenerateRepository(ctx context.Context, packageName string, ent *entity.Entity, repo *entity.RepositoryInterface) (string, error) {
	return g.codeGenerator.GenerateFromInterface(ctx, packageName, ent, repo)
}
