package adapter

import (
	"fmt"
	"strings"

	"github.com/n0xum/structify/internal/domain/entity"
	"github.com/n0xum/structify/internal/parser"
	"github.com/n0xum/structify/internal/util"
)

type ParserAdapter struct {
	patternMatcher *parser.PatternMatcher
	fieldMapper    *parser.FieldMapper
}

func NewParserAdapter() *ParserAdapter {
	return &ParserAdapter{
		patternMatcher: parser.NewPatternMatcher(),
		fieldMapper:    parser.NewFieldMapper(),
	}
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

	// Parse complex tags: check:, default:, index:, enum:, fk:, unique:, etc.
	for _, tag := range tags {
		switch {
		case strings.HasPrefix(tag, "check:"):
			domainField.CheckExpr = strings.TrimPrefix(tag, "check:")
		case strings.HasPrefix(tag, "default:"):
			domainField.DefaultVal = strings.TrimPrefix(tag, "default:")
		case strings.HasPrefix(tag, "enum:"):
			enumStr := strings.TrimPrefix(tag, "enum:")
			domainField.EnumValues = a.parseEnumValues(enumStr)
		case tag == "index":
			// Auto-generate index name
			domainField.IndexName = a.autoGenerateIndexName(pField.Name)
		case strings.HasPrefix(tag, "index:"):
			indexName := strings.TrimPrefix(tag, "index:")
			domainField.IndexName = indexName
			domainField.IndexGroup = indexName
		case tag == "unique_index":
			domainField.IndexName = a.autoGenerateIndexName(pField.Name)
			domainField.IsIndexUnique = true
		case strings.HasPrefix(tag, "unique_index:"):
			indexName := strings.TrimPrefix(tag, "unique_index:")
			domainField.IndexName = indexName
			domainField.IndexGroup = indexName
			domainField.IsIndexUnique = true
		case strings.HasPrefix(tag, "unique:"):
			// Named unique constraint for composite keys
			uniqueName := strings.TrimPrefix(tag, "unique:")
			domainField.IndexGroup = uniqueName
			domainField.IsUnique = true
		case strings.HasPrefix(tag, "fk:"):
			// Parse foreign key: fk:table,column[,on_delete:action][,on_update:action]
			a.parseForeignKey(tag, &domainField)
		}
	}

	return domainField
}

// parseEnumValues splits comma-separated enum values
func (a *ParserAdapter) parseEnumValues(enumStr string) []string {
	if enumStr == "" {
		return nil
	}
	parts := strings.Split(enumStr, ",")
	values := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			values = append(values, trimmed)
		}
	}
	return values
}

// autoGenerateIndexName generates an index name based on field name
func (a *ParserAdapter) autoGenerateIndexName(fieldName string) string {
	return util.ToSnakeCase(fieldName) + "_idx"
}

// parseForeignKey parses a foreign key tag
// Single FK format: fk:table,column[,on_delete:action][,on_update:action]
// Composite FK format: fk:constraint_name,table,column[,on_delete:action][,on_update:action]
// Examples:
//   - fk:users,id,on_delete:CASCADE (single FK with CASCADE)
//   - fk:fk_order,order_items,order_id (composite FK - first field)
//   - fk:fk_order,order_items,item_id (composite FK - second field)
func (a *ParserAdapter) parseForeignKey(tag string, field *entity.Field) {
	// Remove "fk:" prefix
	fkStr := strings.TrimPrefix(tag, "fk:")
	if fkStr == "" {
		return
	}

	// Parse the FK parts (split by comma but not inside known action values)
	parts := a.parseFKParts(fkStr)
	if len(parts) < 2 {
		return
	}

	// Separate cascade options from main parts
	var mainParts []string
	var cascadeParts []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if strings.HasPrefix(trimmed, "on_delete:") || strings.HasPrefix(trimmed, "on_update:") {
			cascadeParts = append(cascadeParts, trimmed)
		} else {
			mainParts = append(mainParts, trimmed)
		}
	}

	var fkRef *entity.FKReference

	if len(mainParts) == 2 {
		// Single FK: fk:table,column
		fkRef = &entity.FKReference{
			Table:  mainParts[0],
			Column: mainParts[1],
		}
	} else if len(mainParts) == 3 {
		// Composite FK: fk:constraint_name,table,column
		field.FKGroup = mainParts[0]
		fkRef = &entity.FKReference{
			Table:  mainParts[1],
			Column: mainParts[2],
		}
	} else {
		return // Invalid format
	}

	// Parse cascade options (only from the first field of composite FK)
	for _, part := range cascadeParts {
		if strings.HasPrefix(part, "on_delete:") {
			field.FKOnDelete = strings.TrimPrefix(part, "on_delete:")
		} else if strings.HasPrefix(part, "on_update:") {
			field.FKOnUpdate = strings.TrimPrefix(part, "on_update:")
		}
	}

	field.FKReference = fkRef
}

// parseFKParts splits a FK definition into parts
// Handles commas inside values (like function calls)
func (a *ParserAdapter) parseFKParts(fkStr string) []string {
	var parts []string
	var current strings.Builder
	inParens := 0

	for i := 0; i < len(fkStr); i++ {
		c := fkStr[i]

		switch c {
		case '(':
			inParens++
			current.WriteByte(c)
		case ')':
			inParens--
			current.WriteByte(c)
		case ',':
			if inParens > 0 {
				// Comma inside function call, keep it
				current.WriteByte(c)
			} else {
				// Comma separates parts
				if current.Len() > 0 {
					parts = append(parts, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteByte(c)
		}
	}

	// Add the last part
	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}

// tagParserState holds the parsing state for the tag parser.
type tagParserState struct {
	parts       []string
	current     strings.Builder
	inBraces    int
	inQuotes    bool
	inTagValue  bool
}

func (a *ParserAdapter) parseTags(tag string) []string {
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}

	state := &tagParserState{
		parts: make([]string, 0),
	}

	for i := 0; i < len(tag); i++ {
		c := tag[i]
		a.handleTagChar(state, c, tag, i)
	}

	a.addFinalPart(state)
	return state.parts
}

// handleTagChar processes a single character during tag parsing.
func (a *ParserAdapter) handleTagChar(state *tagParserState, c byte, tag string, index int) {
	switch c {
	case '{', '}':
		a.handleBrace(state, c)
	case '"', '\'':
		a.handleQuote(state, c)
	case ',':
		a.handleComma(state, c, tag, index)
	case ' ':
		a.handleSpace(state)
	default:
		state.current.WriteByte(c)
	}

	a.updateTagValueState(state)
}

// handleBrace processes brace characters for nested structures.
func (a *ParserAdapter) handleBrace(state *tagParserState, c byte) {
	if c == '{' {
		state.inBraces++
	} else {
		state.inBraces--
	}
	state.current.WriteByte(c)
}

// handleQuote processes quote characters.
func (a *ParserAdapter) handleQuote(state *tagParserState, c byte) {
	state.inQuotes = !state.inQuotes
	state.current.WriteByte(c)
}

// handleComma processes comma characters, which may separate tags or be part of values.
func (a *ParserAdapter) handleComma(state *tagParserState, c byte, tag string, index int) {
	// Keep comma if inside braces or quotes
	if state.inBraces > 0 || state.inQuotes {
		state.current.WriteByte(c)
		return
	}

	// If not in a tag value, comma is a separator
	if !state.inTagValue {
		a.addCurrentPart(state)
		return
	}

	// We're in a tag value - check if this comma should split tags
	currentStr := state.current.String()
	if a.shouldSplitTagAtComma(currentStr, tag, index) {
		a.addCurrentPart(state)
		state.inTagValue = false
	} else {
		state.current.WriteByte(c)
	}
}

// handleSpace processes space characters.
func (a *ParserAdapter) handleSpace(state *tagParserState) {
	// Space is a separator between tags - ignore it unless in quotes or tag value
	if state.inTagValue || state.inQuotes {
		state.current.WriteByte(' ')
	}
}

// shouldSplitTagAtComma determines if a comma should split the current tag.
func (a *ParserAdapter) shouldSplitTagAtComma(currentStr, fullTag string, index int) bool {
	// fk: and enum: tags contain commas as part of their values
	if strings.HasPrefix(currentStr, "fk:") || strings.HasPrefix(currentStr, "enum:") {
		return false
	}

	// Check if the remaining string starts with a known tag prefix
	remaining := strings.TrimSpace(fullTag[index+1:])
	return a.startsWithKnownPrefix(remaining)
}

// updateTagValueState updates whether we're inside a tag value based on current content.
func (a *ParserAdapter) updateTagValueState(state *tagParserState) {
	currentStr := state.current.String()

	// Enter tag value mode when we see a value tag prefix
	for _, prefix := range []string{"check:", "default:", "enum:", "fk:"} {
		if strings.HasSuffix(currentStr, prefix) {
			state.inTagValue = true
			return
		}
	}

	// Exit tag value mode when we complete a simple tag
	if state.inTagValue && len(currentStr) > 0 {
		trimmed := strings.TrimSpace(currentStr)
		if a.isSimpleTag(trimmed) {
			state.inTagValue = false
		}
	}
}

// addCurrentPart adds the current buffer content as a part and resets the buffer.
func (a *ParserAdapter) addCurrentPart(state *tagParserState) {
	if state.current.Len() > 0 {
		trimmed := strings.TrimSpace(state.current.String())
		if trimmed != "" {
			state.parts = append(state.parts, trimmed)
		}
		state.current.Reset()
	}
}

// addFinalPart adds any remaining content as the final part.
func (a *ParserAdapter) addFinalPart(state *tagParserState) {
	a.addCurrentPart(state)
}

// isSimpleTag checks if the string is a simple tag (no value part)
func (a *ParserAdapter) isSimpleTag(s string) bool {
	simpleTags := []string{"pk", "unique", "-", "index", "unique_index"}
	for _, tag := range simpleTags {
		if s == tag {
			return true
		}
	}
	return false
}

// startsWithKnownPrefix checks if the string starts with a known tag prefix.
// Uses O(1) lookup for exact matches and falls back to prefix matching for value tags.
func (a *ParserAdapter) startsWithKnownPrefix(s string) bool {
	// Simple tags (exact match)
	simpleTags := map[string]bool{
		"pk":           true,
		"unique":       true,
		"-":            true,
		"index":        true,
		"unique_index": true,
	}
	if simpleTags[s] {
		return true
	}

	// Value tags (prefix match)
	valuePrefixes := []string{"check:", "default:", "enum:", "fk:"}
	for _, prefix := range valuePrefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
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

// ToRepositoryInterface converts a parsed interface + entity into a domain RepositoryInterface.
func (a *ParserAdapter) ToRepositoryInterface(iface *parser.Interface, ent *entity.Entity) *entity.RepositoryInterface {
	if iface == nil || ent == nil {
		return nil
	}

	repo := &entity.RepositoryInterface{
		Name:       iface.Name,
		EntityName: ent.Name,
		Package:    iface.PackageName,
	}

	for _, m := range iface.Methods {
		rm := entity.RepositoryMethod{
			Name:       m.Name,
			EntityName: ent.Name,
		}

		// Convert params
		for _, p := range m.Params {
			rm.Params = append(rm.Params, entity.MethodParam{
				Name: p.Name,
				Type: p.Type,
			})
		}

		// Determine return characteristics
		for _, ret := range m.Returns {
			if ret.BaseType == "error" || ret.Type == "error" {
				rm.ReturnsError = true
			} else if ret.BaseType == ent.Name {
				rm.HasEntityReturn = true
				rm.ReturnsSingle = !ret.IsSlice
			} else {
				// Scalar return (float64, int64, bool, string, etc.)
				rm.ScalarReturnType = ret.Type
			}
		}

		// Classify method
		rm.Kind = a.classifyMethod(m, ent)

		// Handle CustomSQL (explicit SQL comment)
		if m.SQLComment != "" {
			rm.Kind = entity.MethodCustomSQL
			rm.CustomSQL = m.SQLComment
		}

		// Handle FindBy field extraction
		if rm.Kind == entity.MethodFindBy {
			rm.FindByFields = a.extractFindByFields(m.Name)
		}

		// Handle SmartQuery SQL generation
		if rm.Kind == entity.MethodSmartQuery {
			a.processSmartQueryMethod(&rm, m, ent)
		}

		repo.Methods = append(repo.Methods, rm)
	}

	return repo
}

// classifyMethod determines the MethodKind from the method name.
func (a *ParserAdapter) classifyMethod(m parser.Method, ent *entity.Entity) entity.MethodKind {
	name := m.Name

	// 1. Check for explicit SQL comment (highest priority)
	if m.SQLComment != "" {
		return entity.MethodCustomSQL
	}

	// 2. Check for standard CRUD methods
	switch {
	case name == "Create":
		return entity.MethodCreate
	case name == "GetByID" || name == "Get":
		return entity.MethodGetByID
	case name == "Update":
		return entity.MethodUpdate
	case name == "Delete":
		return entity.MethodDelete
	case name == "List" || name == "ListAll":
		return entity.MethodList
	}

	// 3. Try existing FindBy pattern (for backward compatibility)
	if strings.HasPrefix(name, "FindBy") {
		return entity.MethodFindBy
	}

	// 4. Try smart pattern matching
	if a.patternMatcher.Match(name) != nil {
		return entity.MethodSmartQuery
	}

	// 5. Default to CustomSQL (requires SQL comment)
	return entity.MethodCustomSQL
}

// extractFindByFields parses field names from method names like "FindByEmail" or "FindByStatusAndRole".
func (a *ParserAdapter) extractFindByFields(methodName string) []string {
	suffix := strings.TrimPrefix(methodName, "FindBy")
	if suffix == "" {
		return nil
	}

	parts := strings.Split(suffix, "And")
	fields := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			fields = append(fields, trimmed)
		}
	}
	return fields
}

// processSmartQueryMethod generates SQL for smart query methods
func (a *ParserAdapter) processSmartQueryMethod(rm *entity.RepositoryMethod, m parser.Method, ent *entity.Entity) {
	// Match the pattern
	matched := a.patternMatcher.Match(m.Name)
	if matched == nil {
		// Pattern not found, leave as-is (will fall back to CustomSQL)
		return
	}

	// Generate SQL — use quoted table name so reserved words (e.g. "order") are safe.
	tableName := ent.GetQuotedTableName()
	sql, err := a.generateSmartQuerySQL(m.Name, tableName, ent, matched)
	if err != nil {
		return
	}

	rm.GeneratedSQL = sql
	rm.QueryPattern = matched.Pattern.Regex
}

// generateSmartQuerySQL generates SQL query for a smart query method
func (a *ParserAdapter) generateSmartQuerySQL(methodName, tableName string, ent *entity.Entity, matched *parser.MatchedPattern) (string, error) {
	var sb strings.Builder

	// Determine SELECT clause based on return type
	switch matched.ReturnType {
	case parser.ReturnCount:
		sb.WriteString("SELECT COUNT(*) FROM ")
	case parser.ReturnExists:
		sb.WriteString("SELECT EXISTS(SELECT 1 FROM ")
	case parser.ReturnDelete:
		sb.WriteString("DELETE FROM ")
	default:
		// ReturnSingle, ReturnMany - build column list
		fields := ent.GetGenerateableFields()
		var columns []string
		for _, f := range fields {
			columns = append(columns, util.ToSnakeCase(f.Name))
		}
		sb.WriteString("SELECT " + strings.Join(columns, ", ") + " FROM ")
	}

	sb.WriteString(tableName)

	// Add WHERE clause if conditions exist
	if len(matched.Conditions) > 0 {
		sb.WriteString(" WHERE ")
		var whereParts []string
		for _, cond := range matched.Conditions {
			part := cond.ColumnName + " " + cond.Operator + " $" + fmt.Sprint(cond.ParamIndex)
			if cond.LogicalOp != "" {
				part += " " + cond.LogicalOp
			}
			whereParts = append(whereParts, part)
		}
		sb.WriteString(strings.Join(whereParts, " "))
	}

	// Add ORDER BY clause
	if matched.OrderBy != "" {
		sb.WriteString(" ORDER BY " + matched.OrderBy)
	}

	// Add LIMIT clause
	if matched.Limit > 0 {
		sb.WriteString(" LIMIT ")
		if matched.Limit == 1 {
			sb.WriteString("1")
		} else {
			sb.WriteString(fmt.Sprint(matched.Limit))
		}
	}

	// Close EXISTS subquery if needed
	if matched.ReturnType == parser.ReturnExists {
		sb.WriteString(")")
	}

	return sb.String(), nil
}

// ToInterfaceMap converts parsed interfaces to domain format.
func (a *ParserAdapter) ToInterfaceMap(interfaces map[string][]*parser.Interface) map[string][]*entity.RepositoryInterface {
	if interfaces == nil {
		return nil
	}

	result := make(map[string][]*entity.RepositoryInterface)
	for pkgName, ifaces := range interfaces {
		for _, iface := range ifaces {
			// Create a placeholder without entity binding
			repo := &entity.RepositoryInterface{
				Name:    iface.Name,
				Package: pkgName,
			}
			result[pkgName] = append(result[pkgName], repo)
		}
	}
	return result
}
