package parser

import (
	"regexp"
	"strings"
)

// ReturnType represents the type of result a query returns
type ReturnType int

const (
	ReturnSingle ReturnType = iota // Returns a single entity
	ReturnMany                      // Returns a slice of entities
	ReturnCount                     // Returns a count (int64)
	ReturnExists                    // Returns a boolean
	ReturnDelete                    // Returns deletion result
)

// QueryPattern defines a regex pattern for matching method names and generating SQL
type QueryPattern struct {
	Regex       string
	SQLTemplate string
	ReturnType  ReturnType
	Examples    []string
}

// MatchedPattern represents a successfully matched pattern with extracted components
type MatchedPattern struct {
	Pattern     *QueryPattern
	Matches     []string                 // Regex capture groups
	Conditions  []FieldCondition         // WHERE clause conditions
	OrderBy     string                   // ORDER BY clause
	Limit       int                      // LIMIT value
	ReturnType  ReturnType
}

// FieldCondition represents a single WHERE clause condition
type FieldCondition struct {
	FieldName    string
	ColumnName   string
	Operator     string // "=", ">", "<", "LIKE", "IN", etc.
	ParamIndex   int    // 1-based parameter index
	LogicalOp    string // "AND", "OR"
	Negated      bool
}

// PatternMatcher matches method names against query patterns
type PatternMatcher struct {
	patterns    []*CompiledPattern
	fieldMapper *FieldMapper
}

// CompiledPattern holds a compiled regex pattern
type CompiledPattern struct {
	Pattern *QueryPattern
	Regex   *regexp.Regexp
}

// NewPatternMatcher creates a new pattern matcher with default patterns
func NewPatternMatcher() *PatternMatcher {
	pm := &PatternMatcher{
		fieldMapper: NewFieldMapper(),
	}
	pm.initDefaultPatterns()
	return pm
}

// Match attempts to match a method name against all known patterns
func (pm *PatternMatcher) Match(methodName string) *MatchedPattern {
	for _, cp := range pm.patterns {
		if matches := cp.Regex.FindStringSubmatch(methodName); matches != nil {
			return pm.buildMatchedPattern(cp.Pattern, matches, methodName)
		}
	}
	return nil
}

// buildMatchedPattern creates a MatchedPattern from regex matches
func (pm *PatternMatcher) buildMatchedPattern(pattern *QueryPattern, matches []string, methodName string) *MatchedPattern {
	mp := &MatchedPattern{
		Pattern:    pattern,
		Matches:    matches,
		ReturnType: pattern.ReturnType,
	}

	// Extract conditions, ordering, and limiting based on pattern type
	pm.extractComponents(mp, methodName, matches)

	return mp
}

// extractComponents parses the method name to extract query components
func (pm *PatternMatcher) extractComponents(mp *MatchedPattern, methodName string, matches []string) {
	// Use FieldMapper to extract components
	conditions := pm.fieldMapper.ExtractFields(methodName)
	mp.Conditions = conditions
	mp.OrderBy = pm.fieldMapper.ExtractOrderBy(methodName)
	mp.Limit = pm.fieldMapper.ExtractLimit(methodName)

	// Also set ReturnType based on method name if not already set
	if mp.ReturnType == 0 {
		returnsSlice := false
		for _, ret := range []string{"List", "Find"} {
			if strings.HasPrefix(methodName, ret) {
				returnsSlice = true
				break
			}
		}
		mp.ReturnType = pm.fieldMapper.GetReturnType(methodName, returnsSlice)
	}
}

// GenerateSQL generates SQL query from a matched pattern
func (pm *PatternMatcher) GenerateSQL(methodName, tableName string, ent MinimalEntity) (string, error) {
	mp := pm.Match(methodName)
	if mp == nil {
		return "", nil
	}

	var sb strings.Builder

	// Determine SELECT clause based on return type
	switch mp.ReturnType {
	case ReturnCount:
		sb.WriteString("SELECT COUNT(*) FROM ")
	case ReturnExists:
		sb.WriteString("SELECT EXISTS(SELECT 1 FROM ")
	case ReturnDelete:
		sb.WriteString("DELETE FROM ")
	default:
		// ReturnSingle, ReturnMany - generate column list from entity
		if ent != nil {
			var columns []string
			for _, f := range ent.GetFieldNames() {
				columns = append(columns, pm.fieldMapper.FieldToColumn(f))
			}
			sb.WriteString("SELECT " + strings.Join(columns, ", ") + " FROM ")
		} else {
			sb.WriteString("SELECT * FROM ")
		}
	}

	sb.WriteString(tableName)

	// Add WHERE clause if conditions exist
	if len(mp.Conditions) > 0 {
		sb.WriteString(" WHERE ")
		var whereParts []string
		for _, cond := range mp.Conditions {
			part := cond.ColumnName + " " + cond.Operator + " $" + string(rune('0'+cond.ParamIndex))
			if cond.LogicalOp != "" {
				part += " " + cond.LogicalOp
			}
			whereParts = append(whereParts, part)
		}
		sb.WriteString(strings.Join(whereParts, " "))
	}

	// Add ORDER BY clause
	if mp.OrderBy != "" {
		sb.WriteString(" ORDER BY " + mp.OrderBy)
	}

	// Add LIMIT clause
	if mp.Limit > 0 {
		if mp.Limit == 1 {
			sb.WriteString(" LIMIT 1")
		} else {
			sb.WriteString(" LIMIT " + string(rune('0'+mp.Limit)))
		}
	}

	// Close EXISTS subquery if needed
	if mp.ReturnType == ReturnExists {
		sb.WriteString(")")
	}

	return sb.String(), nil
}

// initDefaultPatterns initializes all default query patterns
func (pm *PatternMatcher) initDefaultPatterns() {
	pm.patterns = []*CompiledPattern{
		// Result Type Patterns: List/Find + Entity + By + Field
		pm.compilePattern(&QueryPattern{
			Regex:       `^(List|Find)(\w+)By(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} = $1",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByEmail", "FindUserByID"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^(List|Find)(\w+)By(\w+)And(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {f1} = $1 AND {f2} = $2",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByEmailAndRole", "FindUserByEmailAndID"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^(List|Find)(\w+)By(\w+)Or(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {f1} = $1 OR {f2} = $2",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByEmailOrUsername", "FindUserByIDOrEmail"},
		}),

		// Count Patterns
		pm.compilePattern(&QueryPattern{
			Regex:       `^Count(\w+)By(\w+)$`,
			SQLTemplate: "SELECT COUNT(*) FROM {table} WHERE {field} = $1",
			ReturnType:  ReturnCount,
			Examples:    []string{"CountUsersByActive", "CountAuditLogsByAction"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^Count(\w+)$`,
			SQLTemplate: "SELECT COUNT(*) FROM {table}",
			ReturnType:  ReturnCount,
			Examples:    []string{"CountUsers", "CountAuditLogs"},
		}),

		// Exists Patterns
		pm.compilePattern(&QueryPattern{
			Regex:       `^Exists(\w+)By(\w+)$`,
			SQLTemplate: "SELECT EXISTS(SELECT 1 FROM {table} WHERE {field} = $1)",
			ReturnType:  ReturnExists,
			Examples:    []string{"ExistsUserByEmail", "ExistsAuditLogByID"},
		}),

		// Boolean/Status Patterns
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)ByActive(True|False)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE active = $1",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByActiveTrue", "ListUsersByActiveFalse"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)ByStatus(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE status = $1",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByStatusActive", "ListOrdersByStatusPending"},
		}),

		// Comparison Patterns
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)By(\w+)GreaterThan(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} > $1",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListAuditLogsByCreatedAtGreaterThan", "ListUsersByAgeGreaterThan18"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)By(\w+)LessThan(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} < $1",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListAuditLogsByCreatedAtLessThan", "ListUsersByAgeLessThan65"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)By(\w+)Between(\w+)And(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} BETWEEN $1 AND $2",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByAgeBetween18And65", "ListAuditLogsByDateBetweenStartAndEnd"},
		}),

		// String Patterns
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)By(\w+)Like(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} LIKE $1",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByEmailLike"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)By(\w+)StartingWith(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} LIKE $1 || '%'",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByEmailStartingWith"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)By(\w+)EndingWith(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} LIKE '%' || $1",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByEmailEndingWith"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)By(\w+)Containing(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} LIKE '%' || $1 || '%'",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByEmailContaining"},
		}),

		// Collection Patterns
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)By(\w+)In(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} IN ($1)",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByIDIn"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)By(\w+)NotIn(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} NOT IN ($1)",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByIDNotIn"},
		}),

		// Ordering Patterns
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)By(\w+)OrderBy(\w+)(Asc|Desc)?$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} = $1 ORDER BY {order}",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByEmailOrderByCreatedAtDesc", "ListUsersByStatusOrderByUsernameAsc"},
		}),

		// Limiting Patterns
		pm.compilePattern(&QueryPattern{
			Regex:       `^(First|Top)(\d+)?(\w+)By(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} = $2 LIMIT {limit}",
			ReturnType:  ReturnMany,
			Examples:    []string{"FirstUserByEmail", "Top10UsersByRole"},
		}),

		// Time-based Patterns
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)ByRecent(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE {field} > $1 ORDER BY {field} DESC",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListAuditLogsByRecentCreatedAt"},
		}),
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)ByLatest(\d+)$`,
			SQLTemplate: "SELECT * FROM {table} ORDER BY created DESC LIMIT $1",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListAuditLogsByLatest10"},
		}),

		// Negation Patterns
		pm.compilePattern(&QueryPattern{
			Regex:       `^List(\w+)ByNot(\w+)$`,
			SQLTemplate: "SELECT * FROM {table} WHERE NOT ({field})",
			ReturnType:  ReturnMany,
			Examples:    []string{"ListUsersByNotActive", "ListUsersByNotDeleted"},
		}),
	}
}

// compilePattern compiles a regex pattern
func (pm *PatternMatcher) compilePattern(pattern *QueryPattern) *CompiledPattern {
	return &CompiledPattern{
		Pattern: pattern,
		Regex:   regexp.MustCompile(pattern.Regex),
	}
}

// MinimalEntity is a minimal interface for entity field access
type MinimalEntity interface {
	GetFieldNames() []string
}
