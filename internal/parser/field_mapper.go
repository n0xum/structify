package parser

import (
	"strings"
	"unicode"

	"github.com/n0xum/structify/internal/util"
)

// FieldMapper handles conversion between Go field names and SQL column names
type FieldMapper struct{}

// NewFieldMapper creates a new field mapper
func NewFieldMapper() *FieldMapper {
	return &FieldMapper{}
}

// FieldToColumn converts a Go field name to SQL column name (snake_case)
// Examples:
//   - ID → i_d
//   - UserID → user_i_d
//   - Email → email
//   - IsActive → is_active
//   - HTMLContent → h_t_m_l_content
func (fm *FieldMapper) FieldToColumn(fieldName string) string {
	return util.ToSnakeCase(fieldName)
}

// ColumnToField converts an SQL column name to Go field name (PascalCase)
// Examples:
//   - i_d → ID
//   - user_i_d → UserID
//   - email → Email
//   - is_active → IsActive
//   - h_t_m_l_content → HTMLContent
func (fm *FieldMapper) ColumnToField(columnName string) string {
	parts := strings.Split(columnName, "_")
	var result string
	for _, part := range parts {
		if part == "" {
			continue
		}
		// Capitalize first letter, keep rest as-is for acronyms
		result += string(unicode.ToUpper(rune(part[0]))) + part[1:]
	}
	return result
}

// ExtractFields extracts field conditions from a method name
// Parses patterns like:
//   - "ByEmail" → [{FieldName: "Email"}]
//   - "ByEmailAndRole" → [{FieldName: "Email"}, {FieldName: "Role"}]
//   - "ByEmailOrUsername" → [{FieldName: "Email", LogicalOp: "OR"}, {FieldName: "Username"}]
//   - "ByCreatedAtGreaterThan" → [{FieldName: "CreatedAt", Operator: ">"}]
//   - "ByActiveOrderByCreatedAtDesc" → [{FieldName: "Active"}]  (OrderBy stripped)
func (fm *FieldMapper) ExtractFields(methodName string) []FieldCondition {
	// Find the "By" keyword
	byIndex := strings.Index(methodName, "By")
	if byIndex == -1 {
		return nil
	}

	suffix := methodName[byIndex+2:] // Everything after "By"
	if suffix == "" {
		return nil
	}

	// Strip the "OrderBy..." suffix before parsing field conditions so that
	// "OrderBy" is never misread as a logical OR separator.
	if obIdx := strings.Index(suffix, "OrderBy"); obIdx != -1 {
		suffix = suffix[:obIdx]
	}
	if suffix == "" {
		return nil
	}

	return fm.parseFieldConditions(suffix)
}

// parseFieldConditions parses the suffix after "By" into field conditions
func (fm *FieldMapper) parseFieldConditions(suffix string) []FieldCondition {
	var conditions []FieldCondition
	current := strings.Builder{}
	paramIndex := 1

	i := 0
	for i < len(suffix) {
		c := suffix[i]

		// Check for "And" operator
		if i+3 <= len(suffix) && suffix[i:i+3] == "And" {
			if current.Len() > 0 {
				cond := fm.buildCondition(current.String(), paramIndex, "AND")
				conditions = append(conditions, cond)
				paramIndex++
			}
			current.Reset()
			i += 3 // Skip "And"
			continue
		}

		// Check for "Or" operator
		if i+2 <= len(suffix) && suffix[i:i+2] == "Or" {
			if current.Len() > 0 {
				cond := fm.buildCondition(current.String(), paramIndex, "OR")
				conditions = append(conditions, cond)
				paramIndex++
			}
			current.Reset()
			i += 2 // Skip "Or"
			continue
		}

		current.WriteByte(c)
		i++
	}

	// Add the last condition
	if current.Len() > 0 {
		cond := fm.buildCondition(current.String(), paramIndex, "")
		conditions = append(conditions, cond)
	}

	return conditions
}

// buildCondition creates a FieldCondition from a parsed string
func (fm *FieldMapper) buildCondition(s string, paramIndex int, logicalOp string) FieldCondition {
	// Check for operators at the end
	cond := FieldCondition{
		ParamIndex: paramIndex,
		LogicalOp:  logicalOp,
		Operator:   "=",
	}

	// Check for comparison operators
	if strings.HasSuffix(s, "GreaterThan") {
		cond.FieldName = strings.TrimSuffix(s, "GreaterThan")
		cond.Operator = ">"
	} else if strings.HasSuffix(s, "LessThan") {
		cond.FieldName = strings.TrimSuffix(s, "LessThan")
		cond.Operator = "<"
	} else if strings.HasSuffix(s, "GreaterThanOrEqual") {
		cond.FieldName = strings.TrimSuffix(s, "GreaterThanOrEqual")
		cond.Operator = ">="
	} else if strings.HasSuffix(s, "LessThanOrEqual") {
		cond.FieldName = strings.TrimSuffix(s, "LessThanOrEqual")
		cond.Operator = "<="
	} else if strings.HasSuffix(s, "Like") {
		cond.FieldName = strings.TrimSuffix(s, "Like")
		cond.Operator = "LIKE"
	} else if strings.HasSuffix(s, "StartingWith") {
		cond.FieldName = strings.TrimSuffix(s, "StartingWith")
		cond.Operator = "LIKE"
	} else if strings.HasSuffix(s, "EndingWith") {
		cond.FieldName = strings.TrimSuffix(s, "EndingWith")
		cond.Operator = "LIKE"
	} else if strings.HasSuffix(s, "Containing") {
		cond.FieldName = strings.TrimSuffix(s, "Containing")
		cond.Operator = "LIKE"
	} else if strings.HasSuffix(s, "NotIn") {
		cond.FieldName = strings.TrimSuffix(s, "NotIn")
		cond.Operator = "NOT IN"
	} else if strings.HasSuffix(s, "In") {
		cond.FieldName = strings.TrimSuffix(s, "In")
		cond.Operator = "IN"
	} else if strings.HasSuffix(s, "IsNull") {
		cond.FieldName = strings.TrimSuffix(s, "IsNull")
		cond.Operator = "IS NULL"
	} else if strings.HasSuffix(s, "IsNotNull") {
		cond.FieldName = strings.TrimSuffix(s, "IsNotNull")
		cond.Operator = "IS NOT NULL"
	} else if strings.HasPrefix(s, "Not") {
		cond.FieldName = strings.TrimPrefix(s, "Not")
		cond.Operator = "!="
		cond.Negated = true
	} else {
		// Default: simple field name
		cond.FieldName = s
	}

	// Convert field name to column name
	cond.ColumnName = fm.FieldToColumn(cond.FieldName)

	return cond
}

// ExtractOrderBy extracts ORDER BY clause from method name
// Examples:
//   - "OrderByCreatedAtDesc" → "created_at DESC"
//   - "OrderByUsername" → "username ASC"
//   - "OrderByCreatedAtAsc" → "created_at ASC"
func (fm *FieldMapper) ExtractOrderBy(methodName string) string {
	// Find "OrderBy" keyword. We use LastIndex to handle cases where "Order"
	// is part of the entity name (e.g., "ExistsOrderByCustomerIDOrderByCreatedAt").
	obIndex := strings.LastIndex(methodName, "OrderBy")
	if obIndex == -1 {
		return ""
	}

	// An "OrderBy" clause is only valid if it's NOT acting as the primary
	// "By" separator for an entity ending in "Order".
	// We check if the part before "OrderBy" is just a standard prefix.
	byIndex := strings.Index(methodName, "By")
	if byIndex != -1 && byIndex >= obIndex {
		prefix := methodName[:obIndex]
		if fm.isStandardPrefix(prefix) {
			return ""
		}
	}

	suffix := methodName[obIndex+7:] // Everything after "OrderBy"
	if suffix == "" {
		return ""
	}

	// Check for direction suffix
	direction := "ASC"
	if strings.HasSuffix(suffix, "Desc") {
		direction = "DESC"
		suffix = strings.TrimSuffix(suffix, "Desc")
	} else if strings.HasSuffix(suffix, "Asc") {
		suffix = strings.TrimSuffix(suffix, "Asc")
	}

	fieldName := suffix
	columnName := fm.FieldToColumn(fieldName)

	return columnName + " " + direction
}

// isStandardPrefix checks if a string is a standard query method prefix
func (fm *FieldMapper) isStandardPrefix(s string) bool {
	switch s {
	case "List", "Find", "Exists", "Count", "Delete", "First":
		return true
	}
	if strings.HasPrefix(s, "Top") {
		rest := s[3:]
		if rest == "" {
			return true
		}
		for _, r := range rest {
			if r < '0' || r > '9' {
				return false
			}
		}
		return true
	}
	return false
}

// ExtractLimit extracts LIMIT value from method name
// Examples:
//   - "FirstUserByEmail" → 1
//   - "Top10UsersByRole" → 10
func (fm *FieldMapper) ExtractLimit(methodName string) int {
	// Check for "First" prefix
	if strings.HasPrefix(methodName, "First") {
		return 1
	}

	// Check for "Top" prefix
	if strings.HasPrefix(methodName, "Top") {
		suffix := strings.TrimPrefix(methodName, "Top")
		// Extract number
		i := 0
		for i < len(suffix) && suffix[i] >= '0' && suffix[i] <= '9' {
			i++
		}
		if i > 0 {
			numStr := suffix[:i]
			var limit int
			for _, c := range numStr {
				limit = limit*10 + int(c-'0')
			}
			return limit
		}
	}

	return 0
}

// IsCountPattern checks if the method name is a count pattern
func (fm *FieldMapper) IsCountPattern(methodName string) bool {
	return strings.HasPrefix(methodName, "Count")
}

// IsExistsPattern checks if the method name is an exists pattern
func (fm *FieldMapper) IsExistsPattern(methodName string) bool {
	return strings.HasPrefix(methodName, "Exists")
}

// IsDeletePattern checks if the method name is a delete pattern
func (fm *FieldMapper) IsDeletePattern(methodName string) bool {
	return strings.HasPrefix(methodName, "Delete")
}

// IsListPattern checks if the method name is a list/find pattern
func (fm *FieldMapper) IsListPattern(methodName string) bool {
	return strings.HasPrefix(methodName, "List") || strings.HasPrefix(methodName, "Find")
}

// GetReturnType determines the return type based on method name
func (fm *FieldMapper) GetReturnType(methodName string, returnsSlice bool) ReturnType {
	if fm.IsCountPattern(methodName) {
		return ReturnCount
	}
	if fm.IsExistsPattern(methodName) {
		return ReturnExists
	}
	if fm.IsDeletePattern(methodName) {
		return ReturnDelete
	}
	if returnsSlice {
		return ReturnMany
	}
	return ReturnSingle
}
