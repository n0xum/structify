package entity

import "errors"

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
	for i := range e.Fields {
		if e.Fields[i].IsPrimary {
			return &e.Fields[i]
		}
	}
	return nil
}

func (e *Entity) HasPrimaryKey() bool {
	return e.GetPrimaryKeyField() != nil
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
	return ToSnakeCase(e.Name)
}
