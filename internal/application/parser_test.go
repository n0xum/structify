package application

import (
	"context"
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

	if result == nil {
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
