package entity

import (
	"errors"

	"github.com/n0xum/structify/internal/util"
)

var ErrInvalidEntity = errors.New("invalid entity")
var ErrEntityNameRequired = errors.New("entity name is required")
var ErrEntityHasNoFields = errors.New("entity must have at least one field")
var ErrDuplicateFieldName = errors.New("duplicate field name")
var ErrInvalidFieldName = errors.New("invalid field name")

type Entity struct {
	Name      string
	Fields    []Field
	TableName string
	Package   string
}

func (e *Entity) Validate() error {
	if e.Name == "" {
		return ErrEntityNameRequired
	}

	if len(e.Fields) == 0 {
		return ErrEntityHasNoFields
	}

	fieldNames := make(map[string]bool)
	for _, field := range e.Fields {
		if field.Name == "" {
			return ErrInvalidFieldName
		}
		if fieldNames[field.Name] {
			return ErrDuplicateFieldName
		}
		fieldNames[field.Name] = true
	}

	return nil
}

func (e *Entity) GetPrimaryKeyField() *Field {
	pkFields := e.GetPrimaryKeyFields()
	if len(pkFields) == 1 {
		return &pkFields[0]
	}
	return nil
}

// GetPrimaryKeyFields returns all primary key fields
func (e *Entity) GetPrimaryKeyFields() []Field {
	var pkFields []Field
	for _, field := range e.Fields {
		if field.IsPrimary {
			pkFields = append(pkFields, field)
		}
	}
	return pkFields
}

func (e *Entity) HasPrimaryKey() bool {
	return len(e.GetPrimaryKeyFields()) > 0
}

// HasCompositePrimaryKey returns true if the entity has a composite primary key
func (e *Entity) HasCompositePrimaryKey() bool {
	return len(e.GetPrimaryKeyFields()) > 1
}

// GetUniqueConstraints returns all unique constraint definitions
func (e *Entity) GetUniqueConstraints() map[string][]Field {
	constraints := make(map[string][]Field)

	// Group fields by their IndexGroup to find composite constraints
	groups := make(map[string][]Field)

	for _, field := range e.Fields {
		if field.IsUnique {
			if field.IndexGroup != "" {
				// Part of a named unique constraint
				groups[field.IndexGroup] = append(groups[field.IndexGroup], field)
			} else {
				// Single unique constraint - use field name as key
				constraints[field.Name] = []Field{field}
			}
		}
	}

	// Add groups to constraints (only composite ones, single fields already added)
	for groupName, fields := range groups {
		if len(fields) > 1 {
			constraints[groupName] = fields
		}
	}

	return constraints
}

// HasCompositeUniqueConstraints returns true if the entity has composite unique constraints
func (e *Entity) HasCompositeUniqueConstraints() bool {
	constraints := e.GetUniqueConstraints()
	for _, fields := range constraints {
		if len(fields) > 1 {
			return true
		}
	}
	return false
}

func (e *Entity) GetGenerateableFields() []Field {
	var fields []Field
	for _, f := range e.Fields {
		if f.ShouldGenerate() {
			fields = append(fields, f)
		}
	}
	return fields
}

func (e *Entity) GetTableName() string {
	if e.TableName != "" {
		return e.TableName
	}
	return util.ToSnakeCase(e.Name)
}

// GetQuotedTableName returns the table name wrapped in double quotes so that
// PostgreSQL reserved words (e.g. "order", "user") are always safe to use in DML.
func (e *Entity) GetQuotedTableName() string {
	return `"` + e.GetTableName() + `"`
}
