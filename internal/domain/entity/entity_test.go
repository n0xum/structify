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

func TestEntityHasCompositePrimaryKey(t *testing.T) {
	tests := []struct {
		name string
		entity *Entity
		want bool
	}{
		{
			name: "single primary key",
			entity: &Entity{
				Name: "User",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
				},
			},
			want: false,
		},
		{
			name: "composite primary key",
			entity: &Entity{
				Name: "OrderItem",
				Fields: []Field{
					{Name: "OrderID", Type: "int64", IsPrimary: true},
					{Name: "ItemID", Type: "int64", IsPrimary: true},
				},
			},
			want: true,
		},
		{
			name: "no primary key",
			entity: &Entity{
				Name: "Log",
				Fields: []Field{
					{Name: "Message", Type: "string"},
				},
			},
			want: false,
		},
		{
			name: "three primary keys",
			entity: &Entity{
				Name: "TriplePK",
				Fields: []Field{
					{Name: "A", Type: "int64", IsPrimary: true},
					{Name: "B", Type: "int64", IsPrimary: true},
					{Name: "C", Type: "int64", IsPrimary: true},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entity.HasCompositePrimaryKey(); got != tt.want {
				t.Errorf("HasCompositePrimaryKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntityGetUniqueConstraints(t *testing.T) {
	tests := []struct {
		name string
		entity *Entity
		wantConstraintCount int
		wantCompositeCount int
	}{
		{
			name: "no unique constraints",
			entity: &Entity{
				Name: "Basic",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
				},
			},
			wantConstraintCount: 0,
			wantCompositeCount: 0,
		},
		{
			name: "single unique constraint",
			entity: &Entity{
				Name: "User",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
					{Name: "Email", Type: "string", IsUnique: true},
				},
			},
			wantConstraintCount: 1,
			wantCompositeCount: 0,
		},
		{
			name: "composite unique constraint",
			entity: &Entity{
				Name: "Assignment",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
					{Name: "TenantID", Type: "int64", IsUnique: true, IndexGroup: "uq_tenant_user"},
					{Name: "UserID", Type: "int64", IsUnique: true, IndexGroup: "uq_tenant_user"},
				},
			},
			wantConstraintCount: 1,
			wantCompositeCount: 1,
		},
		{
			name: "multiple composite unique constraints",
			entity: &Entity{
				Name: "Complex",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
					{Name: "Email", Type: "string", IsUnique: true},
					{Name: "A1", Type: "int64", IsUnique: true, IndexGroup: "uq_a"},
					{Name: "A2", Type: "int64", IsUnique: true, IndexGroup: "uq_a"},
					{Name: "B1", Type: "int64", IsUnique: true, IndexGroup: "uq_b"},
					{Name: "B2", Type: "int64", IsUnique: true, IndexGroup: "uq_b"},
				},
			},
			wantConstraintCount: 3,
			wantCompositeCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			constraints := tt.entity.GetUniqueConstraints()
			if len(constraints) != tt.wantConstraintCount {
				t.Errorf("GetUniqueConstraints() count = %d, want %d", len(constraints), tt.wantConstraintCount)
			}

			compositeCount := 0
			for _, fields := range constraints {
				if len(fields) > 1 {
					compositeCount++
				}
			}
			if compositeCount != tt.wantCompositeCount {
				t.Errorf("GetUniqueConstraints() composite count = %d, want %d", compositeCount, tt.wantCompositeCount)
			}
		})
	}
}

func TestEntityHasCompositeUniqueConstraints(t *testing.T) {
	tests := []struct {
		name string
		entity *Entity
		want bool
	}{
		{
			name: "no unique constraints",
			entity: &Entity{
				Name: "Basic",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
				},
			},
			want: false,
		},
		{
			name: "single unique constraint",
			entity: &Entity{
				Name: "User",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
					{Name: "Email", Type: "string", IsUnique: true},
				},
			},
			want: false,
		},
		{
			name: "composite unique constraint",
			entity: &Entity{
				Name: "Assignment",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
					{Name: "TenantID", Type: "int64", IsUnique: true, IndexGroup: "uq_tenant_user"},
					{Name: "UserID", Type: "int64", IsUnique: true, IndexGroup: "uq_tenant_user"},
				},
			},
			want: true,
		},
		{
			name: "mixed single and composite",
			entity: &Entity{
				Name: "Mixed",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
					{Name: "Email", Type: "string", IsUnique: true},
					{Name: "A1", Type: "int64", IsUnique: true, IndexGroup: "uq_a"},
					{Name: "A2", Type: "int64", IsUnique: true, IndexGroup: "uq_a"},
				},
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.entity.HasCompositeUniqueConstraints(); got != tt.want {
				t.Errorf("HasCompositeUniqueConstraints() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEntityGetPrimaryKeyFields(t *testing.T) {
	tests := []struct {
		name string
		entity *Entity
		wantCount int
	}{
		{
			name: "no primary key",
			entity: &Entity{
				Name: "Log",
				Fields: []Field{
					{Name: "Message", Type: "string"},
				},
			},
			wantCount: 0,
		},
		{
			name: "single primary key",
			entity: &Entity{
				Name: "User",
				Fields: []Field{
					{Name: "ID", Type: "int64", IsPrimary: true},
				},
			},
			wantCount: 1,
		},
		{
			name: "composite primary key",
			entity: &Entity{
				Name: "OrderItem",
				Fields: []Field{
					{Name: "OrderID", Type: "int64", IsPrimary: true},
					{Name: "ItemID", Type: "int64", IsPrimary: true},
				},
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkFields := tt.entity.GetPrimaryKeyFields()
			if len(pkFields) != tt.wantCount {
				t.Errorf("GetPrimaryKeyFields() count = %d, want %d", len(pkFields), tt.wantCount)
			}
		})
	}
}
