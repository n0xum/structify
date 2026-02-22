package application

import (
	"context"

	"github.com/n0xum/structify/internal/adapter"
	"github.com/n0xum/structify/internal/domain/entity"
	"github.com/n0xum/structify/internal/parser"
)

type ParserWrapper struct {
	parser  *parser.Parser
	adapter *adapter.ParserAdapter
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

// ParseInterfaces parses Go files and returns discovered interfaces,
// bound to an entity for repository method classification.
func (p *ParserWrapper) ParseInterfaces(ctx context.Context, interfacePaths []string, ent *entity.Entity) ([]*entity.RepositoryInterface, error) {
	ifaceParser := parser.New()
	if err := ifaceParser.ParseFiles(interfacePaths); err != nil {
		return nil, err
	}

	interfaces := ifaceParser.GetInterfaces()
	var result []*entity.RepositoryInterface
	for _, ifaces := range interfaces {
		for _, iface := range ifaces {
			repo := p.adapter.ToRepositoryInterface(iface, ent)
			if repo != nil {
				result = append(result, repo)
			}
		}
	}
	return result, nil
}
