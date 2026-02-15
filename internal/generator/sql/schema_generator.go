package sql

import (
	"context"
	"fmt"
	"strings"

	"github.com/n0xum/structify/internal/domain/entity"
	"github.com/n0xum/structify/internal/mapper"
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
		sb.WriteString(g.generateIndexes(ent, tableName))
	}

	return sb.String(), nil
}

func (g *SchemaGenerator) generateTable(ent *entity.Entity, tableName string) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("CREATE TABLE %s (\n", tableName))

	fields := ent.GetGenerateableFields()
	var columnDefs []string
	for _, field := range fields {
		colDef := g.generateColumn(ent, field)
		if colDef != "" {
			columnDefs = append(columnDefs, colDef)
		}
	}

	// Check if we have composite FKs
	hasCompositeFK := g.hasCompositeForeignKey(fields)

	// Write column definitions
	for i, col := range columnDefs {
		sb.WriteString("    ")
		sb.WriteString(col)
		if i < len(columnDefs)-1 || ent.HasCompositePrimaryKey() || ent.HasCompositeUniqueConstraints() || hasCompositeFK {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}

	// Add composite PRIMARY KEY if needed
	pkFields := ent.GetPrimaryKeyFields()
	if len(pkFields) > 1 {
		var pkColumns []string
		for _, field := range pkFields {
			if field.ShouldGenerate() {
				pkColumns = append(pkColumns, pq.QuoteIdentifier(entity.ToSnakeCase(field.Name)))
			}
		}
		if len(pkColumns) > 1 {
			sb.WriteString(fmt.Sprintf("    PRIMARY KEY (%s)", strings.Join(pkColumns, ", ")))
			sb.WriteString("\n")
		}
	}

	// Add composite UNIQUE constraints
	uniqueConstraints := ent.GetUniqueConstraints()
	var constraintNames []string
	for constraintName, fields := range uniqueConstraints {
		if len(fields) > 1 {
			constraintNames = append(constraintNames, constraintName)
		}
	}

	for i, constraintName := range constraintNames {
		fields := uniqueConstraints[constraintName]
		var uniqueColumns []string
		for _, field := range fields {
			if field.ShouldGenerate() {
				uniqueColumns = append(uniqueColumns, pq.QuoteIdentifier(entity.ToSnakeCase(field.Name)))
			}
		}
		if len(uniqueColumns) > 1 {
			sb.WriteString("    UNIQUE")
			if strings.HasPrefix(constraintName, "uq_") {
				sb.WriteString(fmt.Sprintf(" %s", pq.QuoteIdentifier(constraintName)))
			}
			sb.WriteString(fmt.Sprintf(" (%s)", strings.Join(uniqueColumns, ", ")))
			if i < len(constraintNames)-1 || ent.HasCompositePrimaryKey() || hasCompositeFK {
				sb.WriteString(",\n")
			} else {
				sb.WriteString("\n")
			}
		}
	}

	// Add composite FOREIGN KEY constraints
	fkGroups := g.groupFieldsByFK(fields)
	var fkGroupNames []string
	for groupName := range fkGroups {
		if len(fkGroups[groupName]) > 1 {
			fkGroupNames = append(fkGroupNames, groupName)
		}
	}

	for i, fkGroupName := range fkGroupNames {
		fkFields := fkGroups[fkGroupName]
		if len(fkFields) < 2 {
			continue
		}

		var localColumns []string
		var refColumns []string
		var refTable string
		var onDelete, onUpdate string

		for _, field := range fkFields {
			if !field.ShouldGenerate() || field.FKReference == nil {
				continue
			}
			localColumns = append(localColumns, pq.QuoteIdentifier(entity.ToSnakeCase(field.Name)))
			refColumns = append(refColumns, pq.QuoteIdentifier(field.FKReference.Column))
			if refTable == "" {
				refTable = field.FKReference.Table
			}
			if field.FKOnDelete != "" {
				onDelete = field.FKOnDelete
			}
			if field.FKOnUpdate != "" {
				onUpdate = field.FKOnUpdate
			}
		}

		if len(localColumns) >= 2 {
			sb.WriteString("    FOREIGN KEY")
			if fkGroupName != "" {
				sb.WriteString(fmt.Sprintf(" %s", pq.QuoteIdentifier(fkGroupName)))
			}
			sb.WriteString(fmt.Sprintf(" (%s) REFERENCES %s (%s)",
				strings.Join(localColumns, ", "),
				pq.QuoteIdentifier(refTable),
				strings.Join(refColumns, ", ")))

			if onDelete != "" {
				sb.WriteString(fmt.Sprintf(" ON DELETE %s", g.formatCascadeAction(onDelete)))
			}
			if onUpdate != "" {
				sb.WriteString(fmt.Sprintf(" ON UPDATE %s", g.formatCascadeAction(onUpdate)))
			}

			if i < len(fkGroupNames)-1 {
				sb.WriteString(",\n")
			} else {
				sb.WriteString("\n")
			}
		}
	}

	sb.WriteString(");\n\n")

	return sb.String()
}

func (g *SchemaGenerator) generateColumn(ent *entity.Entity, field entity.Field) string {
	if !field.ShouldGenerate() {
		return ""
	}

	mapping := g.mapper.MapType(field.Type)
	// Don't add PRIMARY KEY here if it's part of a composite PK
	tags := g.getFieldTags(field)
	if field.IsPrimary && ent.HasCompositePrimaryKey() {
		// Check if this is a composite PK by removing PK tag temporarily
		tempTags := make([]string, 0, len(tags))
		for _, t := range tags {
			if t != "pk" {
				tempTags = append(tempTags, t)
			}
		}
		tags = tempTags
	}
	// Don't add UNIQUE here if it's part of a composite unique constraint
	if field.IsUnique && field.IndexGroup != "" {
		tempTags := make([]string, 0, len(tags))
		for _, t := range tags {
			if t != "unique" {
				tempTags = append(tempTags, t)
			}
		}
		tags = tempTags
	}
	columnDef := g.mapper.FormatColumnDefinition(field.Name, mapping, tags)

	// Add CHECK constraint from field.CheckExpr
	if field.CheckExpr != "" {
		columnDef += fmt.Sprintf(" CHECK (%s)", field.CheckExpr)
	}

	// Add enum CHECK constraint from field.EnumValues
	if len(field.EnumValues) > 0 {
		enumList := make([]string, len(field.EnumValues))
		for i, val := range field.EnumValues {
			enumList[i] = fmt.Sprintf("'%s'", val)
		}
		columnName := pq.QuoteIdentifier(entity.ToSnakeCase(field.Name))
		columnDef += fmt.Sprintf(" CHECK (%s IN (%s))", columnName, strings.Join(enumList, ", "))
	}

	// Add DEFAULT value from field.DefaultVal
	if field.DefaultVal != "" {
		columnDef += fmt.Sprintf(" DEFAULT %s", field.DefaultVal)
	}

	// Add FOREIGN KEY constraint from field.FKReference (only for single-column FKs)
	if field.FKReference != nil && field.FKGroup == "" {
		columnDef += fmt.Sprintf(" REFERENCES %s(%s)",
			pq.QuoteIdentifier(field.FKReference.Table),
			pq.QuoteIdentifier(field.FKReference.Column))

		// Add ON DELETE clause
		if field.FKOnDelete != "" {
			columnDef += fmt.Sprintf(" ON DELETE %s", g.formatCascadeAction(field.FKOnDelete))
		}

		// Add ON UPDATE clause
		if field.FKOnUpdate != "" {
			columnDef += fmt.Sprintf(" ON UPDATE %s", g.formatCascadeAction(field.FKOnUpdate))
		}
	}

	return pq.QuoteIdentifier(entity.ToSnakeCase(field.Name)) + " " + columnDef
}

// formatCascadeAction formats a cascade action (converts underscores to spaces)
func (g *SchemaGenerator) formatCascadeAction(action string) string {
	// Convert SET_NULL to SET NULL, NO_ACTION to NO ACTION
	return strings.ReplaceAll(action, "_", " ")
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

// generateIndexes creates CREATE INDEX statements for fields with index tags
func (g *SchemaGenerator) generateIndexes(ent *entity.Entity, tableName string) string {
	var sb strings.Builder

	// Group fields by index name for composite indexes
	indexGroups := g.groupFieldsByIndex(ent.Fields)

	for indexName, fields := range indexGroups {
		if len(fields) == 0 {
			continue
		}

		// Build column list
		columns := make([]string, len(fields))
		for i, field := range fields {
			columns[i] = pq.QuoteIdentifier(entity.ToSnakeCase(field.Name))
		}

		// Check if this is a unique index
		isUnique := fields[0].IsIndexUnique

		if isUnique {
			sb.WriteString(fmt.Sprintf("CREATE UNIQUE INDEX %s ON %s (%s);\n",
				pq.QuoteIdentifier(indexName),
				tableName,
				strings.Join(columns, ", ")))
		} else {
			sb.WriteString(fmt.Sprintf("CREATE INDEX %s ON %s (%s);\n",
				pq.QuoteIdentifier(indexName),
				tableName,
				strings.Join(columns, ", ")))
		}
	}

	if sb.Len() > 0 {
		sb.WriteString("\n")
	}

	return sb.String()
}

// groupFieldsByIndex groups fields by their index name
func (g *SchemaGenerator) groupFieldsByIndex(fields []entity.Field) map[string][]entity.Field {
	groups := make(map[string][]entity.Field)

	for _, field := range fields {
		if field.IndexName != "" {
			// Use the index name as the group key
			groups[field.IndexName] = append(groups[field.IndexName], field)
		}
	}

	return groups
}

// groupFieldsByFK groups fields by their FK group name for composite FKs
func (g *SchemaGenerator) groupFieldsByFK(fields []entity.Field) map[string][]entity.Field {
	groups := make(map[string][]entity.Field)

	for _, field := range fields {
		if field.FKGroup != "" && field.FKReference != nil {
			groups[field.FKGroup] = append(groups[field.FKGroup], field)
		}
	}

	return groups
}

// hasCompositeForeignKey checks if the entity has composite foreign keys
func (g *SchemaGenerator) hasCompositeForeignKey(fields []entity.Field) bool {
	fkGroups := g.groupFieldsByFK(fields)
	for _, fields := range fkGroups {
		if len(fields) > 1 {
			return true
		}
	}
	return false
}
