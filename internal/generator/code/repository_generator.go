package code

import (
	"context"
	"fmt"
	"strings"

	"github.com/n0xum/structify/internal/domain/entity"
	"github.com/n0xum/structify/internal/mapper"
)

type RepositoryGenerator struct {
	mapper *mapper.Mapper
}

func NewRepositoryGenerator() *RepositoryGenerator {
	return &RepositoryGenerator{
		mapper: mapper.NewMapper(),
	}
}

func (g *RepositoryGenerator) Generate(ctx context.Context, packageName string, entities []*entity.Entity) (string, error) {
	var sb strings.Builder

	sb.WriteString("package ")
	if packageName == "" {
		sb.WriteString("main\n\n")
	} else {
		sb.WriteString(packageName + "\n\n")
	}

	sb.WriteString("import (\n")
	sb.WriteString("    \"database/sql\"\n")
	sb.WriteString("    \"context\"\n")
	sb.WriteString(")\n\n")

	for _, ent := range entities {
		g.generateRepository(&sb, ent)
	}

	return sb.String(), nil
}

func (g *RepositoryGenerator) generateRepository(sb *strings.Builder, ent *entity.Entity) {
	generateableFields := ent.GetGenerateableFields()

	g.generateEntityStruct(sb, ent, generateableFields)
	g.generateCreateMethod(sb, ent, generateableFields)
	g.generateGetByIDMethod(sb, ent, generateableFields)
	g.generateUpdateMethod(sb, ent, generateableFields)
	g.generateDeleteMethod(sb, ent, generateableFields)
	g.generateListMethod(sb, ent, generateableFields)
}

func (g *RepositoryGenerator) generateEntityStruct(sb *strings.Builder, ent *entity.Entity, fields []entity.Field) {
	sb.WriteString(fmt.Sprintf("type %s struct {\n", ent.Name))

	for _, field := range fields {
		sb.WriteString(fmt.Sprintf("    %s %s\n", field.Name, field.Type))
	}

	sb.WriteString("}\n\n")
}

func (g *RepositoryGenerator) generateCreateMethod(sb *strings.Builder, ent *entity.Entity, fields []entity.Field) {
	tableName := ent.GetTableName()

	sb.WriteString(fmt.Sprintf("func Create%s(ctx context.Context, db *sql.DB, item *%s) (*%s, error) {\n", ent.Name, ent.Name, ent.Name))

	var columns []string
	var args []string
	argCount := 1

	for _, field := range fields {
		if field.IsPrimary && field.Type == "int64" {
			continue
		}
		if !field.ShouldGenerate() {
			continue
		}
		colName := entity.ToSnakeCase(field.Name)
		columns = append(columns, colName)
		args = append(args, fmt.Sprintf("item.%s", field.Name))
		argCount++
	}

	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))

	sb.WriteString(fmt.Sprintf("    query := `%s`\n", query))
	sb.WriteString("    var id int64\n")
	sb.WriteString("    err := db.QueryRowContext(ctx, query, ")
	for i, arg := range args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg)
	}
	sb.WriteString(").Scan(&id)\n")
	sb.WriteString("    if err != nil {\n")
	sb.WriteString("        return nil, err\n")
	sb.WriteString("    }\n")
	sb.WriteString(fmt.Sprintf("    return Get%sByID(ctx, db, id)\n", ent.Name))
	sb.WriteString("}\n\n")
}

func (g *RepositoryGenerator) generateGetByIDMethod(sb *strings.Builder, ent *entity.Entity, fields []entity.Field) {
	tableName := ent.GetTableName()

	sb.WriteString(fmt.Sprintf("func Get%sByID(ctx context.Context, db *sql.DB, id int64) (*%s, error) {\n", ent.Name, ent.Name))

	var columns []string
	for _, field := range fields {
		if !field.ShouldGenerate() {
			continue
		}
		colName := entity.ToSnakeCase(field.Name)
		columns = append(columns, colName)
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE id = $1", strings.Join(columns, ", "), tableName)
	sb.WriteString(fmt.Sprintf("    query := `%s`\n", query))

	sb.WriteString(fmt.Sprintf("    var item %s\n", ent.Name))
	sb.WriteString("    err := db.QueryRowContext(ctx, query, id).Scan(")

	var scanFields []string
	for _, field := range fields {
		if !field.ShouldGenerate() {
			continue
		}
		scanFields = append(scanFields, fmt.Sprintf("&item.%s", field.Name))
	}
	sb.WriteString(strings.Join(scanFields, ", "))
	sb.WriteString(")\n")

	sb.WriteString("    if err != nil {\n")
	sb.WriteString("        return nil, err\n")
	sb.WriteString("    }\n")
	sb.WriteString("    return &item, nil\n")
	sb.WriteString("}\n\n")
}

func (g *RepositoryGenerator) generateUpdateMethod(sb *strings.Builder, ent *entity.Entity, fields []entity.Field) {
	tableName := ent.GetTableName()

	sb.WriteString(fmt.Sprintf("func Update%s(ctx context.Context, db *sql.DB, item *%s) error {\n", ent.Name, ent.Name))

	var updates []string
	var args []string

	for _, field := range fields {
		if !field.ShouldGenerate() {
			continue
		}
		if field.IsPrimary {
			continue
		}
		colName := entity.ToSnakeCase(field.Name)
		updates = append(updates, fmt.Sprintf("%s = $%d", colName, len(updates)+1))
		args = append(args, fmt.Sprintf("item.%s", field.Name))
	}

	args = append(args, "item.ID")

	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d",
		tableName,
		strings.Join(updates, ", "),
		len(updates)+1)

	sb.WriteString(fmt.Sprintf("    query := `%s`\n", query))
	sb.WriteString("    _, err := db.ExecContext(ctx, query, ")
	for i, arg := range args {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(arg)
	}
	sb.WriteString(")\n")
	sb.WriteString("    return err\n")
	sb.WriteString("}\n\n")
}

func (g *RepositoryGenerator) generateDeleteMethod(sb *strings.Builder, ent *entity.Entity, fields []entity.Field) {
	tableName := ent.GetTableName()

	sb.WriteString(fmt.Sprintf("func Delete%s(ctx context.Context, db *sql.DB, id int64) error {\n", ent.Name))

	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", tableName)
	sb.WriteString(fmt.Sprintf("    query := `%s`\n", query))

	sb.WriteString("    _, err := db.ExecContext(ctx, query, id)\n")
	sb.WriteString("    return err\n")
	sb.WriteString("}\n\n")
}

func (g *RepositoryGenerator) generateListMethod(sb *strings.Builder, ent *entity.Entity, fields []entity.Field) {
	tableName := ent.GetTableName()

	sb.WriteString(fmt.Sprintf("func List%s(ctx context.Context, db *sql.DB) ([]*%s, error) {\n", ent.Name, ent.Name))

	var columns []string
	for _, field := range fields {
		if !field.ShouldGenerate() {
			continue
		}
		colName := entity.ToSnakeCase(field.Name)
		columns = append(columns, colName)
	}

	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY id", strings.Join(columns, ", "), tableName)
	sb.WriteString(fmt.Sprintf("    query := `%s`\n", query))

	sb.WriteString("    rows, err := db.QueryContext(ctx, query)\n")
	sb.WriteString("    if err != nil {\n")
	sb.WriteString("        return nil, err\n")
	sb.WriteString("    }\n")
	sb.WriteString("    defer rows.Close()\n\n")

	sb.WriteString(fmt.Sprintf("    var items []*%s\n", ent.Name))

	sb.WriteString("    for rows.Next() {\n")
	sb.WriteString(fmt.Sprintf("        var item %s\n", ent.Name))

	var scanFields []string
	for _, field := range fields {
		if !field.ShouldGenerate() {
			continue
		}
		scanFields = append(scanFields, fmt.Sprintf("&item.%s", field.Name))
	}
	sb.WriteString(fmt.Sprintf("        if err := rows.Scan(%s); err != nil {\n", strings.Join(scanFields, ", ")))
	sb.WriteString("            return nil, err\n")
	sb.WriteString("        }\n")
	sb.WriteString("        items = append(items, &item)\n")
	sb.WriteString("    }\n")
	sb.WriteString("    return items, nil\n")
	sb.WriteString("}\n\n")
}
