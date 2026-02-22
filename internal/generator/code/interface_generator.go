package code

import (
	"context"
	"fmt"
	"strings"

	"github.com/n0xum/structify/internal/domain/entity"
	"github.com/n0xum/structify/internal/util"
)

// GenerateFromInterface generates a repository implementation from
// an interface definition and an entity.
func (g *RepositoryGenerator) GenerateFromInterface(ctx context.Context, packageName string, ent *entity.Entity, repo *entity.RepositoryInterface) (string, error) {
	var sb strings.Builder

	if packageName == "" {
		packageName = "main"
	}

	sb.WriteString("package " + packageName + "\n\n")
	sb.WriteString("import (\n")
	sb.WriteString("\t\"context\"\n")
	sb.WriteString("\t\"database/sql\"\n")
	sb.WriteString(")\n\n")

	implName := repo.Name + "Impl"

	// Impl struct
	sb.WriteString(fmt.Sprintf("type %s struct {\n", implName))
	sb.WriteString("\tdb *sql.DB\n")
	sb.WriteString("}\n\n")

	// Constructor
	constructorName := "New" + repo.Name
	sb.WriteString(fmt.Sprintf("func %s(db *sql.DB) %s {\n", constructorName, repo.Name))
	sb.WriteString(fmt.Sprintf("\treturn &%s{db: db}\n", implName))
	sb.WriteString("}\n\n")

	// Generate each method
	for _, method := range repo.Methods {
		g.generateInterfaceMethod(&sb, implName, ent, method)
	}

	return sb.String(), nil
}

func (g *RepositoryGenerator) generateInterfaceMethod(sb *strings.Builder, implName string, ent *entity.Entity, method entity.RepositoryMethod) {
	switch method.Kind {
	case entity.MethodCreate:
		g.genCreate(sb, implName, ent, method)
	case entity.MethodGetByID:
		g.genGetByID(sb, implName, ent, method)
	case entity.MethodUpdate:
		g.genUpdate(sb, implName, ent, method)
	case entity.MethodDelete:
		g.genDelete(sb, implName, ent, method)
	case entity.MethodList:
		g.genList(sb, implName, ent, method)
	case entity.MethodFindBy:
		g.genFindBy(sb, implName, ent, method)
	case entity.MethodSmartQuery:
		g.genSmartQuery(sb, implName, ent, method)
	case entity.MethodCustomSQL:
		g.genCustomSQL(sb, implName, ent, method)
	}
}

func (g *RepositoryGenerator) genCreate(sb *strings.Builder, implName string, ent *entity.Entity, method entity.RepositoryMethod) {
	tableName := ent.GetQuotedTableName()
	fields := ent.GetGenerateableFields()

	// Method signature
	sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, item *%s) (*%s, error) {\n",
		implName, method.Name, ent.Name, ent.Name))

	var columns []string
	var args []string
	for _, field := range fields {
		if field.IsPrimary && field.Type == "int64" && !ent.HasCompositePrimaryKey() {
			continue
		}
		colName := util.ToSnakeCase(field.Name)
		columns = append(columns, colName)
		args = append(args, "item."+field.Name)
	}

	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	// Build RETURNING clause
	var returningCols []string
	for _, field := range fields {
		returningCols = append(returningCols, util.ToSnakeCase(field.Name))
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "),
		strings.Join(returningCols, ", "))

	sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", query))
	sb.WriteString(fmt.Sprintf("\tvar result %s\n", ent.Name))
	sb.WriteString("\terr := r.db.QueryRowContext(ctx, query, ")
	sb.WriteString(strings.Join(args, ", "))
	sb.WriteString(").Scan(")

	var scanFields []string
	for _, field := range fields {
		scanFields = append(scanFields, "&result."+field.Name)
	}
	sb.WriteString(strings.Join(scanFields, ", "))
	sb.WriteString(")\n")

	sb.WriteString("\tif err != nil {\n")
	sb.WriteString("\t\treturn nil, err\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\treturn &result, nil\n")
	sb.WriteString("}\n\n")
}

func (g *RepositoryGenerator) genGetByID(sb *strings.Builder, implName string, ent *entity.Entity, method entity.RepositoryMethod) {
	tableName := ent.GetQuotedTableName()
	fields := ent.GetGenerateableFields()

	// Build param list from method params
	paramStr := g.buildParamString(method)

	sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) (*%s, error) {\n",
		implName, method.Name, paramStr, ent.Name))

	var columns []string
	for _, field := range fields {
		columns = append(columns, util.ToSnakeCase(field.Name))
	}

	// WHERE clause from PK fields
	pkFields := ent.GetPrimaryKeyFields()
	var whereParts []string
	for i, pk := range pkFields {
		whereParts = append(whereParts, fmt.Sprintf("%s = $%d", util.ToSnakeCase(pk.Name), i+1))
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(columns, ", "), tableName, strings.Join(whereParts, " AND "))

	sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", query))
	sb.WriteString(fmt.Sprintf("\tvar item %s\n", ent.Name))
	sb.WriteString("\terr := r.db.QueryRowContext(ctx, query, ")

	var paramNames []string
	for _, p := range method.Params {
		paramNames = append(paramNames, p.Name)
	}
	sb.WriteString(strings.Join(paramNames, ", "))
	sb.WriteString(").Scan(")

	var scanFields []string
	for _, field := range fields {
		scanFields = append(scanFields, "&item."+field.Name)
	}
	sb.WriteString(strings.Join(scanFields, ", "))
	sb.WriteString(")\n")

	sb.WriteString("\tif err != nil {\n")
	sb.WriteString("\t\treturn nil, err\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\treturn &item, nil\n")
	sb.WriteString("}\n\n")
}

func (g *RepositoryGenerator) genUpdate(sb *strings.Builder, implName string, ent *entity.Entity, method entity.RepositoryMethod) {
	tableName := ent.GetQuotedTableName()
	fields := ent.GetGenerateableFields()

	sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, item *%s) error {\n",
		implName, method.Name, ent.Name))

	var updates []string
	var args []string
	for _, field := range fields {
		if field.IsPrimary {
			continue
		}
		colName := util.ToSnakeCase(field.Name)
		updates = append(updates, fmt.Sprintf("%s = $%d", colName, len(updates)+1))
		args = append(args, "item."+field.Name)
	}

	pkFields := ent.GetPrimaryKeyFields()
	var whereParts []string
	for i, pk := range pkFields {
		colName := util.ToSnakeCase(pk.Name)
		whereParts = append(whereParts, fmt.Sprintf("%s = $%d", colName, len(updates)+i+1))
		args = append(args, "item."+pk.Name)
	}

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		tableName, strings.Join(updates, ", "), strings.Join(whereParts, " AND "))

	sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", query))
	sb.WriteString("\t_, err := r.db.ExecContext(ctx, query, ")
	sb.WriteString(strings.Join(args, ", "))
	sb.WriteString(")\n")
	sb.WriteString("\treturn err\n")
	sb.WriteString("}\n\n")
}

func (g *RepositoryGenerator) genDelete(sb *strings.Builder, implName string, ent *entity.Entity, method entity.RepositoryMethod) {
	tableName := ent.GetQuotedTableName()

	paramStr := g.buildParamString(method)

	sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) error {\n",
		implName, method.Name, paramStr))

	pkFields := ent.GetPrimaryKeyFields()
	var whereParts []string
	for i, pk := range pkFields {
		whereParts = append(whereParts, fmt.Sprintf("%s = $%d", util.ToSnakeCase(pk.Name), i+1))
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, strings.Join(whereParts, " AND "))

	sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", query))
	sb.WriteString("\t_, err := r.db.ExecContext(ctx, query, ")

	var paramNames []string
	for _, p := range method.Params {
		paramNames = append(paramNames, p.Name)
	}
	sb.WriteString(strings.Join(paramNames, ", "))
	sb.WriteString(")\n")
	sb.WriteString("\treturn err\n")
	sb.WriteString("}\n\n")
}

func (g *RepositoryGenerator) genList(sb *strings.Builder, implName string, ent *entity.Entity, method entity.RepositoryMethod) {
	tableName := ent.GetQuotedTableName()
	fields := ent.GetGenerateableFields()

	sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context) ([]*%s, error) {\n",
		implName, method.Name, ent.Name))

	var columns []string
	for _, field := range fields {
		columns = append(columns, util.ToSnakeCase(field.Name))
	}

	pkFields := ent.GetPrimaryKeyFields()
	var orderBy []string
	for _, pk := range pkFields {
		orderBy = append(orderBy, util.ToSnakeCase(pk.Name))
	}
	orderClause := "id"
	if len(orderBy) > 0 {
		orderClause = strings.Join(orderBy, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s",
		strings.Join(columns, ", "), tableName, orderClause)

	sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", query))
	g.writeRowsLoop(sb, ent, fields)
}

func (g *RepositoryGenerator) genFindBy(sb *strings.Builder, implName string, ent *entity.Entity, method entity.RepositoryMethod) {
	tableName := ent.GetQuotedTableName()
	fields := ent.GetGenerateableFields()

	paramStr := g.buildParamString(method)

	if method.ReturnsSingle {
		sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) (*%s, error) {\n",
			implName, method.Name, paramStr, ent.Name))
	} else {
		sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) ([]*%s, error) {\n",
			implName, method.Name, paramStr, ent.Name))
	}

	var columns []string
	for _, field := range fields {
		columns = append(columns, util.ToSnakeCase(field.Name))
	}

	// WHERE clause from FindBy fields
	var whereParts []string
	for i, fieldName := range method.FindByFields {
		colName := util.ToSnakeCase(fieldName)
		whereParts = append(whereParts, fmt.Sprintf("%s = $%d", colName, i+1))
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s",
		strings.Join(columns, ", "), tableName, strings.Join(whereParts, " AND "))

	sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", query))

	var paramNames []string
	for _, p := range method.Params {
		paramNames = append(paramNames, p.Name)
	}

	if method.ReturnsSingle {
		sb.WriteString(fmt.Sprintf("\tvar item %s\n", ent.Name))
		sb.WriteString("\terr := r.db.QueryRowContext(ctx, query, ")
		sb.WriteString(strings.Join(paramNames, ", "))
		sb.WriteString(").Scan(")

		var scanFields []string
		for _, field := range fields {
			scanFields = append(scanFields, "&item."+field.Name)
		}
		sb.WriteString(strings.Join(scanFields, ", "))
		sb.WriteString(")\n")
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString("\t\treturn nil, err\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\treturn &item, nil\n")
		sb.WriteString("}\n\n")
	} else {
		g.writeRowsLoopWithParams(sb, ent, fields, paramNames)
	}
}

func (g *RepositoryGenerator) genCustomSQL(sb *strings.Builder, implName string, ent *entity.Entity, method entity.RepositoryMethod) {
	paramStr := g.buildParamString(method)

	var paramNames []string
	for _, p := range method.Params {
		paramNames = append(paramNames, p.Name)
	}

	// Method returns a scalar non-entity type (e.g. float64 for SUM queries)
	if method.ScalarReturnType != "" {
		sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) (%s, error) {\n",
			implName, method.Name, paramStr, method.ScalarReturnType))
		sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", method.CustomSQL))
		sb.WriteString(fmt.Sprintf("\tvar result %s\n", method.ScalarReturnType))
		sb.WriteString("\terr := r.db.QueryRowContext(ctx, query")
		if len(paramNames) > 0 {
			sb.WriteString(", ")
			sb.WriteString(strings.Join(paramNames, ", "))
		}
		sb.WriteString(").Scan(&result)\n")
		sb.WriteString("\treturn result, err\n")
		sb.WriteString("}\n\n")
		return
	}

	// Method returns only error (e.g. UPDATE/DELETE with no result set)
	if !method.HasEntityReturn {
		sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) error {\n",
			implName, method.Name, paramStr))
		sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", method.CustomSQL))
		sb.WriteString("\t_, err := r.db.ExecContext(ctx, query")
		if len(paramNames) > 0 {
			sb.WriteString(", ")
			sb.WriteString(strings.Join(paramNames, ", "))
		}
		sb.WriteString(")\n")
		sb.WriteString("\treturn err\n")
		sb.WriteString("}\n\n")
		return
	}

	if method.ReturnsSingle {
		sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) (*%s, error) {\n",
			implName, method.Name, paramStr, ent.Name))
	} else {
		sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) ([]*%s, error) {\n",
			implName, method.Name, paramStr, ent.Name))
	}

	sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", method.CustomSQL))

	fields := ent.GetGenerateableFields()

	if method.ReturnsSingle {
		sb.WriteString(fmt.Sprintf("\tvar item %s\n", ent.Name))
		sb.WriteString("\terr := r.db.QueryRowContext(ctx, query, ")
		sb.WriteString(strings.Join(paramNames, ", "))
		sb.WriteString(").Scan(")

		var scanFields []string
		for _, field := range fields {
			scanFields = append(scanFields, "&item."+field.Name)
		}
		sb.WriteString(strings.Join(scanFields, ", "))
		sb.WriteString(")\n")
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString("\t\treturn nil, err\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\treturn &item, nil\n")
		sb.WriteString("}\n\n")
	} else {
		g.writeRowsLoopWithParams(sb, ent, fields, paramNames)
	}
}

func (g *RepositoryGenerator) genSmartQuery(sb *strings.Builder, implName string, ent *entity.Entity, method entity.RepositoryMethod) {
	paramStr := g.buildParamString(method)

	// Determine return type based on GeneratedSQL or pattern
	returnsCount := strings.Contains(method.GeneratedSQL, "COUNT(*)")
	returnsExists := strings.Contains(method.GeneratedSQL, "EXISTS(")

	switch {
	case returnsCount:
		sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) (int64, error) {\n",
			implName, method.Name, paramStr))
	case returnsExists:
		sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) (bool, error) {\n",
			implName, method.Name, paramStr))
	case method.ReturnsSingle:
		sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) (*%s, error) {\n",
			implName, method.Name, paramStr, ent.Name))
	default:
		sb.WriteString(fmt.Sprintf("func (r *%s) %s(ctx context.Context, %s) ([]*%s, error) {\n",
			implName, method.Name, paramStr, ent.Name))
	}

	sb.WriteString(fmt.Sprintf("\tquery := `%s`\n", method.GeneratedSQL))

	var paramNames []string
	for _, p := range method.Params {
		paramNames = append(paramNames, p.Name)
	}

	if returnsCount {
		sb.WriteString("\tvar count int64\n")
		sb.WriteString("\terr := r.db.QueryRowContext(ctx, query")
		if len(paramNames) > 0 {
			sb.WriteString(", ")
			sb.WriteString(strings.Join(paramNames, ", "))
		}
		sb.WriteString(").Scan(&count)\n")
		sb.WriteString("\treturn count, err\n")
		sb.WriteString("}\n\n")
		return
	}

	if returnsExists {
		sb.WriteString("\tvar exists bool\n")
		sb.WriteString("\terr := r.db.QueryRowContext(ctx, query")
		if len(paramNames) > 0 {
			sb.WriteString(", ")
			sb.WriteString(strings.Join(paramNames, ", "))
		}
		sb.WriteString(").Scan(&exists)\n")
		sb.WriteString("\treturn exists, err\n")
		sb.WriteString("}\n\n")
		return
	}

	fields := ent.GetGenerateableFields()

	if method.ReturnsSingle {
		sb.WriteString(fmt.Sprintf("\tvar item %s\n", ent.Name))
		sb.WriteString("\terr := r.db.QueryRowContext(ctx, query")
		if len(paramNames) > 0 {
			sb.WriteString(", ")
			sb.WriteString(strings.Join(paramNames, ", "))
		}
		sb.WriteString(").Scan(")

		var scanFields []string
		for _, field := range fields {
			scanFields = append(scanFields, "&item."+field.Name)
		}
		sb.WriteString(strings.Join(scanFields, ", "))
		sb.WriteString(")\n")
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString("\t\treturn nil, err\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\treturn &item, nil\n")
		sb.WriteString("}\n\n")
	} else {
		g.writeRowsLoopWithParams(sb, ent, fields, paramNames)
	}
}

func (g *RepositoryGenerator) buildParamString(method entity.RepositoryMethod) string {
	var parts []string
	for _, p := range method.Params {
		parts = append(parts, p.Name+" "+p.Type)
	}
	return strings.Join(parts, ", ")
}

func (g *RepositoryGenerator) writeRowsLoop(sb *strings.Builder, ent *entity.Entity, fields []entity.Field) {
	sb.WriteString("\trows, err := r.db.QueryContext(ctx, query)\n")
	sb.WriteString("\tif err != nil {\n")
	sb.WriteString("\t\treturn nil, err\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tdefer rows.Close()\n\n")

	sb.WriteString(fmt.Sprintf("\tvar items []*%s\n", ent.Name))
	sb.WriteString("\tfor rows.Next() {\n")
	sb.WriteString(fmt.Sprintf("\t\tvar item %s\n", ent.Name))

	var scanFields []string
	for _, field := range fields {
		scanFields = append(scanFields, "&item."+field.Name)
	}

	sb.WriteString(fmt.Sprintf("\t\tif err := rows.Scan(%s); err != nil {\n", strings.Join(scanFields, ", ")))
	sb.WriteString("\t\t\treturn nil, err\n")
	sb.WriteString("\t\t}\n")
	sb.WriteString("\t\titems = append(items, &item)\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tif err := rows.Err(); err != nil {\n")
	sb.WriteString("\t\treturn nil, err\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\treturn items, nil\n")
	sb.WriteString("}\n\n")
}

func (g *RepositoryGenerator) writeRowsLoopWithParams(sb *strings.Builder, ent *entity.Entity, fields []entity.Field, paramNames []string) {
	sb.WriteString("\trows, err := r.db.QueryContext(ctx, query, ")
	sb.WriteString(strings.Join(paramNames, ", "))
	sb.WriteString(")\n")
	sb.WriteString("\tif err != nil {\n")
	sb.WriteString("\t\treturn nil, err\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tdefer rows.Close()\n\n")

	sb.WriteString(fmt.Sprintf("\tvar items []*%s\n", ent.Name))
	sb.WriteString("\tfor rows.Next() {\n")
	sb.WriteString(fmt.Sprintf("\t\tvar item %s\n", ent.Name))

	var scanFields []string
	for _, field := range fields {
		scanFields = append(scanFields, "&item."+field.Name)
	}

	sb.WriteString(fmt.Sprintf("\t\tif err := rows.Scan(%s); err != nil {\n", strings.Join(scanFields, ", ")))
	sb.WriteString("\t\t\treturn nil, err\n")
	sb.WriteString("\t\t}\n")
	sb.WriteString("\t\titems = append(items, &item)\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\tif err := rows.Err(); err != nil {\n")
	sb.WriteString("\t\treturn nil, err\n")
	sb.WriteString("\t}\n")
	sb.WriteString("\treturn items, nil\n")
	sb.WriteString("}\n\n")
}
