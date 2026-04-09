package entity

// FKReference represents a foreign key reference to another table
type FKReference struct {
	Table   string   // Referenced table name
	Column  string   // Referenced column name
	Columns []string // Multiple columns for composite FK (future)
}

type Field struct {
	Name      string
	Type      string
	IsPrimary bool
	IsUnique  bool
	IsIgnored bool
	TableName string

	// CheckExpr stores the expression for CHECK constraint
	// Parsed from db:"check:expression" tag
	CheckExpr string

	// DefaultVal stores the DEFAULT value
	// Parsed from db:"default:value" tag
	DefaultVal string

	// IndexName stores the index name
	// Parsed from db:"index" (auto-generated), db:"index:name", or db:"unique_index"
	IndexName string

	// IsIndexUnique indicates if this is a unique index
	// Parsed from db:"unique_index" or db:"unique_index:name"
	IsIndexUnique bool

	// IndexGroup groups fields together for composite indexes
	// Fields with the same IndexName belong to the same index
	IndexGroup string

	// EnumValues stores the allowed enum values
	// Parsed from db:"enum:value1,value2,value3" tag
	EnumValues []string

	// FKReference stores the foreign key reference
	// Parsed from db:"fk:Table,Column" tag
	FKReference *FKReference

	// FKOnDelete stores the ON DELETE clause
	// Parsed from db:"fk:...,on_delete:CASCADE"
	FKOnDelete string

	// FKOnUpdate stores the ON UPDATE clause
	// Parsed from db:"fk:...,on_update:CASCADE"
	FKOnUpdate string

	// FKGroup groups fields together for composite foreign keys
	// Fields with the same FKGroup form a composite FK constraint
	// Parsed from db:"fk:constraint_name,table,column"
	FKGroup string
}

func (f *Field) ShouldGenerate() bool {
	if f.IsIgnored {
		return false
	}
	if f.Type == "" {
		return false
	}
	return true
}
