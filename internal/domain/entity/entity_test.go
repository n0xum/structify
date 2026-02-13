package entity

import (
	"testing"
)

func TestEntityValidate(t *testing.T) {
	tests := []struct {
		name    string
		entity  *Entity
		wantErr error
	}{
		{
			name: "valid entity",
			entity: &Entity{
				Name: "User",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
					{Name: "Name", Type: "string"},
				},
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			entity: &Entity{
				Fields: []Field{
					{Name: "ID", Type: "int64"},
				},
			},
			wantErr: ErrEntityNameRequired,
		},
		{
			name: "no fields",
			entity: &Entity{
				Name:   "User",
				Fields: []Field{},
			},
			wantErr: ErrEntityHasNoFields,
		},
		{
			name: "duplicate field names",
			entity: &Entity{
				Name: "User",
				Fields: []Field{
					{Name: "ID", Type: "int64"},
					{Name: "ID", Type: "string"},
				},
			},
			wantErr: ErrDuplicateFieldName,
		},
		{
			name: "empty field name",
			entity: &Entity{
				Name: "User",
				Fields: []Field{
					{Name: "", Type: "int64"},
				},
			},
			wantErr: ErrInvalidFieldName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.entity.Validate()
			if err != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEntityGetPrimaryKeyField(t *testing.T) {
	entity := &Entity{
		Name: "User",
		Fields: []Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Name", Type: "string"},
		},
	}

	pk := entity.GetPrimaryKeyField()
	if pk == nil {
		t.Fatal("GetPrimaryKeyField() returned nil")
	}
	if pk.Name != "ID" {
		t.Errorf("GetPrimaryKeyField() Name = %v, want ID", pk.Name)
	}
}

func TestEntityHasPrimaryKey(t *testing.T) {
	t.Run("has primary key", func(t *testing.T) {
		entity := &Entity{
			Name: "User",
			Fields: []Field{
				{Name: "ID", Type: "int64", IsPrimary: true},
			},
		}
		if !entity.HasPrimaryKey() {
			t.Error("HasPrimaryKey() returned false")
		}
	})

	t.Run("no primary key", func(t *testing.T) {
		entity := &Entity{
			Name: "User",
			Fields: []Field{
				{Name: "Name", Type: "string"},
			},
		}
		if entity.HasPrimaryKey() {
			t.Error("HasPrimaryKey() returned true")
		}
	})
}

func TestEntityGetGenerateableFields(t *testing.T) {
	entity := &Entity{
		Name: "User",
		Fields: []Field{
			{Name: "ID", Type: "int64", IsPrimary: true},
			{Name: "Name", Type: "string"},
			{Name: "Secret", Type: "string", IsIgnored: true},
		},
	}

	fields := entity.GetGenerateableFields()
	if len(fields) != 2 {
		t.Errorf("GetGenerateableFields() returned %d fields, want 2", len(fields))
	}
}

func TestEntityGetTableName(t *testing.T) {
	tests := []struct {
		name         string
		tableName    string
		entityName   string
		wantTableName string
	}{
		{
			name:         "custom table name",
			tableName:    "users",
			entityName:   "User",
			wantTableName: "users",
		},
		{
			name:         "default table name",
			tableName:    "",
			entityName:   "User",
			wantTableName: "user",
		},
		{
			name:         "camel case to snake",
			tableName:    "",
			entityName:   "UserProfile",
			wantTableName: "user_profile",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entity := &Entity{
				Name:      tt.entityName,
				TableName: tt.tableName,
			}
			got := entity.GetTableName()
			if got != tt.wantTableName {
				t.Errorf("GetTableName() = %v, want %v", got, tt.wantTableName)
			}
		})
	}
}

func TestFieldShouldGenerate(t *testing.T) {
	tests := []struct {
		name string
		field Field
		want bool
	}{
		{
			name: "normal field",
			field: Field{Name: "ID", Type: "int64"},
			want: true,
		},
		{
			name: "ignored field",
			field: Field{Name: "Secret", Type: "string", IsIgnored: true},
			want: false,
		},
		{
			name: "empty type",
			field: Field{Name: "Unknown", Type: ""},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.field.ShouldGenerate(); got != tt.want {
				t.Errorf("Field.ShouldGenerate() = %v, want %v", got, tt.want)
			}
		})
	}
}
