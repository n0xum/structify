package query

import (
	"context"

	"github.com/n0xum/structify/internal/domain/entity"
)

type Handler struct {
	parser Parser
}

type Parser interface {
	ParseFiles(ctx context.Context, paths []string) (map[string][]*entity.Entity, error)
}

func NewHandler(parser Parser) *Handler {
	return &Handler{
		parser: parser,
	}
}

type ParseQuery struct {
	Files []string
}

type ParseResult struct {
	Entities   map[string][]*entity.Entity
	Package    string
	Count      int
	EntityList []*entity.Entity
}

func (h *Handler) Parse(ctx context.Context, q *ParseQuery) (*ParseResult, error) {
	entities, err := h.parser.ParseFiles(ctx, q.Files)
	if err != nil {
		return nil, err
	}

	var allEntities []*entity.Entity
	var pkgName string
	for pkg, pkgEntities := range entities {
		pkgName = pkg
		allEntities = append(allEntities, pkgEntities...)
	}

	return &ParseResult{
		Entities:   entities,
		Package:    pkgName,
		Count:      len(allEntities),
		EntityList: allEntities,
	}, nil
}

type FindEntityQuery struct {
	Name      string
	Entities  map[string][]*entity.Entity
}

func (h *Handler) FindEntity(ctx context.Context, q *FindEntityQuery) (*entity.Entity, bool) {
	for _, entities := range q.Entities {
		for _, ent := range entities {
			if ent.Name == q.Name {
				return ent, true
			}
		}
	}
	return nil, false
}

type ListEntitiesQuery struct {
	Entities map[string][]*entity.Entity
}

func (h *Handler) ListEntities(ctx context.Context, q *ListEntitiesQuery) []*entity.Entity {
	var result []*entity.Entity
	for _, entities := range q.Entities {
		result = append(result, entities...)
	}
	return result
}
