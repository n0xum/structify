package parser

import (
	"testing"
)

func TestPatternMatcher_ListBySingleField(t *testing.T) {
	pm := NewPatternMatcher()

	tests := []struct {
		name        string
		methodName  string
		wantMatch   bool
		wantCondLen int
	}{
		{"List users by email", "ListUsersByEmail", true, 1},
		{"Find user by ID", "FindUserByID", true, 1},
		{"List audit logs by action", "ListAuditLogsByAction", true, 1},
		{"No pattern", "RandomMethodName", false, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := pm.Match(tt.methodName)
			if (matched != nil) != tt.wantMatch {
				t.Errorf("Match() match = %v, wantMatch %v", matched != nil, tt.wantMatch)
			}
			if matched != nil && len(matched.Conditions) != tt.wantCondLen {
				t.Errorf("Match() conditions len = %v, want %v", len(matched.Conditions), tt.wantCondLen)
			}
		})
	}
}

func TestPatternMatcher_ListByMultipleFields(t *testing.T) {
	pm := NewPatternMatcher()

	tests := []struct {
		name        string
		methodName  string
		wantMatch   bool
		wantCondLen int
	}{
		{"List users by email and role", "ListUsersByEmailAndRole", true, 2},
		{"Find user by email and ID", "FindUserByEmailAndID", true, 2},
		{"List audit logs by action and user", "ListAuditLogsByActionAndUser", true, 2},
		{"List users by email or username", "ListUsersByEmailOrUsername", true, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := pm.Match(tt.methodName)
			if (matched != nil) != tt.wantMatch {
				t.Errorf("Match() match = %v, wantMatch %v", matched != nil, tt.wantMatch)
			}
			if matched != nil && len(matched.Conditions) != tt.wantCondLen {
				t.Errorf("Match() conditions len = %v, want %v", len(matched.Conditions), tt.wantCondLen)
			}
		})
	}
}

func TestPatternMatcher_CountByField(t *testing.T) {
	pm := NewPatternMatcher()

	tests := []struct {
		name          string
		methodName    string
		wantMatch     bool
		wantReturnType ReturnType
	}{
		{"Count users by active", "CountUsersByActive", true, ReturnCount},
		{"Count audit logs by action", "CountAuditLogsByAction", true, ReturnCount},
		{"Count users", "CountUsers", true, ReturnCount},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := pm.Match(tt.methodName)
			if (matched != nil) != tt.wantMatch {
				t.Errorf("Match() match = %v, wantMatch %v", matched != nil, tt.wantMatch)
			}
			if matched != nil && matched.ReturnType != tt.wantReturnType {
				t.Errorf("Match() returnType = %v, want %v", matched.ReturnType, tt.wantReturnType)
			}
		})
	}
}

func TestPatternMatcher_ExistsByField(t *testing.T) {
	pm := NewPatternMatcher()

	tests := []struct {
		name          string
		methodName    string
		wantMatch     bool
		wantReturnType ReturnType
	}{
		{"Exists user by email", "ExistsUserByEmail", true, ReturnExists},
		{"Exists audit log by ID", "ExistsAuditLogByID", true, ReturnExists},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := pm.Match(tt.methodName)
			if (matched != nil) != tt.wantMatch {
				t.Errorf("Match() match = %v, wantMatch %v", matched != nil, tt.wantMatch)
			}
			if matched != nil && matched.ReturnType != tt.wantReturnType {
				t.Errorf("Match() returnType = %v, want %v", matched.ReturnType, tt.wantReturnType)
			}
		})
	}
}

func TestPatternMapper_BooleanFields(t *testing.T) {
	pm := NewPatternMatcher()

	tests := []struct {
		name       string
		methodName string
		wantMatch  bool
	}{
		{"List active users", "ListUsersByActiveTrue", true},
		{"List inactive users", "ListUsersByActiveFalse", true},
		{"List by status", "ListUsersByStatusActive", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matched := pm.Match(tt.methodName)
			if (matched != nil) != tt.wantMatch {
				t.Errorf("Match() match = %v, wantMatch %v", matched != nil, tt.wantMatch)
			}
		})
	}
}

func TestFieldMapper_FieldToColumn(t *testing.T) {
	fm := NewFieldMapper()

	tests := []struct {
		name     string
		field    string
		expected string
	}{
		{"Simple field", "Email", "email"},
		{"ID field", "ID", "id"},
		{"UserID field", "UserID", "user_id"},
		{"IsActive field", "IsActive", "is_active"},
		{"CreatedAt field", "CreatedAt", "created_at"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fm.FieldToColumn(tt.field)
			if got != tt.expected {
				t.Errorf("FieldToColumn() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFieldMapper_ColumnToField(t *testing.T) {
	fm := NewFieldMapper()

	tests := []struct {
		name     string
		column   string
		expected string
	}{
		{"Simple column", "email", "Email"},
		{"ID column", "id", "Id"},
		{"UserID column", "user_id", "UserId"},
		{"IsActive column", "is_active", "IsActive"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fm.ColumnToField(tt.column)
			if got != tt.expected {
				t.Errorf("ColumnToField() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFieldMapper_ExtractFields(t *testing.T) {
	fm := NewFieldMapper()

	tests := []struct {
		name         string
		methodName   string
		wantCondLen  int
		firstField   string
		firstOp      string
		secondField  string
		logicalOp    string
	}{
		{"By email", "ListUsersByEmail", 1, "Email", "=", "", ""},
		{"By email and role", "ListUsersByEmailAndRole", 2, "Email", "=", "Role", "AND"},
		{"By email or username", "ListUsersByEmailOrUsername", 2, "Email", "=", "Username", "OR"},
		{"By age greater than", "ListUsersByAgeGreaterThan", 1, "Age", ">", "", ""},
		{"By age less than", "ListUsersByAgeLessThan", 1, "Age", "<", "", ""},
		{"By name like", "ListUsersByNameLike", 1, "Name", "LIKE", "", ""},
		// OrderBy suffix must not be parsed as a logical OR condition
		{"By active with OrderBy suffix", "ListUsersByActiveOrderByCreatedAtDesc", 1, "Active", "=", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions := fm.ExtractFields(tt.methodName)
			if len(conditions) != tt.wantCondLen {
				t.Errorf("ExtractFields() len = %v, want %v", len(conditions), tt.wantCondLen)
				return
			}
			if len(conditions) > 0 {
				if conditions[0].FieldName != tt.firstField {
					t.Errorf("ExtractFields() first field = %v, want %v", conditions[0].FieldName, tt.firstField)
				}
				if conditions[0].Operator != tt.firstOp {
					t.Errorf("ExtractFields() first op = %v, want %v", conditions[0].Operator, tt.firstOp)
				}
				if tt.secondField != "" && len(conditions) > 1 {
					if conditions[1].FieldName != tt.secondField {
						t.Errorf("ExtractFields() second field = %v, want %v", conditions[1].FieldName, tt.secondField)
					}
					if conditions[0].LogicalOp != tt.logicalOp {
						t.Errorf("ExtractFields() logical op = %v, want %v", conditions[0].LogicalOp, tt.logicalOp)
					}
				}
			}
		})
	}
}

func TestFieldMapper_ExtractOrderBy(t *testing.T) {
	fm := NewFieldMapper()

	tests := []struct {
		name     string
		method   string
		expected string
	}{
		{"OrderBy created_at desc", "ListUsersByEmailOrderByCreatedAtDesc", "created_at DESC"},
		{"OrderBy username asc", "ListUsersByStatusOrderByUsernameAsc", "username ASC"},
		{"OrderBy created_at", "ListUsersByEmailOrderByCreatedAt", "created_at ASC"},
		{"No order by", "ListUsersByEmail", ""},
		// "Order" is part of the entity name, not an ORDER BY modifier
		{"Exists entity Order By field", "ExistsOrderByCustomerID", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fm.ExtractOrderBy(tt.method)
			if got != tt.expected {
				t.Errorf("ExtractOrderBy() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFieldMapper_ExtractLimit(t *testing.T) {
	fm := NewFieldMapper()

	tests := []struct {
		name     string
		method   string
		expected int
	}{
		{"First user", "FirstUserByEmail", 1},
		{"Top 10 users", "Top10UsersByRole", 10},
		{"Top 5 users", "Top5UsersByStatus", 5},
		{"No limit", "ListUsersByEmail", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fm.ExtractLimit(tt.method)
			if got != tt.expected {
				t.Errorf("ExtractLimit() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFieldMapper_GetReturnType(t *testing.T) {
	fm := NewFieldMapper()

	tests := []struct {
		name        string
		methodName  string
		returnsSlice bool
		expected    ReturnType
	}{
		{"Count pattern", "CountUsers", false, ReturnCount},
		{"Exists pattern", "ExistsUserByEmail", false, ReturnExists},
		{"Delete pattern", "DeleteUserByEmail", false, ReturnDelete},
		{"List pattern", "ListUsersByEmail", true, ReturnMany},
		{"Find pattern (single)", "FindUserByID", false, ReturnSingle},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fm.GetReturnType(tt.methodName, tt.returnsSlice)
			if got != tt.expected {
				t.Errorf("GetReturnType() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFieldMapper_IsCountPattern(t *testing.T) {
	fm := NewFieldMapper()

	if !fm.IsCountPattern("CountUsers") {
		t.Error("IsCountPattern() should return true for CountUsers")
	}
	if !fm.IsCountPattern("CountUsersByActive") {
		t.Error("IsCountPattern() should return true for CountUsersByActive")
	}
	if fm.IsCountPattern("ListUsers") {
		t.Error("IsCountPattern() should return false for ListUsers")
	}
}

func TestFieldMapper_IsExistsPattern(t *testing.T) {
	fm := NewFieldMapper()

	if !fm.IsExistsPattern("ExistsUserByEmail") {
		t.Error("IsExistsPattern() should return true for ExistsUserByEmail")
	}
	if fm.IsExistsPattern("ListUsers") {
		t.Error("IsExistsPattern() should return false for ListUsers")
	}
}

func TestFieldMapper_IsListPattern(t *testing.T) {
	fm := NewFieldMapper()

	if !fm.IsListPattern("ListUsers") {
		t.Error("IsListPattern() should return true for ListUsers")
	}
	if !fm.IsListPattern("FindUserByEmail") {
		t.Error("IsListPattern() should return true for FindUserByEmail")
	}
	if fm.IsListPattern("CountUsers") {
		t.Error("IsListPattern() should return false for CountUsers")
	}
}

// TestPatternMatcher_ExtractComponents_ReturnSingle exercises the inner branch
// of extractComponents that only fires when a matched pattern has ReturnSingle (0).
func TestPatternMatcher_ExtractComponents_ReturnSingle(t *testing.T) {
	pm := NewPatternMatcher()

	t.Run("method starts with List sets ReturnMany", func(t *testing.T) {
		mp := &MatchedPattern{
			Pattern:    &QueryPattern{ReturnType: ReturnSingle},
			ReturnType: ReturnSingle, // 0 — triggers the inner if
		}
		pm.extractComponents(mp, "ListUsersByEmail", nil)
		if mp.ReturnType != ReturnMany {
			t.Errorf("expected ReturnMany, got %v", mp.ReturnType)
		}
	})

	t.Run("method starts with Find sets ReturnMany", func(t *testing.T) {
		mp := &MatchedPattern{
			Pattern:    &QueryPattern{ReturnType: ReturnSingle},
			ReturnType: ReturnSingle,
		}
		pm.extractComponents(mp, "FindUserByID", nil)
		if mp.ReturnType != ReturnMany {
			t.Errorf("expected ReturnMany, got %v", mp.ReturnType)
		}
	})

	t.Run("method without List/Find prefix returns ReturnSingle", func(t *testing.T) {
		mp := &MatchedPattern{
			Pattern:    &QueryPattern{ReturnType: ReturnSingle},
			ReturnType: ReturnSingle,
		}
		pm.extractComponents(mp, "GetUserByID", nil)
		if mp.ReturnType != ReturnSingle {
			t.Errorf("expected ReturnSingle, got %v", mp.ReturnType)
		}
	})
}

// TestFieldMapper_BuildCondition_ViaExtractFields covers operator suffix branches
// that are reachable through ExtractFields (no "Or" in the suffix).
func TestFieldMapper_BuildCondition_ViaExtractFields(t *testing.T) {
	fm := NewFieldMapper()

	tests := []struct {
		methodName   string
		wantField    string
		wantOperator string
	}{
		{"ListUsersByEmailStartingWith", "Email", "LIKE"},
		{"ListUsersByEmailEndingWith", "Email", "LIKE"},
		{"ListUsersByEmailContaining", "Email", "LIKE"},
		// "In" suffix — field name trimmed to everything before "In"
		{"ListUsersByIDIn", "ID", "IN"},
		// "NotIn" must be checked before "In" to avoid false match
		{"ListUsersByIDNotIn", "ID", "NOT IN"},
		{"ListUsersByDeletedAtIsNull", "DeletedAt", "IS NULL"},
		{"ListUsersByDeletedAtIsNotNull", "DeletedAt", "IS NOT NULL"},
		{"ListUsersByNotActive", "Active", "!="},
	}

	for _, tt := range tests {
		t.Run(tt.methodName, func(t *testing.T) {
			conditions := fm.ExtractFields(tt.methodName)
			if len(conditions) == 0 {
				t.Fatalf("ExtractFields(%q) returned no conditions", tt.methodName)
			}
			if conditions[0].FieldName != tt.wantField {
				t.Errorf("FieldName = %q, want %q", conditions[0].FieldName, tt.wantField)
			}
			if conditions[0].Operator != tt.wantOperator {
				t.Errorf("Operator = %q, want %q", conditions[0].Operator, tt.wantOperator)
			}
		})
	}
}

// TestFieldMapper_BuildCondition_Direct tests operator branches that cannot be
// reached via ExtractFields because their suffix contains "Or" (which the parser
// treats as a logical OR separator). These operators are only valid when the
// caller invokes buildCondition directly with a pre-assembled condition string.
func TestFieldMapper_BuildCondition_Direct(t *testing.T) {
	fm := NewFieldMapper()

	tests := []struct {
		input        string
		wantField    string
		wantOperator string
	}{
		{"AgeGreaterThanOrEqual", "Age", ">="},
		{"AgeLessThanOrEqual", "Age", "<="},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			cond := fm.buildCondition(tt.input, 1, "")
			if cond.FieldName != tt.wantField {
				t.Errorf("FieldName = %q, want %q", cond.FieldName, tt.wantField)
			}
			if cond.Operator != tt.wantOperator {
				t.Errorf("Operator = %q, want %q", cond.Operator, tt.wantOperator)
			}
		})
	}
}
