package mapper

import "strings"

type TypeMapping struct {
	GoType       string
	PostgresType string
	Constraints  []string
	IsNotNull    bool
}

var typeMappings = map[string]TypeMapping{
	"int":             {GoType: "int", PostgresType: "INTEGER", IsNotNull: true},
	"int8":            {GoType: "int8", PostgresType: "SMALLINT", IsNotNull: true},
	"int16":           {GoType: "int16", PostgresType: "SMALLINT", IsNotNull: true},
	"int32":           {GoType: "int32", PostgresType: "INTEGER", IsNotNull: true},
	"int64":           {GoType: "int64", PostgresType: "BIGINT", IsNotNull: true},
	"uint":            {GoType: "uint", PostgresType: "BIGINT", Constraints: []string{"CHECK (\" || column_name || \" >= 0)"}, IsNotNull: true},
	"uint8":           {GoType: "uint8", PostgresType: "SMALLINT", Constraints: []string{"CHECK (\" || column_name || \" >= 0)"}, IsNotNull: true},
	"uint16":          {GoType: "uint16", PostgresType: "SMALLINT", Constraints: []string{"CHECK (\" || column_name || \" >= 0)"}, IsNotNull: true},
	"uint32":          {GoType: "uint32", PostgresType: "INTEGER", Constraints: []string{"CHECK (\" || column_name || \" >= 0)"}, IsNotNull: true},
	"uint64":          {GoType: "uint64", PostgresType: "BIGINT", Constraints: []string{"CHECK (\" || column_name || \" >= 0)"}, IsNotNull: true},
	"float32":         {GoType: "float32", PostgresType: "REAL", IsNotNull: true},
	"float64":         {GoType: "float64", PostgresType: "DOUBLE PRECISION", IsNotNull: true},
	"string":          {GoType: "string", PostgresType: "VARCHAR(255)", IsNotNull: false},
	"bool":            {GoType: "bool", PostgresType: "BOOLEAN", IsNotNull: true},
	"time.Time":       {GoType: "time.Time", PostgresType: "TIMESTAMP", IsNotNull: true},
	"json.RawMessage": {GoType: "json.RawMessage", PostgresType: "JSONB", IsNotNull: false},
	"[]byte":          {GoType: "[]byte", PostgresType: "BYTEA", IsNotNull: false},
}

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) MapType(goType string) TypeMapping {
	baseType := m.getBaseType(goType)
	if mapping, ok := typeMappings[baseType]; ok {
		return mapping
	}
	return TypeMapping{GoType: goType, PostgresType: "TEXT", IsNotNull: false}
}

func (m *Mapper) getBaseType(goType string) string {
	goType = strings.TrimPrefix(goType, "*")
	goType = strings.TrimPrefix(goType, "[]")

	if idx := strings.Index(goType, "["); idx != -1 {
		goType = goType[:idx]
	}
	if idx := strings.Index(goType, "}."); idx != -1 {
		goType = goType[:idx]
	}

	parts := strings.Split(goType, ".")
	if len(parts) > 1 {
		goType = parts[len(parts)-1]
	}

	return goType
}

func (m *Mapper) FormatColumnDefinition(fieldName string, mapping TypeMapping, tags []string) string {
	def := mapping.PostgresType

	constraints := m.parseConstraints(tags)
	constraints = append(constraints, mapping.Constraints...)

	if m.HasTag(tags, "-") {
		return ""
	}

	if m.HasTag(tags, "pk") {
		def += " PRIMARY KEY"
	}

	if mapping.IsNotNull && !m.HasTag(tags, "pk") {
		def += " NOT NULL"
	}

	if m.HasTag(tags, "unique") {
		def += " UNIQUE"
	}

	for _, constraint := range constraints {
		if constraint != "" {
			def += " " + constraint
		}
	}

	return def
}

func (m *Mapper) parseConstraints(tags []string) []string {
	var constraints []string
	for _, tag := range tags {
		if strings.HasPrefix(tag, "check:") {
			constraint := strings.TrimPrefix(tag, "check:")
			constraints = append(constraints, "CHECK ("+constraint+")")
		}
		if strings.HasPrefix(tag, "default:") {
			defaultVal := strings.TrimPrefix(tag, "default:")
			constraints = append(constraints, "DEFAULT "+defaultVal)
		}
		if strings.HasPrefix(tag, "enum:") {
			enumVals := strings.TrimPrefix(tag, "enum:")
			constraints = append(constraints, "CHECK (column_name IN ("+enumVals+"))")
		}
	}
	return constraints
}

func (m *Mapper) HasTag(tags []string, tag string) bool {
	for _, t := range tags {
		if t == tag {
			return true
		}
	}
	return false
}

func toSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, r)
	}
	return strings.ToLower(string(result))
}
