package code

import (
	"context"
	"fmt"
	"strings"

	"github.com/n0xum/structify/internal/domain/entity"
	"github.com/n0xum/structify/internal/mapper"
	"github.com/n0xum/structify/internal/util"
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
		g.generateRepository(&sb, ent, entities)
	}

	return sb.String(), nil
}

func (g *RepositoryGenerator) generateRepository(sb *strings.Builder, ent *entity.Entity, allEntities []*entity.Entity) {
	generateableFields := ent.GetGenerateableFields()

	g.generateEntityStruct(sb, ent, generateableFields)
	g.generateCreateMethod(sb, ent, generateableFields)

	// Generate different CRUD methods based on PK type
	if ent.HasCompositePrimaryKey() {
		g.generateGetByCompositePKMethod(sb, ent, generateableFields)
		g.generateUpdateMethod(sb, ent, generateableFields)
		g.generateDeleteByCompositePKMethod(sb, ent, generateableFields)
	} else {
		g.generateGetByIDMethod(sb, ent, generateableFields)
		g.generateUpdateMethod(sb, ent, generateableFields)
		g.generateDeleteMethod(sb, ent, generateableFields)
	}

	g.generateListMethod(sb, ent, generateableFields)

	// Generate JOIN methods for FK relationships
	g.generateJoinMethods(sb, ent, allEntities)
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

	for _, field := range fields {
		// Skip auto-increment int64 PK fields
		if field.IsPrimary && field.Type == "int64" && !ent.HasCompositePrimaryKey() {
			continue
		}
		if !field.ShouldGenerate() {
			continue
		}
		colName := util.ToSnakeCase(field.Name)
		columns = append(columns, colName)
		args = append(args, fmt.Sprintf("item.%s", field.Name))
	}

	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}

	var query string
	pkFields := ent.GetPrimaryKeyFields()
	if len(pkFields) == 1 && pkFields[0].Type == "int64" {
		// Single int64 PK - use RETURNING id
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id",
			tableName,
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "))
	} else {
		// Composite PK or non-int64 PK - use RETURNING all PK columns
		var returningColumns []string
		for _, pkField := range pkFields {
			returningColumns = append(returningColumns, util.ToSnakeCase(pkField.Name))
		}
		query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING %s",
			tableName,
			strings.Join(columns, ", "),
			strings.Join(placeholders, ", "),
			strings.Join(returningColumns, ", "))
	}

	sb.WriteString(fmt.Sprintf("    query := `%s`\n", query))

	if len(pkFields) == 1 && pkFields[0].Type == "int64" {
		// Single int64 PK
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
	} else {
		// Composite PK - scan all PK values
		for _, pkField := range pkFields {
			sb.WriteString(fmt.Sprintf("    var %s %s\n", util.ToSnakeCase(pkField.Name), pkField.Type))
		}
		sb.WriteString("    err := db.QueryRowContext(ctx, query, ")
		for i, arg := range args {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(arg)
		}
		sb.WriteString(").Scan(")
		for i, pkField := range pkFields {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("&%s", util.ToSnakeCase(pkField.Name)))
		}
		sb.WriteString(")\n")
		sb.WriteString("    if err != nil {\n")
		sb.WriteString("        return nil, err\n")
		sb.WriteString("    }\n")
		// Call Get method with all PK values
		sb.WriteString(fmt.Sprintf("    return Get%s(ctx, db, ", ent.Name))
		for i, pkField := range pkFields {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(util.ToSnakeCase(pkField.Name))
		}
		sb.WriteString(")\n")
	}
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
		colName := util.ToSnakeCase(field.Name)
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
		colName := util.ToSnakeCase(field.Name)
		updates = append(updates, fmt.Sprintf("%s = $%d", colName, len(updates)+1))
		args = append(args, fmt.Sprintf("item.%s", field.Name))
	}

	pkFields := ent.GetPrimaryKeyFields()
	var whereClause []string
	var whereArgs []string
	for i, pkField := range pkFields {
		colName := util.ToSnakeCase(pkField.Name)
		whereClause = append(whereClause, fmt.Sprintf("%s = $%d", colName, len(updates)+i+1))
		whereArgs = append(whereArgs, fmt.Sprintf("item.%s", pkField.Name))
	}

	args = append(args, whereArgs...)

	query := fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		tableName,
		strings.Join(updates, ", "),
		strings.Join(whereClause, " AND "))

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
		colName := util.ToSnakeCase(field.Name)
		columns = append(columns, colName)
	}

	// Build ORDER BY clause from primary key fields
	var orderBy []string
	pkFields := ent.GetPrimaryKeyFields()
	for _, pkField := range pkFields {
		orderBy = append(orderBy, util.ToSnakeCase(pkField.Name))
	}
	orderClause := "id"
	if len(orderBy) > 0 {
		orderClause = strings.Join(orderBy, ", ")
	}

	query := fmt.Sprintf("SELECT %s FROM %s ORDER BY %s", strings.Join(columns, ", "), tableName, orderClause)
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

// generateJoinMethods generates JOIN query methods for entities with FK relationships
func (g *RepositoryGenerator) generateJoinMethods(sb *strings.Builder, ent *entity.Entity, allEntities []*entity.Entity) {
	// Find all FK fields in the entity
	var fkFields []entity.Field
	for _, field := range ent.Fields {
		if field.FKReference != nil && field.ShouldGenerate() {
			fkFields = append(fkFields, field)
		}
	}

	if len(fkFields) == 0 {
		return
	}

	// Group FK fields by referenced table
	fkGroups := g.groupFKsByTable(fkFields)

	// Generate individual JOIN methods for each relationship
	for _, fk := range fkFields {
		relatedEntity := g.findEntity(fk.FKReference.Table, allEntities)
		if relatedEntity == nil {
			continue
		}
		g.generateSingleJoinMethod(sb, ent, relatedEntity, fk)
	}

	// Generate combined JOIN method if there are multiple FKs
	if len(fkGroups) > 1 {
		g.generateMultiJoinMethod(sb, ent, allEntities, fkGroups)
	}
}

// groupFKsByTable groups FK fields by their referenced table
func (g *RepositoryGenerator) groupFKsByTable(fkFields []entity.Field) map[string][]entity.Field {
	groups := make(map[string][]entity.Field)
	for _, fk := range fkFields {
		if fk.FKReference != nil {
			table := fk.FKReference.Table
			groups[table] = append(groups[table], fk)
		}
	}
	return groups
}

// findEntity finds an entity by table name
func (g *RepositoryGenerator) findEntity(tableName string, entities []*entity.Entity) *entity.Entity {
	for _, ent := range entities {
		if ent.GetTableName() == tableName {
			return ent
		}
	}
	return nil
}

// generateSingleJoinMethod generates a JOIN method for a single FK relationship
func (g *RepositoryGenerator) generateSingleJoinMethod(sb *strings.Builder, ent, relatedEnt *entity.Entity, fkField entity.Field) {
	// Generate result struct name
	resultStructName := fmt.Sprintf("%sWith%s", ent.Name, relatedEnt.Name)

	g.generateJoinResultStruct(sb, resultStructName, ent, relatedEnt)

	// Generate the JOIN method
	sb.WriteString(fmt.Sprintf("func Get%sWith%s(ctx context.Context, db *sql.DB, %sID int64) (*%s, error) {\n",
		ent.Name, relatedEnt.Name, util.ToSnakeCase(ent.Name), resultStructName))

	// Build JOIN query
	tableName := ent.GetTableName()
	relatedTableName := relatedEnt.GetTableName()

	var columns []string
	entColumns := g.getEntityColumns(ent)
	relatedColumns := g.getEntityColumns(relatedEnt)

	// Prefix columns with table aliases to avoid conflicts
	for _, col := range entColumns {
		columns = append(columns, fmt.Sprintf("%s.%s", tableName, col))
	}
	for _, col := range relatedColumns {
		columns = append(columns, fmt.Sprintf("%s.%s", relatedTableName, col))
	}

	query := fmt.Sprintf("SELECT %s FROM %s JOIN %s ON %s.%s = %s.%s WHERE %s.id = $1",
		strings.Join(columns, ", "),
		tableName,
		relatedTableName,
		tableName, util.ToSnakeCase(fkField.Name),
		relatedTableName, fkField.FKReference.Column,
		tableName)

	sb.WriteString(fmt.Sprintf("    query := `%s`\n", query))

	// Scan into result struct
	var scanFields []string
	entFields := ent.GetGenerateableFields()
	relatedFields := relatedEnt.GetGenerateableFields()

	for _, field := range entFields {
		scanFields = append(scanFields, fmt.Sprintf("&result.%s", field.Name))
	}
	for _, field := range relatedFields {
		scanFields = append(scanFields, fmt.Sprintf("&result.%s", field.Name))
	}

	sb.WriteString(fmt.Sprintf("    var result %s\n", resultStructName))
	sb.WriteString("    err := db.QueryRowContext(ctx, query, ")
	sb.WriteString(fmt.Sprintf("%sID", util.ToSnakeCase(ent.Name)))
	sb.WriteString(").Scan(")
	sb.WriteString(strings.Join(scanFields, ", "))
	sb.WriteString(")\n")

	sb.WriteString("    if err != nil {\n")
	sb.WriteString("        return nil, err\n")
	sb.WriteString("    }\n")
	sb.WriteString("    return &result, nil\n")
	sb.WriteString("}\n\n")
}

// generateMultiJoinMethod generates a JOIN method for all FK relationships
func (g *RepositoryGenerator) generateMultiJoinMethod(sb *strings.Builder, ent *entity.Entity, allEntities []*entity.Entity, fkGroups map[string][]entity.Field) {
	// Build result struct name
	resultStructName := fmt.Sprintf("%sWithRelations", ent.Name)

	// Collect all related entities
	var relatedEntities []*entity.Entity
	tablesJoined := make(map[string]bool)
	for _, fk := range ent.Fields {
		if fk.FKReference != nil && fk.ShouldGenerate() {
			table := fk.FKReference.Table
			if !tablesJoined[table] {
				if relatedEnt := g.findEntity(table, allEntities); relatedEnt != nil {
					relatedEntities = append(relatedEntities, relatedEnt)
					tablesJoined[table] = true
				}
			}
		}
	}

	g.generateMultiJoinResultStruct(sb, resultStructName, ent, relatedEntities)

	// Generate the JOIN method
	sb.WriteString(fmt.Sprintf("func Get%sWithRelations(ctx context.Context, db *sql.DB, id int64) (*%s, error) {\n",
		ent.Name, resultStructName))

	// Build JOIN query with multiple tables
	tableName := ent.GetTableName()

	var columns []string
	var joins []string
	entColumns := g.getEntityColumns(ent)

	// Add main entity columns
	for _, col := range entColumns {
		columns = append(columns, fmt.Sprintf("%s.%s", tableName, col))
	}

	// Add JOINs and related entity columns
	for _, relatedEnt := range relatedEntities {
		relatedTableName := relatedEnt.GetTableName()
		relatedColumns := g.getEntityColumns(relatedEnt)

		for _, col := range relatedColumns {
			columns = append(columns, fmt.Sprintf("%s.%s", relatedTableName, col))
		}
	}

	// Build JOIN clauses
	for _, relatedEnt := range relatedEntities {
		relatedTableName := relatedEnt.GetTableName()
		// Find FK field that references this table
		for _, fk := range ent.Fields {
			if fk.FKReference != nil && fk.FKReference.Table == relatedTableName && fk.ShouldGenerate() {
				joins = append(joins, fmt.Sprintf("JOIN %s ON %s.%s = %s.%s",
					relatedTableName,
					tableName, util.ToSnakeCase(fk.Name),
					relatedTableName, fk.FKReference.Column))
				break
			}
		}
	}

	joinClause := strings.Join(joins, " ")
	query := fmt.Sprintf("SELECT %s FROM %s %s WHERE %s.id = $1",
		strings.Join(columns, ", "),
		tableName,
		joinClause,
		tableName)

	sb.WriteString(fmt.Sprintf("    query := `%s`\n", query))

	// Scan into result struct
	var scanFields []string
	entFields := ent.GetGenerateableFields()
	for _, field := range entFields {
		scanFields = append(scanFields, fmt.Sprintf("&result.%s", field.Name))
	}
	for _, relatedEnt := range relatedEntities {
		relatedFields := relatedEnt.GetGenerateableFields()
		for _, field := range relatedFields {
			scanFields = append(scanFields, fmt.Sprintf("&result.%s", field.Name))
		}
	}

	sb.WriteString(fmt.Sprintf("    var result %s\n", resultStructName))
	sb.WriteString("    err := db.QueryRowContext(ctx, query, id).Scan(")
	sb.WriteString(strings.Join(scanFields, ", "))
	sb.WriteString(")\n")

	sb.WriteString("    if err != nil {\n")
	sb.WriteString("        return nil, err\n")
	sb.WriteString("    }\n")
	sb.WriteString("    return &result, nil\n")
	sb.WriteString("}\n\n")
}

// generateJoinResultStruct generates a struct for JOIN results
func (g *RepositoryGenerator) generateJoinResultStruct(sb *strings.Builder, structName string, ent1, ent2 *entity.Entity) {
	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	sb.WriteString(fmt.Sprintf("    %s\n", ent1.Name))
	sb.WriteString(fmt.Sprintf("    %s\n", ent2.Name))
	sb.WriteString("}\n\n")
}

// generateMultiJoinResultStruct generates a struct for multi-table JOIN results
func (g *RepositoryGenerator) generateMultiJoinResultStruct(sb *strings.Builder, structName string, mainEnt *entity.Entity, relatedEnts []*entity.Entity) {
	sb.WriteString(fmt.Sprintf("type %s struct {\n", structName))
	sb.WriteString(fmt.Sprintf("    %s\n", mainEnt.Name))
	for _, ent := range relatedEnts {
		sb.WriteString(fmt.Sprintf("    %s\n", ent.Name))
	}
	sb.WriteString("}\n\n")
}

// getEntityColumns returns all column names for an entity
func (g *RepositoryGenerator) getEntityColumns(ent *entity.Entity) []string {
	fields := ent.GetGenerateableFields()
	var columns []string
	for _, field := range fields {
		columns = append(columns, util.ToSnakeCase(field.Name))
	}
	return columns
}

// generateGetByCompositePKMethod generates a GET method for entities with composite primary keys
func (g *RepositoryGenerator) generateGetByCompositePKMethod(sb *strings.Builder, ent *entity.Entity, fields []entity.Field) {
	tableName := ent.GetTableName()
	pkFields := ent.GetPrimaryKeyFields()

	sb.WriteString(fmt.Sprintf("func Get%s(ctx context.Context, db *sql.DB", ent.Name))

	// Build parameter list
	var params []string
	for _, pkField := range pkFields {
		params = append(params, fmt.Sprintf("%s %s", util.ToSnakeCase(pkField.Name), pkField.Type))
	}

	sb.WriteString(fmt.Sprintf(", %s) (*%s, error {\n", strings.Join(params, ", "), ent.Name))

	// Build column list
	var columns []string
	for _, field := range fields {
		if !field.ShouldGenerate() {
			continue
		}
		colName := util.ToSnakeCase(field.Name)
		columns = append(columns, colName)
	}

	// Build WHERE clause for composite PK
	var whereClause []string
	for i, pkField := range pkFields {
		whereClause = append(whereClause, fmt.Sprintf("%s = $%d", util.ToSnakeCase(pkField.Name), i+1))
	}

	query := fmt.Sprintf("SELECT %s FROM %s WHERE %s", strings.Join(columns, ", "), tableName, strings.Join(whereClause, " AND "))
	sb.WriteString(fmt.Sprintf("    query := `%s`\n", query))

	sb.WriteString(fmt.Sprintf("    var item %s\n", ent.Name))
	sb.WriteString("    err := db.QueryRowContext(ctx, query, ")

	// Add PK parameters (use snake_case for parameter names)
	var scanParams []string
	for _, pkField := range pkFields {
		scanParams = append(scanParams, util.ToSnakeCase(pkField.Name))
	}
	sb.WriteString(strings.Join(scanParams, ", "))
	sb.WriteString(").Scan(")

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

// generateDeleteByCompositePKMethod generates a DELETE method for entities with composite primary keys
func (g *RepositoryGenerator) generateDeleteByCompositePKMethod(sb *strings.Builder, ent *entity.Entity, fields []entity.Field) {
	tableName := ent.GetTableName()
	pkFields := ent.GetPrimaryKeyFields()

	sb.WriteString(fmt.Sprintf("func Delete%s(ctx context.Context, db *sql.DB", ent.Name))

	// Build parameter list
	var params []string
	for _, pkField := range pkFields {
		params = append(params, fmt.Sprintf("%s %s", util.ToSnakeCase(pkField.Name), pkField.Type))
	}

	sb.WriteString(fmt.Sprintf(", %s) error {\n", strings.Join(params, ", ")))

	// Build WHERE clause for composite PK
	var whereClause []string
	for i, pkField := range pkFields {
		whereClause = append(whereClause, fmt.Sprintf("%s = $%d", util.ToSnakeCase(pkField.Name), i+1))
	}

	query := fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, strings.Join(whereClause, " AND "))
	sb.WriteString(fmt.Sprintf("    query := `%s`\n", query))

	sb.WriteString("    _, err := db.ExecContext(ctx, query, ")

	// Add PK parameters
	var paramsList []string
	for _, pkField := range pkFields {
		paramsList = append(paramsList, util.ToSnakeCase(pkField.Name))
	}
	sb.WriteString(strings.Join(paramsList, ", "))
	sb.WriteString(")\n")

	sb.WriteString("    return err\n")
	sb.WriteString("}\n\n")
}
