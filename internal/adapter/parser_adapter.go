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
	return entity.ToSnakeCase(fieldName) + "_idx"
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

func (a *ParserAdapter) parseTags(tag string) []string {
	if tag == "" {
		return nil
	}
	tag = strings.TrimSpace(tag)
	if tag == "" {
		return nil
	}

	// Smart split that handles commas inside {...}, quotes, and tag values with spaces
	var parts []string
	var current strings.Builder
	inBraces := 0
	inQuotes := false
	inTagValue := false  // Inside a tag value (check:, default:, enum:, etc.)

	for i := 0; i < len(tag); i++ {
		c := tag[i]

		switch c {
		case '{':
			inBraces++
			current.WriteByte(c)
		case '}':
			inBraces--
			current.WriteByte(c)
		case '"', '\'':
			inQuotes = !inQuotes
			current.WriteByte(c)
		case ',':
			// Keep comma if inside braces or quotes
			if inBraces > 0 || inQuotes {
				current.WriteByte(c)
			} else if inTagValue {
				// Check if we're inside an fk: tag (which contains commas)
				currentStr := current.String()
				if strings.HasPrefix(currentStr, "fk:") {
					// Keep comma - we're inside an fk: tag
					current.WriteByte(c)
				} else if strings.HasPrefix(currentStr, "enum:") {
					// Keep comma - we're inside an enum: tag
					current.WriteByte(c)
				} else {
					// Check if the remaining string starts with a known tag prefix
					remaining := strings.TrimSpace(tag[i+1:])
					if a.startsWithKnownPrefix(remaining) {
						// Split here - comma separates two tags (e.g., check:..., default:...)
						if current.Len() > 0 {
							trimmed := strings.TrimSpace(current.String())
							if trimmed != "" {
								parts = append(parts, trimmed)
							}
							current.Reset()
							inTagValue = false
						}
					} else {
						// Keep comma as part of tag value
						current.WriteByte(c)
					}
				}
			} else {
				// Comma is a separator between tags
				if current.Len() > 0 {
					trimmed := strings.TrimSpace(current.String())
					if trimmed != "" {
						parts = append(parts, trimmed)
					}
					current.Reset()
				}
			}
		case ' ':
			// Keep space if inside tag value or quotes
			if inTagValue || inQuotes {
				current.WriteByte(c)
			}
			// Otherwise space is a separator between tags
		default:
			current.WriteByte(c)
		}

		// Track if we're inside a tag value (check:, default:, enum:, fk:, etc.)
		currentStr := current.String()
		// Check if we just finished typing a known tag prefix
		for _, prefix := range []string{"check:", "default:", "enum:", "fk:"} {
			if strings.HasSuffix(currentStr, prefix) {
				inTagValue = true
				break
			}
		}
		// Exit tag value mode when we hit a comma (already handled above) or
		// when the current word is a known simple tag
		if inTagValue && len(currentStr) > 0 {
			trimmed := strings.TrimSpace(currentStr)
			if a.isSimpleTag(trimmed) {
				inTagValue = false
			}
		}
	}

	// Add the last part
	if current.Len() > 0 {
		trimmed := strings.TrimSpace(current.String())
		if trimmed != "" {
			parts = append(parts, trimmed)
		}
	}

	return parts
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

// startsWithKnownPrefix checks if the string starts with a known tag prefix
func (a *ParserAdapter) startsWithKnownPrefix(s string) bool {
	knownPrefixes := []string{"pk", "unique", "check:", "default:", "enum:", "index", "unique_index", "fk:"}
	for _, prefix := range knownPrefixes {
		if strings.HasPrefix(s, prefix) {
			return true
		}
	}
	return false
}

// isKnownTagPrefix checks if the string starts with a known tag prefix
func (a *ParserAdapter) isKnownTagPrefix(s string) bool {
	knownPrefixes := []string{"pk", "unique", "check:", "default:", "index", "fk:"}
	for _, prefix := range knownPrefixes {
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

// ToSnakeCase converts a PascalCase or camelCase string to snake_case
func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
