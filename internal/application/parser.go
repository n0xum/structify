package application

import (
	"context"

	"github.com/ak/structify/internal/adapter"
	"github.com/ak/structify/internal/domain/entity"
	"github.com/ak/structify/internal/parser"
)

type ParserWrapper struct {
	parser   *parser.Parser
	adapter  *adapter.ParserAdapter
}

func NewParserWrapper() *ParserWrapper {
	return &ParserWrapper{
		parser:  parser.New(),
		adapter: adapter.NewParserAdapter(),
	}
}

func (p *ParserWrapper) ParseFiles(ctx context.Context, paths []string) (map[string][]*entity.Entity, error) {
	if err := p.parser.ParseFiles(paths); err != nil {
		return nil, err
	}

	structs := p.parser.GetStructs()
	return p.adapter.ToMap(structs), nil
}
