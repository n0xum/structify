package adapter

import (
	"strings"

	"github.com/n0xum/structify/internal/domain/entity"
	"github.com/n0xum/structify/internal/parser"
)

type ParserAdapter struct{}

func NewParserAdapter() *ParserAdapter {
	return &ParserAdapter{}
}

func (a *ParserAdapter) ToDomain(pStruct *parser.Struct) *entity.Entity {
	if pStruct == nil {
		return nil
	}

	domainFields := make([]entity.Field, 0, len(pStruct.Fields))

	for _, pField := range pStruct.Fields {
		domainField := a.toDomainField(pField)
		domainFields = append(domainFields, domainField)
	}

	domainEntity := &entity.Entity{
		Name:      pStruct.Name,
		Fields:    domainFields,
		TableName: pStruct.TableName,
		Package:   pStruct.PackageName,
	}

	if customTable := a.extractCustomTableName(pStruct.Fields); customTable != "" {
		domainEntity.TableName = customTable
	}

	return domainEntity
}

func (a *ParserAdapter) ToDomainSlice(pStructs []*parser.Struct) []*entity.Entity {
	if pStructs == nil {
		return nil
	}

	domainEntities := make([]*entity.Entity, 0, len(pStructs))
	for _, pStruct := range pStructs {
		if ent := a.ToDomain(pStruct); ent != nil {
			domainEntities = append(domainEntities, ent)
		}
	}
	return domainEntities
}

func (a *ParserAdapter) toDomainField(pField parser.Field) entity.Field {
	tags := a.parseTags(pField.DatabaseTag)

	domainField := entity.Field{
		Name:      pField.Name,
		Type:      pField.Type,
		IsPrimary: a.hasTag(tags, "pk"),
		IsUnique:  a.hasTag(tags, "unique"),
		IsIgnored: a.hasTag(tags, "-"),
	}

	return domainField
}

func (a *ParserAdapter) parseTags(tag string) []string {
	if tag == "" {
		return nil
	}
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}
	parts := strings.Split(tag, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func (a *ParserAdapter) hasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (a *ParserAdapter) extractCustomTableName(fields []parser.Field) string {
	for _, field := range fields {
		tags := a.parseTags(field.DatabaseTag)
		for _, tag := range tags {
			if strings.HasPrefix(tag, "table:") {
				return strings.TrimPrefix(tag, "table:")
			}
		}
	}
	return ""
}

func (a *ParserAdapter) ToMap(structs map[string][]*parser.Struct) map[string][]*entity.Entity {
	if structs == nil {
		return nil
	}

	result := make(map[string][]*entity.Entity)
	for pkgName, pStructs := range structs {
		result[pkgName] = a.ToDomainSlice(pStructs)
	}
	return result
}
