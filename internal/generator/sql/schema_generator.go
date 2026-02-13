package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/ak/structify/internal/domain/entity"
	"github.com/ak/structify/internal/mapper"
	"github.com/lib/pq"
)

type SchemaGenerator struct {
	mapper *mapper.Mapper
}

func NewSchemaGenerator() *SchemaGenerator {
	return &SchemaGenerator{
		mapper: mapper.NewMapper(),
	}
}

func (g *SchemaGenerator) Generate(ctx context.Context, entities []*entity.Entity) (string, error) {
	var sb strings.Builder

	for _, ent := range entities {
		tableName := g.getTableName(ent)
		sb.WriteString(g.generateTable(ent, tableName))
	}

	return sb.String(), nil
}

func (g *SchemaGenerator) generateTable(ent *entity.Entity, tableName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))

	fields := ent.GetGenerateableFields()
	var columnDefs []string
	for _, field := range fields {
		colDef := g.generateColumn(field)
		if colDef != "" {
			columnDefs = append(columnDefs, "    "+colDef)
		}
	}

	for i, col := range columnDefs {
		sb.WriteString(col)
		if i < len(columnDefs)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}

	sb.WriteString(");\n\n")

	return sb.String()
}

func (g *SchemaGenerator) generateColumn(field entity.Field) string {
	if !field.ShouldGenerate() {
		return ""
	}

	mapping := g.mapper.MapType(field.Type)
	columnDef := g.mapper.FormatColumnDefinition(field.Name, mapping, g.getFieldTags(field))

	return pq.QuoteIdentifier(entity.ToSnakeCase(field.Name)) + " " + columnDef
}

func (g *SchemaGenerator) getFieldTags(field entity.Field) []string {
	var tags []string
	if field.IsPrimary {
		tags = append(tags, "pk")
	}
	if field.IsUnique {
		tags = append(tags, "unique")
	}
	if field.IsIgnored {
		tags = append(tags, "-")
	}
	return tags
}

func (g *SchemaGenerator) getTableName(ent *entity.Entity) string {
	if ent.TableName != "" {
		return pq.QuoteIdentifier(ent.TableName)
	}
	return pq.QuoteIdentifier(ent.GetTableName())
}
