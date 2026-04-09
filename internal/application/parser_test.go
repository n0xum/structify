package application

import (
	"context"
	"github.com/n0xum/structify/internal/domain/entity"
	"testing"
)

func TestNewParserWrapper(t *testing.T) {
	wrapper := NewParserWrapper()

	if wrapper == nil {
		t.Fatal("NewParserWrapper() returned nil")
	}

	if wrapper.parser == nil {
		t.Error("parser field is nil")
	}

	if wrapper.adapter == nil {
		t.Error("adapter field is nil")
	}
}

func TestParserWrapperParseFiles(t *testing.T) {
	wrapper := NewParserWrapper()

	ctx := context.Background()

	result, err := wrapper.ParseFiles(ctx, []string{"../../test/fixtures/user.go"})
	if err != nil {
		t.Fatalf("ParseFiles() error = %v", err)
	}

	if len(result) == 0 {
		t.Fatal("ParseFiles() returned nil")
	}

	totalCount := 0
	for _, entities := range result {
		totalCount += len(entities)
	}

	if totalCount == 0 {
		t.Error("ParseFiles() returned no entities")
	}
}

func TestParserWrapperParseFilesEmpty(t *testing.T) {
	wrapper := NewParserWrapper()

	ctx := context.Background()

	result, err := wrapper.ParseFiles(ctx, []string{})
	if err != nil {
		t.Fatalf("ParseFiles() error = %v", err)
	}

	if result != nil {
		totalCount := 0
		for _, entities := range result {
			totalCount += len(entities)
		}
		if totalCount != 0 {
			t.Error("ParseFiles() should return empty result map")
		}
	}
}

func TestParserWrapperParseFilesNonExistent(t *testing.T) {
	wrapper := NewParserWrapper()

	ctx := context.Background()

	_, err := wrapper.ParseFiles(ctx, []string{"nonexistent.go"})
	if err == nil {
		t.Error("ParseFiles() should return error for non-existent file")
	}
}

func TestParserWrapperParseInterfaces(t *testing.T) {
	wrapper := NewParserWrapper()
	ctx := context.Background()

	// First parse the entity
	entitiesMap, err := wrapper.ParseFiles(ctx, []string{"../../test/fixtures/user.go"})
	if err != nil {
		t.Fatalf("ParseFiles() error = %v", err)
	}

	// Find the User entity
	var userEnt *entity.Entity
	for _, ents := range entitiesMap {
		for _, ent := range ents {
			if ent.Name == "User" {
				userEnt = ent
				break
			}
		}
	}
	if userEnt == nil {
		t.Fatal("User entity not found")
	}

	// Parse interfaces
	result, err := wrapper.ParseInterfaces(ctx, []string{"../../test/fixtures/user_repository.go"}, userEnt)
	if err != nil {
		t.Fatalf("ParseInterfaces() error = %v", err)
	}

	if len(result) == 0 {
		t.Fatal("ParseInterfaces() returned nil")
	}

	if len(result) == 0 {
		t.Error("ParseInterfaces() returned no interfaces")
	}

	// Verify the UserRepository interface was parsed
	found := false
	for _, repo := range result {
		if repo.Name == "UserRepository" {
			found = true
			if len(repo.Methods) == 0 {
				t.Error("UserRepository has no methods")
			}
			break
		}
	}
	if !found {
		t.Error("UserRepository interface not found")
	}
}

func TestParserWrapperParseInterfacesNonExistent(t *testing.T) {
	wrapper := NewParserWrapper()
	ctx := context.Background()

	userEnt := &entity.Entity{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}}}

	_, err := wrapper.ParseInterfaces(ctx, []string{"nonexistent.go"}, userEnt)
	if err == nil {
		t.Error("ParseInterfaces() should return error for non-existent file")
	}
}

func TestParserWrapperParseInterfacesNoInterfaces(t *testing.T) {
	wrapper := NewParserWrapper()
	ctx := context.Background()

	// Parse a file without interfaces (user.go has no interfaces)
	userEnt := &entity.Entity{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}}}

	result, err := wrapper.ParseInterfaces(ctx, []string{"../../test/fixtures/user.go"}, userEnt)
	if err != nil {
		t.Fatalf("ParseInterfaces() error = %v", err)
	}

	// Accept nil or empty slice for files with no interfaces
	if len(result) != 0 {
		t.Errorf("ParseInterfaces() should return nil or empty slice for file with no interfaces, got %d items", len(result))
	}
}

func TestParserWrapperParseInterfacesEmptyPaths(t *testing.T) {
	wrapper := NewParserWrapper()
	ctx := context.Background()

	userEnt := &entity.Entity{Name: "User", Fields: []entity.Field{{Name: "ID", Type: "int64", IsPrimary: true}}}

	result, err := wrapper.ParseInterfaces(ctx, []string{}, userEnt)
	if err != nil {
		t.Fatalf("ParseInterfaces() error = %v", err)
	}

	// Accept nil or empty slice for empty paths
	if len(result) != 0 {
		t.Errorf("ParseInterfaces() should return nil or empty slice for empty paths, got %d items", len(result))
	}
}
